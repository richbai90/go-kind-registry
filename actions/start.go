package actions

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/richbai90/bundle-containers/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type Config struct {
	RegName    string `json:"name"`
	RegPort    int    `json:"port"`
	RegPath    string `json:"path"`
	RegVersion string `json:"version"`
}

func Start(cmd *cobra.Command, args []string) {
	// Print verbose output during run
	var verbose = false
	// Hold the exec command at each step
	var osCmd *exec.Cmd
	// Track if starting from import
	fmt.Println("Starting, this may take a few moments...")
	if v, err := cmd.Flags().GetBool("verbose"); err == nil {
		// Set verbose to whatever value the user provided
		verbose = v
	}

	if file, err := helpers.FS.Open("/config.json"); err == nil {
		// An import has been run so use the configs from the import instead
		config := helpers.ReadJsonFile(file)
		file.Close()
		// prefer config values except for port which is safe to change between envs
		helpers.FlagValue(cmd, "version").Set(string(config["version"].(string)))
		helpers.FlagValue(cmd, "regPath").Set(string(config["path"].(string)))
		helpers.FlagValue(cmd, "name").Set(string(config["name"].(string)))
		// We have collected all the useful information from the import, remove the config file now to prevent errors later
		helpers.FS.Remove("/config.json")
	} else {
		// first time running. Presumably in dev env.
		cmd.Flags().Set("version", GetKind(helpers.FlagValue(cmd, "version").String()))
	}

	// Get the installer script embedded with the executable
	regInstaller := helpers.OpenResource("ociregistry.sh")
	// close the file when we're done
	defer regInstaller.Close()

	// On mac we have to install a package since we cannot guarantee brew is installed
	// On linux kind is a self contained executable named kind
	// TODO: Give windows some love too
	// TODO: Dynamically update the kind exe as well as kind image -- Right now kind is embedded in the binary
	filenameMap := map[string]string{
		"linux":  "kind",
		"darwin": "kind-0.14.0.pkg",
	}

	// Get the correct kind file for the running OS
	kind, filename := helpers.OpenWhenOS(filenameMap)
	kindPath := path.Join(helpers.FSRoot, filename)

	// Copy the kind file to the config path to be executed by the OS
	// All fs helpers execute from the $HOME/.config/kind root
	helpers.CopyFile(kind, fmt.Sprintf("/%s", filename), func(inFile, outFile fs.File) {
		// close the resource file
		inFile.Close()
		// close the copied file
		outFile.Close()
	})

	// Collect all the config values used to run this command in an encoding/json compatible struct
	config := Config{
		RegName:    helpers.FlagValue(cmd, "name").String(),
		RegPort:    helpers.Atoi(helpers.FlagValue(cmd, "port").String()),
		RegPath:    helpers.FlagValue(cmd, "regPath").String(),
		RegVersion: helpers.FlagValue(cmd, "version").String(),
	}

	// Apply the appropriate env vars to influence the install script based on command flags
	env, _ := helpers.CreateEnvVars(
		"reg_name", config.RegName,
		"reg_port", fmt.Sprintf("%d", config.RegPort),
		"reg_path", config.RegPath,
		"reg_version", config.RegVersion,
		"PATH", fmt.Sprintf("%s:%s", os.Getenv("PATH"), path.Dir(kindPath)),
		"BUNDLE_DIR", helpers.FSRoot,
	)

	// save the configuration for bundle and cleanup processes
	helpers.MakeJsonFile("/config.json", config, func(f afero.File) {
		f.Close()
		// Make sure the config file is readable by users other than owner.
		// This is a necessetiy when running on macos as file gets owned by root
		os.Chmod(helpers.AbsFilePath("/config.json"), 666)
	})

	if strings.HasSuffix(filename, ".pkg") {
		// If the kind filename returned is a .pkg then we are running on a mac
		// We need to install it before continuing
		osCmd = exec.Command("/usr/sbin/installer", "-pkg", kindPath, "-target", "/")
		ensurePermissions()
		helpers.Run(osCmd, verbose, "Unable to install the kind package. Error: {ERROR}")
	}

	// Defer to the OS shell for the install of kind
	osCmd = exec.Command("sh")
	// Read commands from the embedded install script
	osCmd.Stdin = regInstaller
	// Update the environment with the appropriate vars
	osCmd.Env = env

	// Run the installer
	helpers.Run(osCmd, verbose, "Failed to execute installer. Error: {ERROR}")

}

