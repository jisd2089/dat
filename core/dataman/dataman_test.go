package dataman

/**
    Author: luzequan
    Created: 2017-12-28 10:00:16
*/

import (
	"testing"
	"io"
	"fmt"
	"os"
	"strings"
	"bufio"
)

func TestDataMan(t *testing.T) {
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

		fmt.Println(strings.Contains(line, "SOH"))

	}
}
