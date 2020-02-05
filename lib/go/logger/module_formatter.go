package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37

	defaultTimestampFormat = time.RFC3339
	FieldKeyMsg            = "msg"
	FieldKeyLevel          = "level"
	FieldKeyTime           = "time"
	FieldKeyLogrusError    = "logrus_error"
	FieldKeyFunc           = "func"
	FieldKeyFile           = "file"

	PanicLevel logrus.Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

var baseTimestamp time.Time

type fieldKey string
type FieldMap map[fieldKey]string
func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}

type ModuleFormatter struct {
	logrus.TextFormatter
	BugsnagHook
	module string
	isTerminal bool
	terminalInitOnce sync.Once
	FieldMap
}

func (f *ModuleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(Fields)
	for k, v := range entry.Data {
		data[k] = v
	}
	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	var funcVal, fileVal string

	fixedKeys := make([]string, 0)
	fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyTime))
	fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyLevel))
	if entry.Message != "" {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(FieldKeyMsg))
	}
	if entry.HasCaller() {
		fixedKeys = append(fixedKeys,
			f.FieldMap.resolve(FieldKeyFunc), f.FieldMap.resolve(FieldKeyFile))
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		} else {
			funcVal = entry.Caller.Function
			fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		}
	}

	if !f.DisableSorting {
		if f.SortingFunc == nil {
			sort.Strings(keys)
			fixedKeys = append(fixedKeys, keys...)
		} else {
			if !f.isColored() {
				fixedKeys = append(fixedKeys, keys...)
				f.SortingFunc(fixedKeys)
			} else {
				f.SortingFunc(keys)
			}
		}
	} else {
		fixedKeys = append(fixedKeys, keys...)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.terminalInitOnce.Do(func() { f.init(entry) })

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if f.isColored() {
		f.printColored(b, entry, keys, data, timestampFormat)
	} else {
		for _, key := range fixedKeys {
			var value interface{}
			switch {
			case key == f.FieldMap.resolve(FieldKeyTime):
				value = entry.Time.Format(timestampFormat)
			case key == f.FieldMap.resolve(FieldKeyLevel):
				value = entry.Level.String()
			case key == f.FieldMap.resolve(FieldKeyMsg):
				value = entry.Message
			//case key == f.FieldMap.resolve(FieldKeyLogrusError):
			//	value = entry.err
			case key == f.FieldMap.resolve(FieldKeyFunc) && entry.HasCaller():
				value = funcVal
			case key == f.FieldMap.resolve(FieldKeyFile) && entry.HasCaller():
				value = fileVal
			default:
				value = data[key]
			}
			f.appendKeyValue(b, key, value)
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

//LOGRUS BELOW

//var baseTimestamp time.Time
//
//func init() {
//	baseTimestamp = time.Now()
//}
//
//// TextFormatter formats logs into text
//type TextFormatter struct {
//	// Set to true to bypass checking for a TTY before outputting colors.
//	ForceColors bool
//
//	// Force disabling colors.
//	DisableColors bool
//
//	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
//	EnvironmentOverrideColors bool
//
//	// Disable timestamp logging. useful when output is redirected to logging
//	// system that already adds timestamps.
//	DisableTimestamp bool
//
//	// Enable logging the full timestamp when a TTY is attached instead of just
//	// the time passed since beginning of execution.
//	FullTimestamp bool
//
//	// TimestampFormat to use for display when a full timestamp is printed
//	TimestampFormat string
//
//	// The fields are sorted by default for a consistent output. For applications
//	// that log extremely frequently and don't use the JSON formatter this may not
//	// be desired.
//	DisableSorting bool
//
//	// The keys sorting function, when uninitialized it uses sort.Strings.
//	SortingFunc func([]string)
//
//	// Disables the truncation of the level text to 4 characters.
//	DisableLevelTruncation bool
//
//	// QuoteEmptyFields will wrap empty fields in quotes if true
//	QuoteEmptyFields bool
//
//	// Whether the logger's out is to a terminal
//	isTerminal bool
//
//	// FieldMap allows users to customize the names of keys for default fields.
//	// As an example:
//	// formatter := &TextFormatter{
//	//     FieldMap: FieldMap{
//	//         FieldKeyTime:  "@timestamp",
//	//         FieldKeyLevel: "@level",
//	//         FieldKeyMsg:   "@message"}}
//	FieldMap FieldMap
//
//	// CallerPrettyfier can be set by the user to modify the content
//	// of the function and file keys in the json data when ReportCaller is
//	// activated. If any of the returned value is the empty string the
//	// corresponding key will be removed from json fields.
//	CallerPrettyfier func(*runtime.Frame) (function string, file string)
//
//	terminalInitOnce sync.Once
//}
//
//func (f *ModuleFormatter) init(entry *Entry) {
//	if entry.Logger != nil {
//		f.isTerminal = checkIfTerminal(entry.Logger.Out)
//
//		if f.isTerminal {
//			initTerminal(entry.Logger.Out)
//		}
//	}
//}
//
//func (f *ModuleFormatter) isColored() bool {
//	isColored := f.ForceColors || (f.isTerminal && (runtime.GOOS != "windows"))
//
//	if f.EnvironmentOverrideColors {
//		if force, ok := os.LookupEnv("CLICOLOR_FORCE"); ok && force != "0" {
//			isColored = true
//		} else if ok && force == "0" {
//			isColored = false
//		} else if os.Getenv("CLICOLOR") == "0" {
//			isColored = false
//		}
//	}
//
//	return isColored && !f.DisableColors
//}
//
//// Format renders a single log entry
//
//
//func (f *ModuleFormatter) printColored(b *bytes.Buffer, entry *Entry, keys []string, data Fields, timestampFormat string) {
//	var levelColor int
//	switch entry.Level {
//	case DebugLevel, TraceLevel:
//		levelColor = gray
//	case WarnLevel:
//		levelColor = yellow
//	case ErrorLevel, FatalLevel, PanicLevel:
//		levelColor = red
//	default:
//		levelColor = blue
//	}
//
//	levelText := strings.ToUpper(entry.Level.String())
//	if !f.DisableLevelTruncation {
//		levelText = levelText[0:4]
//	}
//
//	// Remove a single newline if it already exists in the message to keep
//	// the behavior of logrus text_formatter the same as the stdlib log package
//	entry.Message = strings.TrimSuffix(entry.Message, "\n")
//
//	caller := ""
//
//	if entry.HasCaller() {
//		funcVal := fmt.Sprintf("%s()", entry.Caller.Function)
//		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
//
//		if f.CallerPrettyfier != nil {
//			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
//		}
//		caller = fileVal + " " + funcVal
//	}
//
//	if f.DisableTimestamp {
//		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m%s %-44s ", levelColor, levelText, caller, entry.Message)
//	} else if !f.FullTimestamp {
//		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d]%s %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), caller, entry.Message)
//	} else {
//		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s]%s %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), caller, entry.Message)
//	}
//	for _, k := range keys {
//		v := data[k]
//		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
//		f.appendValue(b, v)
//	}
//}
//
func (f *ModuleFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *ModuleFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *ModuleFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

func prefixFieldClashes(data Fields, fieldMap FieldMap, reportCaller bool) {
	timeKey := fieldMap.resolve(FieldKeyTime)
	if t, ok := data[timeKey]; ok {
		data["fields."+timeKey] = t
		delete(data, timeKey)
	}

	msgKey := fieldMap.resolve(FieldKeyMsg)
	if m, ok := data[msgKey]; ok {
		data["fields."+msgKey] = m
		delete(data, msgKey)
	}

	levelKey := fieldMap.resolve(FieldKeyLevel)
	if l, ok := data[levelKey]; ok {
		data["fields."+levelKey] = l
		delete(data, levelKey)
	}

	logrusErrKey := fieldMap.resolve(FieldKeyLogrusError)
	if l, ok := data[logrusErrKey]; ok {
		data["fields."+logrusErrKey] = l
		delete(data, logrusErrKey)
	}

	// If reportCaller is not set, 'func' will not conflict.
	if reportCaller {
		funcKey := fieldMap.resolve(FieldKeyFunc)
		if l, ok := data[funcKey]; ok {
			data["fields."+funcKey] = l
		}
		fileKey := fieldMap.resolve(FieldKeyFile)
		if l, ok := data[fileKey]; ok {
			data["fields."+fileKey] = l
		}
	}
}

func (f *ModuleFormatter) isColored() bool {
	isColored := f.ForceColors || (f.isTerminal && (runtime.GOOS != "windows"))

	if f.EnvironmentOverrideColors {
		if force, ok := os.LookupEnv("CLICOLOR_FORCE"); ok && force != "0" {
			isColored = true
		} else if ok && force == "0" {
			isColored = false
		} else if os.Getenv("CLICOLOR") == "0" {
			isColored = false
		}
	}

	return isColored && !f.DisableColors
}

func (f *ModuleFormatter) init(entry *Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
}


func (f *ModuleFormatter) printColored(b *bytes.Buffer, entry *Entry, keys []string, data Fields, timestampFormat string) {
	var levelColor int
	switch entry.Level {
	case DebugLevel, TraceLevel:
		levelColor = gray
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	levelText := strings.ToUpper(entry.Level.String())
	if !f.DisableLevelTruncation {
		levelText = levelText[0:4]
	}

	// Remove a single newline if it already exists in the message to keep
	// the behavior of logrus text_formatter the same as the stdlib log package
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	caller := ""

	if entry.HasCaller() {
		funcVal := fmt.Sprintf("%s()", entry.Caller.Function)
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)

		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}
		caller = fileVal + " " + funcVal
	}

	if f.DisableTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m%s %-44s ", levelColor, levelText, caller, entry.Message)
	} else if !f.FullTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d]%s %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), caller, entry.Message)
	} else {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s]%s %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), caller, entry.Message)
	}
	fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m", blue, f.module)
	for _, k := range keys {
		v := data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	}
}

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
