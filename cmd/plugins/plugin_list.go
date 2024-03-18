package plugins

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := getPluginList(c.VersionDetail.SemVer)
		if err != nil {
			c.Logger.WithError(err).Error("failed to fetch a list of plugins")
		}
		for _, v := range list {
			fmt.Println(v)
		}
	},
}

func init() {
	pluginCmd.AddCommand(listCmd)
}
