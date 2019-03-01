package logrus_formatter

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"testing"
)

type LogOutput struct {
	buffer string
}

func (o *LogOutput) Write(p []byte) (int, error) {
	o.buffer += string(p[:])
	return len(p), nil
}

func (o *LogOutput) GetValue() string {
	return o.buffer
}

func TestLogrusFormatter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogrusFormatter Suite")
}

var _ = Describe("Formatter", func() {
	var formatter *TextFormatter
	var log *logrus.Logger
	var output *LogOutput

	BeforeEach(func() {
		output = new(LogOutput)
		formatter = new(TextFormatter)
		log = logrus.New()
		log.Out = output
		log.Formatter = formatter
		log.Level = logrus.DebugLevel
	})

	Describe("logfmt output", func() {
		It("should output simple message", func() {
			formatter.SetFormat(TagBL, FieldKeyLevel, TagBR, FieldKeyMsg)
			log.Debug("test")
			Ω(output.GetValue()).Should(Equal("[DEBG] test\n"))
		})

		It("should output message with additional field", func() {
			formatter.SetFormat(FieldKeyLevel, FieldKeyMsg)
			log.WithFields(logrus.Fields{"animal": "walrus"}).Debug("test")
			Ω(output.GetValue()).Should(Equal("DEBG test (animal:\"walrus\")\n"))
		})
	})

	Describe("Formatted output", func() {
		It("should output formatted message", func() {
			formatter.SetFormat(TagBL, FieldKeyTime, TagBR, FieldKeyLevel, FieldKeyFile, TaGColon, FieldKeyFunc, TaGColon, FieldKeyLine, FieldKeyMsg)
			log.Warnln("warnning test ")
			fmt.Println(output.GetValue())
		})
	})

	Describe("Theming support", func() {

	})
})
