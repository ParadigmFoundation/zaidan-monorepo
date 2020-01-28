package logger

import (
	"fmt"

	ethlog "github.com/ethereum/go-ethereum/log"
	"github.com/sirupsen/logrus"
)

type EthLogHandler struct {
	*logrus.Logger
}

func (log EthLogHandler) Log(record *ethlog.Record) error {
	entry := &logrus.Entry{Logger: log.Logger}
	for i := 0; i < len(record.Ctx); i += 2 {
		key, val := record.Ctx[i].(string), record.Ctx[i+1]
		if key == "err" {
			var err error
			switch t := val.(type) {
			case error:
				err = t
			default:
				err = fmt.Errorf("%v", val)
			}
			entry = entry.WithError(err)
		} else {
			entry = entry.WithField(key, val)
		}
	}

	var logFn func(...interface{})
	switch record.Lvl {
	case ethlog.LvlCrit:
		logFn = entry.Fatal
	case ethlog.LvlError:
		logFn = entry.Error
	case ethlog.LvlWarn:
		logFn = entry.Warn
	case ethlog.LvlInfo:
		logFn = entry.Info
	case ethlog.LvlDebug:
		logFn = entry.Debug
	case ethlog.LvlTrace:
		logFn = entry.Trace
	}
	logFn(record.Msg)
	return nil
}

func HandleEthLog() LogOpt {
	fn := func(logger *logrus.Logger) {
		ethlog.Root().SetHandler(EthLogHandler{logger})
	}
	return fn
}
