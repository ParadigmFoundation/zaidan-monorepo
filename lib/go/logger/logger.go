package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type LogOpt func(*Logger)

// New returns an initialized logrus.Entry
func New(module string, opts ...LogOpt) *Logger {
	logger := logrus.New()
	if key := os.Getenv("BUGSNAG_APIKEY"); key != "" {
		ConfigureBugsnag(logger, key, module)
	}

	for _, opt := range opts {
		opt(logger)
	}

	logger.SetFormatter(&ModuleFormatter{
		module: module,
	})

	return logger
}

type Logger = logrus.Logger
type Entry = logrus.Entry
type Fields = logrus.Fields
