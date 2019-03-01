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
