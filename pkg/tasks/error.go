package tasks

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func logTaskError(taskName string, err error) {
	logrus.WithFields(logrus.Fields{
		"tasks":       taskName,
		"stack_trace": fmt.Sprintf("%+v", err),
	}).WithError(err).Error("Error Task")
}
