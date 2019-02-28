# logrus-formatter 
 text formatter for logrus

[Logrus](https://github.com/sirupsen/logrus) 

## Installation
To install formatter, use `go get`:

```sh
$ go get github.com/0x00b/logrus-formatter
```

## Usage
Here is how it should be used:

```go
package main

import (
	"fmt"
	logf "github.com/0x00b/logrus_formatter"
	"github.com/sirupsen/logrus"
	"strings"
)

func formatFuncName(name string) string {
	funcLen := 10
	l := len(name)
	if l > funcLen {
		return name[l-funcLen : l]
	}
	return strings.Repeat(" ", funcLen-l) + name
}
func TestFunctionNameLoooong(log *logrus.Entry) {
	log.Infoln("just test long function name")
}
func main() {
	log := logrus.NewEntry(logrus.New())
	formatter := &logf.TextFormatter{}

	// logf.FieldKeyFile/logf.FieldKeyFunc/logf.FieldKeyLine ,you must call log.Logger.SetReportCaller(true)
	formatter.SetFormat(logf.TagBR, logf.FieldKeyTime, logf.TagBL, logf.FieldKeyLevel, logf.FieldKeyFile, logf.TaGColon, logf.FieldKeyFunc, logf.TaGColon, logf.FieldKeyLine, logf.FieldKeyMsg)
	log.Logger.SetReportCaller(true)
	log.Logger.SetFormatter(formatter)
	log.Data["test"] = "log-formatter"

	log.Printf("format test")
	TestFunctionNameLoooong(log)

	fmt.Println()

	formatter.SetFormat(logf.TagBR, logf.FieldKeyTime, logf.TagBL, logf.FieldKeyLevel, logf.FieldKeyFunc, logf.TaGColon, logf.FieldKeyLine, logf.FieldKeyMsg)
	formatter.FormatFuncName = formatFuncName //自定义函数名的格式。也可以自定义文件名
	formatter.TimestampFormat = "2006-01-02 15:04:05"

	log.Printf("format test")
	TestFunctionNameLoooong(log)
}


will output:
["2019-02-28T14:48:13+08:00"] INFO .ter/example/main.go:      main.main:31   "format test" test:"log-formatter"
["2019-02-28T14:48:13+08:00"] INFO .ter/example/main.go:.ionNameLoooong:19   "just test long function name" test:"log-formatter"

["2019-02-28 14:48:13"] INFO  main.main:40   "format test" test:"log-formatter"
["2019-02-28 14:48:13"] INFO ameLoooong:19   "just test long function name" test:"log-formatter"
```
## API
`logf.TextFormatter` exposes the following fields and methods.

### Fields

* `TimestampFormat string` — timestamp format to use for display when a full timestamp is printed.
* `QuoteEmptyFields bool` — wrap empty fields in quotes if true.
* `FormatFuncName HandlerFormatFile` — custom function name.
* `FormatFileName HandlerFormatFunc` — custom file name.

### Methods

#### `SetFormat(args ...string) (format string)`

Sets an alternative formatting string for output. use the following definition:
* FieldKeyMsg            = logrus.FieldKeyMsg
* FieldKeyLevel          = logrus.FieldKeyLevel
* FieldKeyTime           = logrus.FieldKeyTime
* FieldKeyLogrusError    = logrus.FieldKeyLogrusError
* FieldKeyFunc           = logrus.FieldKeyFunc
* FieldKeyFile           = logrus.FieldKeyFile
* FieldKeyLine           = "line"
* TagBR                  = "["
* TagBL                  = "]"
* TaGColon               = ":"


It's not necessary to call log.Logger.SetReportCaller(true) if you use:
* FieldKeyFunc
* FieldKeyFile
* FieldKeyLine

# License
MIT
