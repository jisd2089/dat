package formatter

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestParse(t *testing.T) {
	fields := make(map[string]interface{}, 6)
	fields["msg"] = "this is message"
	fields["time"] = time.Now()
	fields["level"] = logrus.DebugLevel
	fields["line"] = 10
	fields["file"] = "/home/xiaolie/workspace/go/src/dds/log.go"
	fields["func"] = "dds/log/formatter.Benchmark_parse"
	fields["error"] = errors.New("valueerror")
	// pattern := "%msg %m %d [%date{20060102 150405}] %p %level %l %line %F %file\n"
	pattern := "%d{2006-01-02 15:04:05}[%!.4p] %m #%M.%l (%F) error=%e\n"
	converter := parse(pattern)
	buffer := NewRuneBuffer()
	converter.Convert(fields, buffer)
	fmt.Println("=====")
	fmt.Print(buffer.String())
}
