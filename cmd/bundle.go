/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"path"

	"github.com/richbai90/bundle-containers/actions"
	"github.com/spf13/cobra"
)

// bundleCmd represents the bundle command
var bundleCmd = &cobra.Command{
	Use:   "bundle [flags] <Output File Path>",
	Short: "bundle the containers",
	Args: cobra.MatchAll(cobra.ExactArgs(1), func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(path.Dir(args[0])); os.IsNotExist(err) {
			return err
		}
		return nil
	}),
	Run: actions.Bundle,
}

func init() {
	bundleCmd.Flags().String("name", "kr", "The name given to the docker container hosting the registry. Defaults to 'kr'. ")
	rootCmd.AddCommand(bundleCmd)

}
