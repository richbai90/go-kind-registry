package helpers

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Run(cmd *exec.Cmd, debug bool, errMsg ...string) {
	if debug {
		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stdout
		PrettyPrint(cmd.Env, os.Stdout)
	}
	if err := cmd.Run(); err != nil {
		replacer := strings.NewReplacer("{ERROR}", err.Error())
		for i, s := range(errMsg) {
			errMsg[i] = replacer.Replace(s)
		}
		log.Fatal(errMsg)
	}
}

// Lookup the cobra command flag. If it is not set, use the default pflag.Value
func FlagValue(cmd *cobra.Command, flag string) pflag.Value {
	Flag := cmd.Flag(flag)

	if Flag == nil {
		Flag = &pflag.Flag{}
	}

	return Flag.Value
}

// Wrap the init function in a check for a testing environment
// If the environment is Testing, do not execute the init fn
func InitCommand(init func()) {
	if Testing() {
		return
	}

	init()
}


