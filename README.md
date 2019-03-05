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
	logf "github.com/0x00b/logrus-formatter"
	"github.com/sirupsen/logrus"
	"strings"
)

func TestFunctionNameLoooong(log *logrus.Entry) {
	log.Infoln("just test long function name")
}
func main() {
	log := logrus.NewEntry(logrus.New())

	formatter := &logf.TextFormatter{}
	log.Logger.SetFormatter(formatter)

	log.Printf("format test")    //[2019-03-01T19:39:22+08:00] INFO "format test"
	TestFunctionNameLoooong(log) //[2019-03-01T19:39:22+08:00] INFO "just test long function name"

	//set logrus Data
	log.Data["name"] = "ice"
	log.Data["age"] = 18

	//set format as : [time] level file:func:line msg
	formatter.SetFormat(logf.TagBL, logf.FieldKeyTime, logf.TagBR, logf.FieldKeyLevel, logf.FieldKeyFile, logf.TaGColon, logf.FieldKeyFunc, logf.TaGColon, logf.FieldKeyLine, logf.FieldKeyMsg)
	// logf.FieldKeyFile/logf.FieldKeyFunc/logf.FieldKeyLine , must call log.Logger.SetReportCaller(true)
	log.Logger.SetReportCaller(true)
	//set timestamp format
	formatter.TimestampFormat = "2006-01-02 15:04:05"

	log.Printf("format test")    //[2019-03-01 19:40:48] INFO .ter/example/main.go:                main.main:37   "format test" (name:"ice" age:18)
	TestFunctionNameLoooong(log) //[2019-03-01 19:40:48] INFO .ter/example/main.go:..TestFunctionNameLoooong:18   "just test long function name" (name:"ice" age:18)

	logf.FileNameLength = 10     //
	logf.FunctionNameLength = 10 //
	log.Printf("format test")    //[2019-03-01 19:51:10] INFO .e/main.go: main.main:37   "format test" (name:"ice" age:18)

	//自定义函数名的格式。也可以自定义文件名
	formatter.FormatFuncName = func(name string) string {
		funcLen := 5
		l := len(name)
		if l > funcLen {
			return name[l-funcLen : l]
		}
		return strings.Repeat(" ", funcLen-l) + name
	}

	log.Printf("format test") //[2019-03-01 19:51:10] INFO .e/main.go:.main:49   "format test" (name:"ice" age:18)
}

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
```go
FieldKeyMsg            = logrus.FieldKeyMsg		//"msg"
FieldKeyLevel          = logrus.FieldKeyLevel		//"level"
FieldKeyTime           = logrus.FieldKeyTime		//"time"
FieldKeyFunc           = logrus.FieldKeyFunc		//"func"
FieldKeyFile           = logrus.FieldKeyFile		//"file"
FieldKeyLine           = "line"	
TagBL                  = "["
TagBR                  = "]"
TaGColon               = ":"

You can use the above fields to combine the format you want.
eg: 
[2019-03-01 19:51:10] INFO "format test"
SetFormat(TagBL, FieldKeyTime, TagBR, FieldKeyLevel, FieldKeyMsg)

2019-03-05 11:48:08 INFO .ter/example/main.go  main.main 48   "format test"
SetFormat(FieldKeyTime, FieldKeyLevel, FieldKeyFile, FieldKeyFunc, FieldKeyLine, FieldKeyMsg)

2019-03-05 11:49:47 INFO .ter/example/main.go:48    main.main "format test"
SetFormat(FieldKeyTime, FieldKeyLevel, FieldKeyFile, TaGColon, FieldKeyLine, FieldKeyFunc, FieldKeyMsg)

It's necessary to call log.Logger.SetReportCaller(true) if you use:
* FieldKeyFunc
* FieldKeyFile
* FieldKeyLine
```
# License
MIT
