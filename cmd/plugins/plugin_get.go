package plugins

import (
	"binary-version-switcher/ask"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// getCmd represents the list command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		selectPlugins()
	},
}

func selectPlugins() {

	a := ask.New(c.LogLevel)
	pList, perr := getPluginList(c.VersionDetail.SemVer)
	if perr != nil {
		c.Logger.WithError(perr).Error("failed to get the list of plugins")
	}
	selectedPlugins := a.PromptForMultipleString(pList, "Select the plugins to download")
	for _, name := range selectedPlugins {
		downloadPlugin(c.VersionDetail.SemVer, name, c.PluginPath)
	}

}

func downloadPlugin(semver string, name string, dst string) error {
	// TODO: eventually build and populate Release Artifacts with OS/ARCH in the names
	// targetOS := runtime.GOOS
	// if targetOS != "linux" && targetOS != "darwin" {
	// 	return fmt.Errorf("unsupported OS: %s", targetOS)
	// }

	// targetArch := runtime.GOARCH
	// if targetArch == "amd64" {
	// 	targetArch = "amd64"
	// } else if targetArch == "arm64" {
	// 	targetArch = "arm64"
	// } else {
	// 	return fmt.Errorf("unsupported architecture: %s", targetArch)
	// }

	url := fmt.Sprintf("https://github.com/Maahsome/binary-version-switcher-plugins/releases/download/%s/bvs-%s.so", semver, name)

	destFile := filepath.Join(dst, fmt.Sprintf("bvs-%s.so", name))
	c.Logger.Infof("Downloading Plugin to %s", destFile)
	client := resty.New()
	resp, err := client.R().SetOutput(destFile).Get(url)
	if err != nil {
		return fmt.Errorf("error downloading kubectl: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("error downloading kubectl: status code %d", resp.StatusCode())
	}

	return nil
}

func getPluginList(semver string) ([]string, error) {

	client := resty.New()
	url := fmt.Sprintf("https://github.com/Maahsome/binary-version-switcher-plugins/releases/download/%s/plugins.txt", semver)
	resp, err := client.R().
		Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	tmpNames := strings.Split(string(resp.Body()), "\n")
	tagNames := []string{}
	for _, v := range tmpNames {
		if len(v) > 0 {
			tagNames = append(tagNames, v)
		}
	}
	return tagNames, nil

}

func init() {
	pluginCmd.AddCommand(getCmd)
}
