package utils

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func LogError(err error, message string) {
	log.WithError(err).Error(message)
}

func LogInfo(message string) {
	log.Info(message)
}

func WrapError(err error, message string) error {
	return errors.Wrap(err, message)
}
