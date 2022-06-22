package cmd

import (
	"testing"

	"github.com/richbai90/bundle-containers/helpers"
	"github.com/spf13/cobra"
)

func TestStartInit(t *testing.T) {
	cmd := helpers.NewTestCommand()
	helpers.MockInit(cmd, func(c *cobra.Command) {
		startCmd.Flags().String("name", "", "Specify a name for the container. Default is 'kr'")
		startCmd.Flags().StringP("port", "p", "", "Specify a port to bind on the host to access the OCI. Default is 5015")
		startCmd.Flags().String("regPath", "", "Specify a location to bind the container storage to on the host. Default is $HOME/.config/kind/registry")
		startCmd.Flags().String("version", "", "Specify a kind image to use. Default is the current latest version")
		c.AddCommand(startCmd)
	})

	// call the start command
	cmd.SetArgs([]string{
		"start",
	})

	cmd.Execute()
}
