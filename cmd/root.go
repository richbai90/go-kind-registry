/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bundle-containers",
	Short: "Bundle all the containers saved in the registry",
	Long: `bundle-containers is a CLI tool designed to work with porter.sh.
	The bundle-containers tool creates a local docker repository where all your containers can be pushed.
	When ready, call "bundle-containers bundle ./<your_container_bundle>.tgz" to get a tar file that can
	Be ported anywhere you need. Restore it by calling bundle-containers start <path_to_bundle.tgz>
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bundle_containers.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Help message for bundle-containers")
	rootCmd.PersistentFlags().BoolP("verbose", "V", false, "Print verbose messages to stdout")
}


