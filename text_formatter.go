package logrus_formatter

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	FieldKeyMsg            = logrus.FieldKeyMsg
	FieldKeyLevel          = logrus.FieldKeyLevel
	FieldKeyTime           = logrus.FieldKeyTime
	FieldKeyLogrusError    = logrus.FieldKeyLogrusError
	FieldKeyFunc           = logrus.FieldKeyFunc
	FieldKeyFile           = logrus.FieldKeyFile
	FieldKeyLine           = "line"
	TagBR                  = "["
	TagBL                  = "]"
	TaGColon               = ":"
	defaultTimestampFormat = time.RFC3339
)

var (
	defaultFormat = fmt.Sprintf("[%%%v%%] %%%v%% %%%v%%:%%line%% - %%%v%%", FieldKeyTime, FieldKeyLevel, FieldKeyFunc, FieldKeyMsg)
)

// TextFormatter formats logs into text
type TextFormatter struct {
	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	//LogFormat
	LogFormat string

	FormatFuncName HandlerFormatFile
	FormatFileName HandlerFormatFunc

	hasTime  bool
	hasLevel bool
	hasMsg   bool
	hasFunc  bool
	hasFile  bool
	hasLine  bool
}

//HandlerFormatFunc format function name
type HandlerFormatFunc func(funcName string) string

//HandlerFormatFile format file name
type HandlerFormatFile func(fileName string) string

func defaultFormatFunc(funcName string) string {
	funcLen := 15
	l := len(funcName)
	if l > funcLen {
		return "." + funcName[l-funcLen+1:l]
	}
	return strings.Repeat(" ", funcLen-l) + funcName
}
func defaultFormatFile(fileName string) string {
	fileLen := 20
	r := []rune(fileName)
	l := len(r)
	if l > fileLen {
		return "." + fileName[l-fileLen+1:l]
	}
	return strings.Repeat(" ", fileLen-l) + fileName
}

func isTag(s string) bool {
	if s == TagBR || s == TagBL || s == TaGColon {
		return true
	}
	return false
}

func isBL(s string) bool {
	if s == TagBL {
		return true
	}
	return false
}

func (f *TextFormatter) setHasKey(k string) {
	switch true {
	case k == FieldKeyTime:
		f.hasTime = true
	case k == FieldKeyLevel:
		f.hasLevel = true
	case k == FieldKeyMsg:
		f.hasMsg = true
	case k == FieldKeyFunc:
		f.hasFunc = true
	case k == FieldKeyFile:
		f.hasFile = true
	case k == FieldKeyLine:
		f.hasLine = true
	}
}

func (f *TextFormatter) SetFormat(args ...string) (format string) {
	f.LogFormat = ""
	for idx, k := range args {
		if isTag(k) {
			f.LogFormat += k
			if isBL(k) {
				f.LogFormat += " "
			}
		} else {
			key := "%" + k + "%"
			f.LogFormat += key
			if idx < len(args)-1 && args[idx+1] != TagBL && args[idx+1] != TaGColon {
				f.LogFormat += " "
			}
			f.setHasKey(k)
		}
	}
	return f.LogFormat
}

type Level logrus.Level

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	} else {
		return "unknown"
	}
}

func (level Level) MarshalText() ([]byte, error) {
	switch logrus.Level(level) {
	case logrus.TraceLevel:
		return []byte("TRAC"), nil
	case logrus.DebugLevel:
		return []byte("DEBG"), nil
	case logrus.InfoLevel:
		return []byte("INFO"), nil
	case logrus.WarnLevel:
		return []byte("WARN"), nil
	case logrus.ErrorLevel:
		return []byte("ERRO"), nil
	case logrus.FatalLevel:
		return []byte("FATA"), nil
	case logrus.PanicLevel:
		return []byte("PANC"), nil
	}

	return nil, fmt.Errorf("not a valid lorus level %q", level)
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = defaultFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if f.hasTime {
		output = strings.Replace(output, "%"+FieldKeyTime+"%", f.quoteValue(entry.Time.Format(timestampFormat)), 1)
	}
	if f.hasMsg {
		output = strings.Replace(output, "%"+FieldKeyMsg+"%", f.quoteValue(entry.Message), 1)
	}
	if f.hasLevel {
		level := Level(entry.Level).String()
		output = strings.Replace(output, "%"+FieldKeyLevel+"%", (level), 1)
	}
	if entry.Caller != nil {
		if f.hasFunc {
			if f.FormatFuncName == nil {
				f.FormatFuncName = defaultFormatFunc
			}
			output = strings.Replace(output, "%"+FieldKeyFunc+"%", (f.FormatFuncName(entry.Caller.Function)), 1)

		}
		if f.hasFile {
			if f.FormatFileName == nil {
				f.FormatFileName = defaultFormatFile
			}
			output = strings.Replace(output, "%"+FieldKeyFile+"%", (f.FormatFileName(entry.Caller.File)), 1)
		}
		if f.hasLine {
			line := fmt.Sprintf("%-4v", strconv.FormatInt(int64(entry.Caller.Line), 10))
			output = strings.Replace(output, "%"+FieldKeyLine+"%", (line), 1)
		}
	}
	for k, v := range entry.Data {
		if s, ok := v.(string); ok {
			output += fmt.Sprintf(" %v:%q", k, s)
		}
	}
	output += "\n"
	return []byte(output), nil
}

func (f *TextFormatter) needsQuoting(text string) bool {
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

func (f *TextFormatter) quoteValue(value interface{}) string {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		return (stringVal)
	} else {
		return (fmt.Sprintf("%q", stringVal))
	}
}
