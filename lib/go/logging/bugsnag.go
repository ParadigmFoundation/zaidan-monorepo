package logging

import (
	"fmt"
	bs "github.com/bugsnag/bugsnag-go"
	"log"
	"os"
)

var errorLog = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
var infoLog = log.New(os.Stdout, "[INFO] ", log.LstdFlags)

func ConfigureBugsnag (apiKey string) {
	bs.Configure(bs.Configuration{
		APIKey:          apiKey,
		// The import paths for the Go packages containing your source files
		ProjectPackages: []string{"github.com/ParadigmFoundation/zaidan-monorepo"},
		Logger:          errorLog,
	})
}

func FatalString(err string) {
	Fatal(fmt.Errorf(err))
}

func Fatal(err error) {
	SafeError(err)
	os.Exit(1)
}

func SafeErrorString(err string) {
	SafeError(fmt.Errorf(err))
}

func SafeError (err error) {
	if len(bs.Config.APIKey) == 0 {
		errorLog.Println(err)
	} else {
		err := bs.Notify(err)
		if err != nil {
			errorLog.Println(err)
		}
	}
}

func Info(v ...interface{}) {
	infoLog.Println(v...)
}
