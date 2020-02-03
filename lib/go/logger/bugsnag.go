package logger

import (
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/bugsnag/bugsnag-go"
	bugsnag_errors "github.com/bugsnag/bugsnag-go/errors"
	"github.com/sirupsen/logrus"
)

type BugsnagHook struct{}

func (*BugsnagHook) Fire(entry *logrus.Entry) error {
	var notifyErr error
	// .Data["error"] is populated when .WithError() from logrus is called
	err, ok := entry.Data["error"].(error)
	if ok {
		notifyErr = err
	} else {
		notifyErr = errors.New(entry.Message)
	}

	// Build Metadata an save it as `extra` in bugsnag
	md := make(bugsnag.MetaData)
	md["extra"] = make(map[string]interface{})
	for key, val := range entry.Data {
		if key == "error" {
			continue
		}
		md["extra"][key] = val
	}

	// In order to get a stack trace right were our function fired, we need to skip some frames from the backtrace
	// We need to skip certain number of frames depending on the way we called logrus
	// for logrus.WithFields() we'll skip 6, for logrus.<LogLevel>() we'll skip one extra frame
	skip := 6
	if len(md["extra"]) == 0 {
		skip++
	}
	errWithStack := bugsnag_errors.New(notifyErr, skip)

	// Build extra .Notify() args
	var args = []interface{}{md}
	if ctx := entry.Context; ctx != nil {
		args = append(args, bugsnag.StartSession(ctx))
	}

	return bugsnag.Notify(errWithStack, args...)
}

func (*BugsnagHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel,
	}
}

var bsOnce sync.Once

func ConfigureBugsnag(log *logrus.Logger, key string, module string) {
	wd, _ := os.Getwd()
	wd = strings.SplitAfter(wd, "zaidan-monorepo")[0]
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:          key,
		SourceRoot:      wd,
		Synchronous:     false,
		ProjectPackages: []string{"github.com/ParadigmFoundation/zaidan-monorepo/**"},
		AppType: module,
		// Logger is disabled cause logrus will take care of it
		Logger: &noop{},
	})

	log.AddHook(&BugsnagHook{})
}

type noop struct{}

func (*noop) Printf(string, ...interface{}) {}
