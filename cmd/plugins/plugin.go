package plugins

import (
	"binary-version-switcher/config"

	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:     "plugins",
	Aliases: []string{"plugin"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "List and Download Plugins",
	Long: `EXAMPLES
	binary-version-switcher plugins get`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return pluginCmd
}
