package log

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

func ConfigureBugsnag(key string) {
	wd, _ := os.Getwd()
	wd = strings.SplitAfter(wd, "zaidan-monorepo")[0]
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:          key,
		SourceRoot:      wd,
		Synchronous:     false,
		ProjectPackages: []string{"github.com/ParadigmFoundation/zaidan-monorepo/**"},
		// Logger is disabled cause logrus will take care of it
		Logger: &noop{},
	})
	bugsnag.Config.APIKey = key

	bsOnce.Do(func() {
		logrus.AddHook(&BugsnagHook{})
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	})
}

type noop struct{}

func (*noop) Printf(string, ...interface{}) {}

func init() {
	if key := os.Getenv("BUGSNAG_APIKEY"); key != "" {
		ConfigureBugsnag(key)
	}
}
