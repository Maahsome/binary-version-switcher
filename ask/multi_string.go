package ask

import (
	"os"
	"sort"

	"github.com/AlecAivazis/survey/v2"
)

type (
	StringsAnswer struct {
		Strings []string `survey:"string"`
	}
)

func (as *askSurvey) PromptForMultipleString(list []string, prompt string) []string {

	sort.Strings(list)

	var msSurvey = []*survey.Question{
		{
			Name: "string",
			Prompt: &survey.MultiSelect{
				Message: prompt,
				Options: list,
			},
		},
	}

	opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

	msAnswer := &StringsAnswer{}
	if err := survey.Ask(msSurvey, msAnswer, opts); err != nil {
		as.Logger.WithError(err).Fatal("No strings selected")
	}
	return msAnswer.Strings
}
