package logrus_formatter

import (
	"bytes"
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
	TagBL                  = "["
	TagBR                  = "]"
	TaGColon               = ":"
	defaultTimestampFormat = time.RFC3339
)

var (
	defaultFormat      = fmt.Sprintf("[%%%v%%] %%%v%% %%%v%%:%%line%% - %%%v%%", FieldKeyTime, FieldKeyLevel, FieldKeyFunc, FieldKeyMsg)
	defaultFormatArray = []string{TagBL, FieldKeyTime, TagBR, FieldKeyLevel, FieldKeyMsg}

	FunctionNameLength = 25
	FileNameLength     = 20
)

// TextFormatter formats logs into text
type TextFormatter struct {
	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	//LogFormat
	//LogFormat string

	FormatFuncName HandlerFormatFile
	FormatFileName HandlerFormatFunc

	TagSource bool

	hasTime  bool
	hasLevel bool
	hasMsg   bool
	hasFunc  bool
	hasFile  bool
	hasLine  bool

	keyArray []string
}

//HandlerFormatFunc format function name
type HandlerFormatFunc func(funcName string) string

//HandlerFormatFile format file name
type HandlerFormatFile func(fileName string) string

func defaultFormatFunc(funcName string) string {
	length := len(funcName)
	if length > FunctionNameLength {
		return "." + funcName[length-FunctionNameLength+1:length]
	}
	return strings.Repeat(" ", FunctionNameLength-length) + funcName
}
func defaultFormatFile(fileName string) string {
	r := []rune(fileName)
	length := len(r)
	if length > FileNameLength {
		return "." + fileName[length-FileNameLength+1:length]
	}
	return strings.Repeat(" ", FileNameLength-length) + fileName
}

func isTag(s string) bool {
	if s == TagBR || s == TagBL || s == TaGColon {
		return true
	}
	return false
}

func isBR(s string) bool {
	if s == TagBR {
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

func (f *TextFormatter) SetFormat(args ...string) {
	f.keyArray = args
}
func (f *TextFormatter) SetFormatAndTagSource(args ...string) {
	f.keyArray = args
	f.TagSource = true
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer

	if len(f.keyArray) == 0 {
		f.keyArray = defaultFormatArray
	}

	for idx, k := range f.keyArray {
		if isTag(k) {
			buf.WriteString(k)
			if isBR(k) {
				buf.WriteByte(' ')
			}
		} else {
			switch k {
			case FieldKeyMsg:
				buf.WriteString(f.quoteValue(entry.Message))
			case FieldKeyLevel:
				level := Level(entry.Level).String()
				buf.WriteString(level)
			case FieldKeyTime:
				timestampFormat := f.TimestampFormat
				if timestampFormat == "" {
					timestampFormat = defaultTimestampFormat
				}
				buf.WriteString(entry.Time.Format(timestampFormat))
			case FieldKeyFunc:
				if entry.Caller != nil {
					if f.FormatFuncName == nil {
						f.FormatFuncName = defaultFormatFunc
					}
					buf.WriteString(f.FormatFuncName(entry.Caller.Function))
				}
			case FieldKeyFile:
				if entry.HasCaller() && entry.Caller != nil {
					if f.FormatFileName == nil {
						f.FormatFileName = defaultFormatFile
					}
					buf.WriteString(f.FormatFileName(entry.Caller.File))
				}
			case FieldKeyLine:
				if entry.HasCaller() && entry.Caller != nil {
					line := fmt.Sprintf("%-4v", strconv.FormatInt(int64(entry.Caller.Line), 10))
					buf.WriteString(line)
				}
			}
			if idx < len(f.keyArray)-1 && f.keyArray[idx+1] != TagBR && f.keyArray[idx+1] != TaGColon {
				buf.WriteByte(' ')
			}
		}
	}
	length := len(entry.Data)
	if length > 0 {
		var idx int
		buf.WriteString(" (")
		if f.TagSource && entry.HasCaller() && entry.Caller != nil {
			buf.WriteString(fmt.Sprintf("source=%v:%v:%v", entry.Caller.File, entry.Caller.Function, entry.Caller.Line))
			if len(entry.Data) > 0 {
				buf.WriteByte(' ')
			}
		}
		for k, v := range entry.Data {
			if s, ok := v.(string); ok {
				buf.WriteString(fmt.Sprintf("%v=%q", k, s))
			} else {
				buf.WriteString(fmt.Sprintf("%v=%v", k, v))
			}
			if idx < length-1 {
				buf.WriteByte(' ')
			}
			idx++
		}
		buf.WriteByte(')')
	}

	buf.WriteByte('\n')

	return buf.Bytes(), nil
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

// func (f *TextFormatter) SetFormat(args ...string) (format string) {
// 	f.LogFormat = ""
// 	for idx, k := range args {
// 		if isTag(k) {
// 			f.LogFormat += k
// 			if isBR(k) {
// 				f.LogFormat += " "
// 			}
// 		} else {
// 			key := "%" + k + "%"
// 			f.LogFormat += key
// 			if idx < len(args)-1 && args[idx+1] != TagBR && args[idx+1] != TaGColon {
// 				f.LogFormat += " "
// 			}
// 			f.setHasKey(k)
// 		}
// 	}
// 	return f.LogFormat
// }

// func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
// 	output := f.LogFormat
// 	if output == "" {
// 		output = defaultFormat
// 	}

// 	timestampFormat := f.TimestampFormat
// 	if timestampFormat == "" {
// 		timestampFormat = defaultTimestampFormat
// 	}
// 	if f.hasTime {
// 		output = strings.Replace(output, "%"+FieldKeyTime+"%", (entry.Time.Format(timestampFormat)), 1)
// 	}
// 	if f.hasMsg {
// 		output = strings.Replace(output, "%"+FieldKeyMsg+"%", f.quoteValue(entry.Message), 1)
// 	}
// 	if f.hasLevel {
// 		level := Level(entry.Level).String()
// 		output = strings.Replace(output, "%"+FieldKeyLevel+"%", (level), 1)
// 	}
// 	if entry.Caller != nil {
// 		if f.hasFunc {
// 			if f.FormatFuncName == nil {
// 				f.FormatFuncName = defaultFormatFunc
// 			}
// 			output = strings.Replace(output, "%"+FieldKeyFunc+"%", (f.FormatFuncName(entry.Caller.Function)), 1)

// 		}
// 		if f.hasFile {
// 			if f.FormatFileName == nil {
// 				f.FormatFileName = defaultFormatFile
// 			}
// 			output = strings.Replace(output, "%"+FieldKeyFile+"%", (f.FormatFileName(entry.Caller.File)), 1)
// 		}
// 		if f.hasLine {
// 			line := fmt.Sprintf("%-4v", strconv.FormatInt(int64(entry.Caller.Line), 10))
// 			output = strings.Replace(output, "%"+FieldKeyLine+"%", (line), 1)
// 		}
// 	}
// 	for k, v := range entry.Data {
// 		if s, ok := v.(string); ok {
// 			output += fmt.Sprintf(" %v:%q", k, s)
// 		}
// 	}
// 	output += "\n"
// 	return []byte(output), nil
// }

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
