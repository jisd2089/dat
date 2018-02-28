package formatter

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func Benchmark_parse(b *testing.B) {
	fields := make(map[string]interface{}, 6)
	fields["msg"] = "this is message"
	fields["time"] = time.Now()
	fields["level"] = logrus.DebugLevel
	fields["line"] = 10
	fields["file"] = "/home/xiaolie/workspace/go/src/dds/log.go"
	fields["func"] = "dds/log/formatter.Benchmark_parse"

	pattern := "%d{2006-01-02 15:04:05}[%5p] %m #%M.%l (%F)\n"
	// pattern := "%d{2006-01-02 15:04:05}[%5p] %m #%l\n"
	converter := parse(pattern)
	for i := 0; i < b.N; i++ {
		// buffer.Reset()
		buffer := NewRuneBuffer()
		converter.Convert(fields, buffer)
	}
}
