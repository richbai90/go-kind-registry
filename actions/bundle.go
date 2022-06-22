package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/richbai90/bundle-containers/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func Bundle(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	name := helpers.FlagValue(cmd, "name").String()
	filename := helpers.TrimExtension(args[0])
	bundlepath := path.Join(path.Dir(args[0]), fmt.Sprintf("%s.tar.gz", filename))

	if err := helpers.FS.MkdirAll("/staging", 0760); err != nil {
		log.Fatal("Unable to create staging directory. Error: ", err)
	}
	var config map[string]interface{}
	cfg := helpers.OpenFile("/config.json")
	helpers.CopyFile(cfg, "/config/config.json", func(f1, f2 afero.File) {
		f1.Seek(0, 0)
		config = helpers.ReadJsonFile(f1)
		f1.Close()
		f2.Close()
	})
	if verbose {
		helpers.PrettyPrint(config, os.Stdout)
	}

	osCmd := exec.Command("sh")
	bundleScript := helpers.OpenResource("ociregistrybackup.sh")
	osCmd.Stdin = bundleScript
	envs, _ := helpers.CreateEnvVars(
		"staging_dir", helpers.AbsFilePath("/staging"), 
		"outfile", bundlepath, 
		"reg_version", string(config["version"].(string)), 
		"reg_name", name,
	)
	osCmd.Env = envs
	helpers.Run(osCmd, verbose, "Failed to bundle ", name, ". Error: {ERROR}")

}
