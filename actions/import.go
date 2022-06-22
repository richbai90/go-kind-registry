package actions

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/richbai90/bundle-containers/helpers"
	"github.com/spf13/cobra"
)

// Import the provided bundle
// Changing SHA values from one registry to another makes the import process less predicatable
// For this reason the process is carried out in go without defering to shell scripts
func Import(cmd *cobra.Command, args []string) {
	// Print verbose output during run
	var verbose = false
	// Hold the exec command at each step
	var osCmd *exec.Cmd

	bundlePath := args[0]

	if v, err := cmd.Flags().GetBool("verbose"); err == nil {
		// Set verbose to whatever value the user provided
		verbose = v
	}
	// Create a temporary directory to extract the bundle to
	tmp, err := os.MkdirTemp("/tmp", "bundle")
	helpers.HandleError(err, "Failed to create temporary directory in /tmp. Error: {ERROR}")
	// extract the bundle to the temporary directory
	osCmd = exec.Command("tar", "-zxvf", bundlePath, "--directory", tmp)
	helpers.Run(osCmd, verbose, "Unable to extract bundle. Error: {ERROR}")

	// extract the volumes backup to get the config file
	osCmd = exec.Command("tar", "-zxvf", path.Join(tmp, "backup", "volumes.tar.gz"), "--directory", path.Join(tmp, "backup"))
	helpers.Run(osCmd, verbose, "Failed to extract volumes. Error: {ERROR}")

	// Open the config file for reading
	file, err := os.Open(path.Join(tmp, "backup", "root", "config.json"))
	helpers.HandleError(err, "Failed to open config file. Error: {ERROR}")

	// Parse the config file and save it as config
	cfg := helpers.ReadJsonFile(file)

	config := Config{
		RegName:    string(cfg["name"].(string)),
		RegVersion: string(cfg["version"].(string)),
		RegPort:    int(cfg["port"].(float64)),
		RegPath:    string(cfg["path"].(string)),
	}

	if verbose {
		helpers.PrettyPrint(config, os.Stderr)
	}

	// Use the provided config file if a regPath was not supplied
	if helpers.FlagValue(cmd, "regPath").String() != "" {
		config.RegPath = helpers.FlagValue(cmd, "regPath").String()
	}

	// Use the provided config file if a regPath was not supplied
	if helpers.FlagValue(cmd, "name").String() != "" {
		config.RegName = helpers.FlagValue(cmd, "name").String()
	}

	// Make sure the regPath directory exists
	os.MkdirAll(helpers.FlagValue(cmd, "regPath").String(), 0760)

	// Make sure config directory exists
	os.MkdirAll(path.Join(helpers.FlagValue(cmd, "regPath").String(), "..", "config"), 0760)

	// Move the extracted bundle to the config directory for restoring
	if _, err := helpers.FS.Stat("/restore"); err != fs.ErrNotExist {
		// if we had a failed attempt to restore, delete the restore folder before continuing
		helpers.FS.RemoveAll("/restore")
	}
	if err := os.Rename(path.Join(tmp, "backup"), helpers.AbsFilePath("/restore")); err != nil {
		log.Fatal("Failed to create restore directory. Error: ", err.Error())
	}

	// Import the registry image
	osCmd = exec.Command(
		"docker",
		"load",
		"-i",
		helpers.AbsFilePath("/restore/registry.tar"),
	)
	
	helpers.Run(osCmd, verbose, "Failed to import registry container. Error: {ERROR}")
	
	// Import the alpine image
	osCmd = exec.Command(
		"docker",
		"load",
		"-i",
		helpers.AbsFilePath("/restore/alpine.tar"),
	)
	helpers.Run(osCmd, verbose, "Failed to import registry container. Error: {ERROR}")

	// Import the kind image
	osCmd = exec.Command(
		"docker",
		"load",
		"-i",
		helpers.AbsFilePath("/restore/kind.tar"),
	)

	// We need to capture the output so cannot use helper method Run
	out, err := osCmd.Output()
	helpers.HandleError(err, "Failed to restore kind image. Error: {ERROR}")
	out = out[:len(out)-1]
	// Check for verbose option (helpers.Run does this automtaically)
	if verbose {
		os.Stdout.Write(out)
	}

	// Extract the new SHA
	sha := strings.TrimPrefix(string(out), "Loaded image ID: ")
	// Tag the kind image - must leave out SHA during tag
	osCmd = exec.Command("docker", "tag", sha, fmt.Sprintf("kindest/node:%s", strings.Split(config.RegVersion, "@")[0]))
	helpers.Run(osCmd, verbose, 
	`Failed to tag kind image.
	Command: `, osCmd.String(), `
	Error: {ERROR}`)

	osCmd = exec.Command("sh")
	env, _ := helpers.CreateEnvVars(
		"reg_path", config.RegPath,
		"reg_version", config.RegVersion,
		"reg_name", config.RegName,
		"restore_dir", helpers.AbsFilePath("/restore"),
	)

	osCmd.Env = env
	osCmd.Stdin = helpers.OpenResource("ociregistryrestore.sh")
	helpers.Run(osCmd, verbose, "Failed to import registry volumes. Error: {ERROR}")
}
