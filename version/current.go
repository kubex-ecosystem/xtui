package version

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
	"time"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of XTui",
		Long:  "Print the version number of XTui",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(GetVersionInfo())
		},
	}
	subLatestCmd = &cobra.Command{
		Use:   "latest",
		Short: "Print the latest version number of XTui",
		Long:  "Print the latest version number of XTui",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(GetLatestVersionInfo())
		},
	}
	subCmdCheck = &cobra.Command{
		Use:   "check",
		Short: "Check if the current version is the latest version of XTui",
		Long:  "Check if the current version is the latest version of XTui",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(GetVersionInfoWithLatestAndCheck())
		},
	}
)

const gitModelUrl = "https://github.com/faelmori/xtui.git"
const currentVersionFallback = "v1.1.3" // First version with the version file

//go:embed CLI_VERSION
var currentVersion string

func GetVersion() string {
	if currentVersion == "" {
		return currentVersionFallback
	}
	return currentVersion
}

func GetGitModelUrl() string {
	return gitModelUrl
}

func GetVersionInfo() string {
	return "Version: " + GetVersion() + "\n" + "Git repository: " + GetGitModelUrl()
}

func GetLatestVersionFromGit() string {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	response, err := netClient.Get(gitModelUrl + "/releases/latest")
	if err != nil {
		return "Error: " + err.Error()
	}

	if response.StatusCode != 200 {
		return "Error: " + response.Status
	}

	tag := strings.Split(response.Request.URL.Path, "/")

	return tag[len(tag)-1]
}

func GetLatestVersionInfo() string {
	return "Latest version: " + GetLatestVersionFromGit()
}

func GetVersionInfoWithLatestAndCheck() string {
	if GetVersion() == GetLatestVersionFromGit() {
		return GetVersionInfo() + "\n" + GetLatestVersionInfo() + "\n" + "You are using the latest version."
	} else {
		return GetVersionInfo() + "\n" + GetLatestVersionInfo() + "\n" + "You are using an outdated version.\n" + "Please upgrade your XTui to prevent any issues."
	}
}

func CliCommand() *cobra.Command {
	versionCmd.AddCommand(subLatestCmd)
	versionCmd.AddCommand(subCmdCheck)
	return versionCmd
}
