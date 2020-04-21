package core

import (
	"encoding/gob"

	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&User{})

	logrus.WithFields(logrus.Fields{
		"module": "core",
		"action": "init",
	}).Info("Core module initialized")
}
