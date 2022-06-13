package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/richbai90/bundle-containers/helpers"
	"github.com/spf13/cobra"
)

func Bundle(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	name := helpers.FlagValue(cmd, "name").String()
	// Open the container bundle json file that holds all the conf values used to set up the docker registry
	config, err := helpers.FS.Open("/config")
	// Have to handle error directly because CopyToArchive depends on afero.File and helper methods return fs.File for compatibility
	if err != nil {
		log.Fatal("Unable to open version file. Error: ", err)
	}
	
	
	osCmd := exec.Command("docker", "export", helpers.FlagValue(cmd, "name").String(), "-o", path.Join(helpers.FSRoot, path.Base(args[0])))

	helpers.Run(osCmd, fmt.Sprintf("Failed to export docker container %s.", helpers.FlagValue(cmd, "name").String()), verbose)

}
