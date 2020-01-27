package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type LogOpt func(*logrus.Entry)

// New returns an initialized logrus.Entry
func New(module string, opts ...LogOpt) *Entry {
	log := logrus.New()
	if key := os.Getenv("BUGSNAG_APIKEY"); key != "" {
		ConfigureBugsnag(log, key)
	}

	entry := log.WithField("module", module)

	for _, opt := range opts {
		opt(entry)
	}

	return entry
}

type Entry = logrus.Entry
