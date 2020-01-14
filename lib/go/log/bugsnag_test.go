package log

import (
	"errors"
	"testing"
	"time"

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

	<-time.After(10 * time.Second)
}

func newError(str string) error {
	return errors.New(str)
}
