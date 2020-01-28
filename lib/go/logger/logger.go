package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type ModuleFormatter struct {
	module string
	baseFormatter *logrus.TextFormatter
}

func (f *ModuleFormatter) Format(entry *Entry) ([]byte, error) {
	arr, err := f.baseFormatter.Format(entry)

	headerBytes := []byte(fmt.Sprintf("%v|", f.module))
	return append(headerBytes, arr...), err
}

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

	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	return entry
}

type Entry = logrus.Entry
