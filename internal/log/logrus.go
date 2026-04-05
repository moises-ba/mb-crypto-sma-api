package log

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Debug(message string) {
	log.Debug(message)
}

func Info(message string) {
	log.Info(message)
}

func Warn(message string) {
	log.Warn(message)
}

func Error(message string, err error) {
	log.Errorf(message+": \n%+v", err)
}
