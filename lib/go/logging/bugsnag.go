package logging

import (
	bs "github.com/bugsnag/bugsnag-go"
	"log"
	"os"
	"strings"
)

var errorLog = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
var infoLog = log.New(os.Stdout, "[INFO] ", log.LstdFlags)

func ConfigureBugsnag (apiKey string) {
	wd, _ := os.Getwd()
	wd = strings.SplitAfter(wd, "zaidan-monorepo")[0]
	bs.Configure(bs.Configuration{
		APIKey:          apiKey,
		SourceRoot: 	 wd,
		ProjectPackages: []string{"github.com/ParadigmFoundation/zaidan-monorepo/**"},
		Logger:          errorLog,
	})
}

func Fatal(err error, bugsnagRawData ...interface{}) {
	SafeError(err, bugsnagRawData...)
	os.Exit(1)
}

func SafeError (err error, bugsnagRawData ...interface{}) {
	if len(bs.Config.APIKey) == 0 {
		errorLog.Println(err)
	} else {
		err := bs.Notify(err, bugsnagRawData...)
		if err != nil {
			errorLog.Println(err)
		}
	}
}

func Info(v ...interface{}) {
	infoLog.Println(v...)
}
