package common

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var log = logrus.New()

//todo config을 통해 log 설정을 수행한다. (level 등등)
func init() {
	customFormatter := new(prefixed.TextFormatter)
	customFormatter.FullTimestamp = true

	log.Formatter = customFormatter
	log.Level = logrus.DebugLevel
}

func GetLogger(name string) *logrus.Entry{
	return log.WithFields(logrus.Fields{
		"File": name,
	})
}