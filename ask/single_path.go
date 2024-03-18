package ask

import (
	"os"
	"sort"

	"github.com/AlecAivazis/survey/v2"
)

type (
	PathAnswer struct {
		Path string `survey:"path"`
	}
)

func (as *askSurvey) PromptForPath(paths []string, prompt string) string {

	sort.Strings(paths)

	var pathSurvey = []*survey.Question{
		{
			Name: "path",
			Prompt: &survey.Select{
				Message: prompt,
				Options: paths,
			},
		},
	}

	opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

	pathAnswer := &PathAnswer{}
	if err := survey.Ask(pathSurvey, pathAnswer, opts); err != nil {
		as.Logger.WithError(err).Fatal("No path selected")
	}
	return pathAnswer.Path
}
