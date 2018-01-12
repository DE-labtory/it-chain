package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

var ilogger = InitLogger()

func InitLogger() *logrus.Logger{
//todo config을 통해 log 설정을 수행한다. (level 등등)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	return logrus.New()
}


func GetLogger(name string) *logrus.Entry{
	return ilogger.WithFields(logrus.Fields{
		"File": name,
	})
}