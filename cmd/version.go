package cmd

import (
	"encoding/json"

	"binary-version-switcher/objects"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Express the 'version' of binary-version-switcher.",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := expressVersion()
		if err != nil {
			logrus.WithError(err).Error("Failed to express the version")
		}
		if !c.FormatOverridden {
			c.OutputFormat = "json"
		}
		c.OutputData(&version)
	},
}

func expressVersion() (objects.Version, error) {
	var verData objects.Version
	err := json.Unmarshal([]byte(c.VersionJSON), &verData)
	if err != nil {
		return verData, errors.Wrap(err, "Failed to unmarshal JSON")
	}

	return verData, nil
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
