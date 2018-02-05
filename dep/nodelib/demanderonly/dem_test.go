package demanderonly

/**
    Author: luzequan
    Created: 2018-01-02 19:06:23
*/
import (
	"testing"
	"fmt"
	"os"
	"crypto/md5"
	"io"
	"bufio"
)

func TestDem(t *testing.T) {
	//testRule := DEM.RuleTree.Trunk["ruleTest"].ParseFunc
	//fmt.Println("**********", testRule)

	path := "D:/dds_send/JON20171102000000276_ID010201_20171213175701_0001.TARGET"
	targetPath := "D:/dds_send/tmp/JON20171102000000276_ID010201_20171213175701_0001.TARGET"

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
	}

	targetfile, err := os.Open(targetPath)
	defer targetfile.Close()

	io.Copy(targetfile, file)

	buf := make([]byte, 2048)

	md5Str := md5.New()
	for {
		_, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		md5Str.Write(buf)
	}

	md5Strline := fmt.Sprintf("%x", md5Str.Sum(nil))

	fmt.Println(md5Str.Sum(nil))
	fmt.Println(md5Strline)

	file.Seek(0, 0)

	inputReader := bufio.NewReader(file)
	line, _, _:=inputReader.ReadLine()
	fmt.Println("inputReader.ReadLine()", string(line))


	//buf := make([]byte, 1024)
	//
	//seek, err := file.Seek(-1024, os.SEEK_END)
	//
	//_, err = file.Read(buf)
	//if err != nil {
	//	fmt.Println("file read error")
	//}
	//fmt.Println(seek)
	//fmt.Println(string(buf))
	//fmt.Println(md5.Sum(buf))
}