// Track at which step the last  prompt was given
// This way we can avoid displaying the prompt excessively
type prompt struct {
	step  int
	value any
}

func GetKind(version string) string {
	// Hold the version found from the github API query
	// This value begins as the version number or "latest"
	// It will end holding the full version string that number represents
	var Version string = "latest"
	if version != "" {
		Version = version
	}

	var Prompt prompt
	// Hold the decoded github API response as a map
	var release map[string]interface{}
	// url used to lookup the specified version metadata from the github API
	url := fmt.Sprintf("https://api.github.com/repos/kubernetes-sigs/kind/releases/%s", Version)
	// Begin by looking up the latest version using the github API
	// If the response is an error skip the rest of the lookup process
	if resp, err := http.Get(url); err == nil {
		// decode the server response and store it in the release variable
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&release); err != nil {
			// If there was an error decoding the response, warn the user with a prompt to continue
			if err := helpers.WarnWithPrompt(
				err,
				"Unable to parse JSON response from server.",
				"Proceed using default version as a fallback? [y/Y/n/N]",
				// callback fn gets called after response is recieved
				func(resp string) {
					if !helpers.IsIn(resp, []string{"y", "Y"}) {
						// If the response from the user was not y/Y than exit with a failure
						log.Fatal("Unable to parse JSON response from server. Error: ", err.Error())
					} else {
						// Otherwise mark this as procedure 1 and continue
						Prompt = prompt{step: 1}
					}
				},
			); err != nil {
				// There was an invalid response from the user. Stop with a fatal error
				log.Fatal("Did not understand response. Error was: ", err.Error())
			}
		}
		if Prompt.step != 1 {
			// Skip further processing if we already failed to query the version, and the user has decided to proceed
			// Body holds the entire release description. We use a regex query to find the version information
			body := release["body"].(string)
			// The image version is of the form kindest/node:vd.d.d@sha256:<sha>
			r := regexp.MustCompile("`kindest\\/node:v(\\d+\\.\\d+\\.\\d+@sha256:\\w+)")
			matches := r.FindStringSubmatch(body)
			if matches != nil && len(matches) > 1 {
				// if the version string is found, get just the version number and sha
				Version = matches[1]
			} else {
				// Otherwise use default
				// TODO: prompt the user here also?
				log.Print("Could not find version number from response. Perhaps the format has changed? Using default value as fallback")
			}
		}
	} else {
		helpers.WarnWithPrompt(err, "Failed to query the github API. Did you mean to run import first? ", "yY/nN", func(resp string) {
			if strings.ToLower(resp) == "y" {
				log.Fatal("Run an import with bundle-containers import and then rerun start. See import -h for more info.")
			} else {
				print("Proceeding with default version information.\n")
			}
		})
	}

	if Version == "latest" {
		// Failed to update the version somewhere along the way. Use fallback
		// TODO: This is the latest version at time of writing. Could be useful to let the user supply a fallback instead
		Version = "1.24.0@sha256:0866296e693efe1fed79d5e6c7af8df71fc73ae45e3679af05342239cdc5bc8e"
	}

	// Pull the kind image
	cmd := exec.Command("docker", "pull", fmt.Sprintf("kindest/node:v%s", Version))
	if err := cmd.Run(); err != nil {
		// If the pull failed exit with message to user
		log.Fatal("Unable to pull docker from repo")
	}

	return Version
}

func ensurePermissions() {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	helpers.HandleError(err, "Failed to determine user. Error: {ERROR}")

	// output has trailing \n
	// need to remove the \n
	// otherwise it will cause error for strconv.Atoi
	i := helpers.Atoi(string(output[:len(output)-1]))

	if i != 0 {
		// 0 = root
		log.Fatal("This program must be run with root (sudo) permissions.")
	}
}
