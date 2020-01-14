package log

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestBugSnag(t *testing.T) {
	err := newError("error from logrus using hook")
	logrus.Error("Raw error")
	logrus.WithFields(logrus.Fields{
		"foo":    "bar",
		"ticker": "BTC/USD",
		"number": 1,
	}).WithError(err).Error(err)
	select {}
}

func newError(str string) error {
	return errors.New(str)
}
