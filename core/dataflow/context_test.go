package dataflow

/**
    Author: luzequan
    Created: 2017-12-28 09:50:30
*/
import (
	"testing"
	"io"
	"fmt"
	"os"
	"strings"
	"bufio"
)

func TestContext(t *testing.T) {
	inputFile := "D:\\document\\接力配送\\after_140.data"
	f, err := os.Open(inputFile)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	buf := bufio.NewReader(f)

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		fmt.Println(line)
	}
}