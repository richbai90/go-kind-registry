/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/richbai90/bundle-containers/actions"
	"github.com/richbai90/bundle-containers/helpers"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [flags] [ path/to/bundle.tar ]",
	Short: "Start the registry",
	Long: `Start the registry at the specified path. 
	If a positional argument is specified, and it is a valid path to an exported docker container, 
	that will be used instead of creating a new one.`,
	Args: cobra.MatchAll(cobra.MaximumNArgs(1), func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				return err
			}
		}
		return nil
	}),
	Run: actions.Start,
}

func init() {

	helpers.InitCommand(func() {
		startCmd.Flags().String("name", "", "Specify a name for the container. Default is 'kind-registry'")
		startCmd.Flags().IntP("port", "p", 5015, "Specify a port to bind on the host to access the OCI. Default is 5015")
		startCmd.Flags().String("regPath", helpers.FSRoot, "Specify a location to bind the container storage to on the host. Default is $HOME/.config/kind/registry")
		startCmd.Flags().String("version", "", "Specify a kind image to use. Default is the current latest version")
		rootCmd.AddCommand(startCmd)
	})

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
