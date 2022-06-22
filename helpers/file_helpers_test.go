package helpers

import (
	"testing"

	"github.com/mholt/archiver/v4"
	"github.com/spf13/afero"
)

func TestMakeJsonFile(t *testing.T) {
	JSON := struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}{
		Foo: "bar",
		Bar: "baz",
	}

	filepath := "/test.json"

	type args struct {
		filepath string
		JSON     interface{}
		cb       func(afero.File)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "File exists",
			args: args{
				filepath: filepath,
				JSON:     JSON,
				cb: func(f afero.File) {

					if _, err := f.Stat(); err != nil {
						t.Errorf("Expected file to be stattable. Got error %s", err.Error())
					}
				},
			},
		},
		{
			name: "File is correct size",
			args: args{
				filepath: filepath,
				JSON:     JSON,
				cb: func(f afero.File) {
					stats, _ := f.Stat()
					if stats.Size() < 22 {
						t.Errorf("Expected file size of 22. Got %d", stats.Size())
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MakeJsonFile(tt.args.filepath, tt.args.JSON, tt.args.cb)
		})
	}
}

func createTar() afero.File {
	file, err := FS.Create("/test.tar")
	if err != nil {
		panic(err)
	}

	return file
}

func Test_copyToTarArchive(t *testing.T) {
	type args struct {
		file    afero.File
		archive afero.File
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Copy To Archive",
			args: args{
				file:    CreateTestFile(),
				archive: createTar(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, _, err := archiver.Identify(tt.args.archive.Name(), tt.args.archive)
			if err != nil {
				t.Error(err.Error())
			}
			if _, ok := format.(archiver.Extractor); !ok {
				t.Error("Unable to open archive")
			}
		})
	}
}
