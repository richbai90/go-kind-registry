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

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a container bundle",
	Run: actions.Import,
	Args: cobra.MatchAll(cobra.MaximumNArgs(1), func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				return err
			}
		}
		return nil
	}),
}

func init() {

	helpers.InitCommand(func() {
		importCmd.Flags().String("name", "", "Specify a name for the container. Default is the name specified in the bundle")
		importCmd.Flags().String("regPath", "", "Specify a location to bind the container storage to on the host. Default is the path specified in the bundle")
	
		rootCmd.AddCommand(importCmd)
	})

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
