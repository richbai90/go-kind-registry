package helpers

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Only in test files",
		Short: "Test a command",
		Long:  `This is a generic command. Use it to create a new command for testing.`,
	}

	return cmd
}

func NewTestCommandWithFlags(Flags []*pflag.Flag) *cobra.Command {
	cmd := NewTestCommand()
	for _, flag := range Flags {
		cmd.Flags().AddFlag(flag)
	}

	return cmd
}

// Call init with the provided callback and root command
func MockInit(rootCommand *cobra.Command, init func(*cobra.Command)) {
	init(rootCommand)
}

func PrettyPrint(v interface{}, writer io.Writer) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	b = append(b, '\n')
	if err == nil {
			writer.Write(b)
	}
	return
}

func CreateTestFile() afero.File {
	filename := GenRandomString(5)
	file, err := FS.Create(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	cmd := exec.Command("head", "-c", "100KB", "/dev/urandom")
	cmd.Stdout = file

	if err := cmd.Run(); err != nil {
		log.Fatal(err.Error())
	}

	return file

}

func GenRandomString(len int) string {
	letters := [26]byte {'a','b','c','d','e','f','g','h','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z'}
	str := make([]byte, len)
	
	for i := 0; i < len; i++ {
		rand.Seed(time.Now().UnixNano())
		char := letters[rand.Intn(27)]
		str[i] = char
	}

	return string(str)
}

// Change the file system after it has been intialized. This really only makes sense for testing
func ResetFSBackend() {
	FS = getAferoBackend(true)
}
