package ask

import (
	glog "github.com/maahsome/golang-logger"
	"github.com/sirupsen/logrus"
)

type AskSurvey interface {
	PromptForPath(paths []string, prompt string) string
	PromptForMultipleString(list []string, prompt string) []string
}

type askSurvey struct {
	Logger *logrus.Logger
}

func New(loglevel string) AskSurvey {

	log := glog.CreateStandardLogger()
	log.SetLevel(glog.LogLevelFromString(loglevel))
	return &askSurvey{
		Logger: log,
	}
}
