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

type LogOpt func(*Logger)

// New returns an initialized logrus.Entry
func New(module string, opts ...LogOpt) *Logger {
	logger := logrus.New()
	if key := os.Getenv("BUGSNAG_APIKEY"); key != "" {
		ConfigureBugsnag(logger, key)
	}

	for _, opt := range opts {
		opt(logger)
	}

	logger.SetFormatter(&ModuleFormatter{
		module,
		&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}})

	return logger
}

type Logger = logrus.Logger
type Entry = logrus.Entry
type Fields = logrus.Fields
