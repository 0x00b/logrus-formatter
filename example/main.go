package main

import (
	"fmt"
	logf "github.com/0x00b/logrus-formatter"
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

	// logf.FieldKeyFile/logf.FieldKeyFunc/logf.FieldKeyLine , must call log.Logger.SetReportCaller(true)
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
