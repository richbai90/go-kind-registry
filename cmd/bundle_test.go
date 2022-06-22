package cmd

import (
	"os"
	"testing"

	"github.com/richbai90/bundle-containers/helpers"
	"github.com/spf13/cobra"
)

func TestBundleCmd(t *testing.T) {
	cmd := helpers.NewTestCommand()
	helpers.MockInit(cmd, func(c *cobra.Command) {
		// bundleCmd.Flags().String("name", "kr", "The name given to the docker container hosting the registry. Defaults to 'kr'. ")
		bundleCmd.PreRun = func(cmd *cobra.Command, args []string) {
			os.Setenv("GO_ENVIRON", "PROD")
			helpers.ResetFSBackend()
			helpers.FS.Create("/config")
		}
		c.AddCommand(bundleCmd)
	})

	// call the start command
	cmd.SetArgs([]string{
		"bundle",
		"/tmp/bundle_test.tar",
	})

	cmd.Execute()
}
