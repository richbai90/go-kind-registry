package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"runtime"
	"time"
	"github.com/richbai90/bundle-containers/resources"
	"github.com/spf13/afero"
)

// The root of the standard file system using chroot for simplicity
var FSRoot string
// Standard file system using chroot for simplicity
var FS afero.Fs

// Return a wrapper around the native OS calls when in production
// Return a memmapped FS for testing
func getAferoBackend(chroot bool) afero.Fs {
	var fs afero.Fs
	cacheFs := afero.NewMemMapFs()
	if Testing() {
		fs = afero.NewMemMapFs()
	} else {
		fs = afero.NewOsFs()
	}

	// cache files in memory for faster access
	fs = afero.NewCacheOnReadFs(fs, cacheFs, time.Duration(86.64e13))

	if chroot {
		// chroot the filesystem to the config directory
		fs = afero.NewBasePathFs(fs, FSRoot)
	}

	return fs

}

func getConfigDir() string {
	root := os.Getenv("BUNDLE_DIR")
	if root == "" {
		root = os.Getenv("HOME")
	}

	return path.Join(root, ".config", "kind")
}

// Open a file from the embedded file system
func OpenResource(filename string) fs.File {
	f, e := resources.Resources.Open(fmt.Sprintf("resources/%s", filename))
	if e != nil {
		log.Fatal(fmt.Sprintf("Unable to open file %s: %s", filename, e.Error()))
	}

	return f
}

// Open the filename corresponding to the os currently running
func OpenWhenOS(filenameMap map[string]string) (fs.File, string) {
	os := runtime.GOOS
	Filename := filenameMap[os]
	if Filename == "" {
		log.Fatal("Unsupported OS: ", os)
	}
	return OpenResource(Filename), Filename
}

// Read a file from the FS or fatally error
func ReadFile(filepath string) []byte {
	bytes, err := afero.ReadFile(FS, filepath);
	if  err != nil {
		log.Fatal("Unable to read file ", filepath, ". Error: ", err.Error())
	}
	
	return bytes
}

// Create a JSON file on the FS from an interface or fatally error.
// If successful, call cb with the new file.
// cb must close the file when it is done.
func MakeJsonFile(filepath string, JSON interface{}, cb func (fs.File)) {
	str, err := json.Marshal(JSON)
	if err != nil {
		log.Fatal("Failed to encode JSON from object. Error: ", err.Error())
	}

	WriteFile(filepath, str, cb)
}

// Copy the given file to the specified filepath
// When copy is complete executes the callback as cb(inputFile, copiedFile)
// cb must close both files when it is done
func CopyFile(file fs.File, filepath string, cb func(fs.File, fs.File)) {
	outFile, err := FS.Create(filepath)
	if err != nil {
		log.Fatal("Unable to unpack embedded file. Error: ", err.Error())
	}

	// Callback when copy is complete
	defer cb(file, outFile)

	buf := make([]byte, 1024)

	for {
		// read a chunk
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			outFile.Close()
			file.Close()
			log.Fatal("There was a problem reading the provided file during Copy operation: ", err.Error())
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := outFile.Write(buf[:n]); err != nil {
			outFile.Close()
			file.Close()
			log.Fatal("There was a problem writing the new file during Copy operation: ", err.Error())
		}
	}
}

// Create a file on the FS and write to it, or fatally fail
// If successful, call cb with the created file
// cb must close the file when it is done
func WriteFile(filepath string, value []byte, cb func (fs.File)) {
	if file, err := FS.Create(filepath); err != nil {
		log.Fatal("Unable to create file ", filepath, " Error: ", err.Error())
	} else {
		defer cb(file)
		file.Write(value)
	}
}



func getFileInfo(file fs.File) fs.FileInfo {
	info, err := file.Stat()
	if err != nil {
		log.Fatal("Could not stat file. Error: ", err.Error())
	}
	return info
}

func init() {
	FSRoot = getConfigDir()
	FS = getAferoBackend(true)
}