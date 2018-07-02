package batchdistribution

import (
	. "drcs/core/databox"

	"fmt"
	"drcs/runtime/status"
	"os"
	"drcs/core/interaction/request"
	"drcs/dep/member"
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
	SeqNo      string
	TaskIdStr  string
	TaskIdList []string
	UserId     string
	JobId      string
	IdType     string
	DataRange  string
	MaxDelay   int
	DmpSeqNo   string
	MD5        string
	LineCount  int
}

func genMD5beforeSend(ctx *Context, nextRule string, batchRequest *BatchRequest) {

	md5Str, total, err := getMD5(ctx.GetDataBox().DataFilePath)
	if err != nil {
		errEnd(ctx)
		return
	}

	fmt.Println("rcv md5: ", md5Str)

	batchRequest.MD5 = md5Str
	batchRequest.LineCount = total

	ctx.AddQueue(&request.DataRequest{
		Rule:         nextRule,
		TransferType: request.NONETYPE,
		Priority:     1,
		Reloadable:   true,
	})
}

func getMD5(dataFilePath string) (string, int, error) {

	dataFile, err := os.Open(dataFilePath)
	defer dataFile.Close()
	if err != nil {
		return "", 0, err
	}

	buf := bufio.NewReader(dataFile)

	md5Hash := md5.New()
	lineCnt := 300
	cntBuf := &bytes.Buffer{}
	c := lineCnt
	t := 0
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
				return "", 0, fmt.Errorf("read data file error")
			}
		}

		t ++

		cntBuf.Write(line)
		cntBuf.WriteByte('\n')

		if c == 0 {

			md5Hash.Write(cntBuf.Bytes())

			c = lineCnt
			cntBuf.Reset()
		}
	}

	md5Str := fmt.Sprintf("%x", md5Hash.Sum(nil))
	return md5Str, t, nil
}

func getMemberUrls(taskInfoMap map[string]string) ([]string, []string) {
	var svcUrls []string
	var supMemId []string

	for k, _ := range taskInfoMap {
		if p, err := member.GetPartnerInfoById(k); err == nil {
			url := p.SvrURL
			if len(url) > 0 {
				svcUrls = append(svcUrls, url)
				supMemId = append(supMemId, k)
			}
		}
	}
	return svcUrls, supMemId
}

func getPartnerUrl(memberId string) (string, string, error) {
	p, err := member.GetPartnerInfoById(memberId)
	if err != nil {
		return "", "", err
	}
	return p.SvrURL, p.MemberId, nil
}
