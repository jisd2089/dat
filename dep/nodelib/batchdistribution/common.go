package batchdistribution

import (
	. "drcs/core/databox"

	"fmt"
	"drcs/runtime/status"
	"os"
	"drcs/core/interaction/request"
	"bufio"
	"crypto/md5"
	"bytes"
	"io"
)

/**
    Author: luzequan
    Created: 2018-05-15 19:30:06
*/
func procEndFunc(ctx *Context) {
	fmt.Println("end start ...")

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
}

func errEnd(ctx *Context) {
	ctx.AddQueue(&request.DataRequest{
		Rule:         "end",
		TransferType: request.NONETYPE,
		Priority:     1,
		Reloadable:   true,
	})
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

type BatchRequest struct {
	SeqNo     string
	TaskId    string
	UserId    string
	JobId     string
	IdType    string
	DataRange string
	MaxDelay  int
	DmpSeqNo  string
	MD5       string
}

func getMD5(dataFilePath string) (string, error) {

	dataFile, err := os.Open(dataFilePath)
	defer dataFile.Close()
	if err != nil {
		return "", err
	}

	buf := bufio.NewReader(dataFile)

	md5Hash := md5.New()
	lineCnt := 300
	cntBuf := &bytes.Buffer{}
	c := lineCnt
	for {
		c--
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF && cntBuf.Len() == 0 {
				break
			} else if err == io.EOF && cntBuf.Len() > 0 {
				md5Hash.Write(cntBuf.Bytes())
				break
			} else {
				return "", fmt.Errorf("read data file error")
			}
		}

		cntBuf.Write(line)
		cntBuf.WriteByte('\n')

		if c == 0 {

			md5Hash.Write(cntBuf.Bytes())

			c = lineCnt
			cntBuf.Reset()
		}
	}

	md5Str := fmt.Sprintf("%x", md5Hash.Sum(nil))
	return md5Str, nil
}
