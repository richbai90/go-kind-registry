package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type InvalidArgumentError struct {
	Message string
}

func (err InvalidArgumentError) Error() string {
	return fmt.Sprintf("InvalidArgumentError: %s", err.Message)
}

func Run(cmd *exec.Cmd, debug bool, errMsg ...string) {
	if debug {
		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stdout
		PrettyPrint(cmd.Env, os.Stdout)
	}
	err := cmd.Run()
	HandleError(err, errMsg...)
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

func CreateEnvVars(args ...string) ([]string, error) {
	if len(args)%2 > 0 {
		return nil, InvalidArgumentError{
			Message: "Variable provided without a corresponding value. Arguments must come as key value pairs.",
		}
	}

	if len(args) == 0 {
		return []string{}, nil
	}

	vars := make([]string, len(args)/2)
	envVar := strings.Builder{}
	envVar.WriteString(args[0] + "=")
	idx := 0

	for i := 1; i < len(args); i++ {
		envVar.WriteString(args[i])
		if i%2 > 0 {
			vars[idx] = envVar.String()
			envVar.Reset()
		} else {
			idx++
			envVar.WriteRune('=')
		}
	}

	return vars, nil
}

// 0 1 2 3 4 5 6 7
// x 0 x 1 x 2 x 3
