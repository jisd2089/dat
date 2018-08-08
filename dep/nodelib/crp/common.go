package crp

import (
	. "drcs/core/databox"
	"fmt"
	"drcs/runtime/status"
	"os"
	"drcs/core/interaction/request"
	"drcs/dep/nodelib/crp/common"
	"drcs/dep/member"
	"bufio"
	"crypto/md5"
	"bytes"
	"io"
	"encoding/json"
	"strings"
	logger "drcs/log"
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
	logger.Error(" return for abnormal reason ")

	responseByte := []byte("response error")
	ctx.GetDataBox().BodyChan <- responseByte
	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "end",
		TransferType: request.NONETYPE,
		Priority:     1,
		Reloadable:   true,
	})
}

func buildResponseFunc(ctx *Context) {
	fmt.Println("buildResponseFunc rule...")

	pubRespMsg := ctx.DataResponse.Bobject
	pubResInfo := &common.PubResInfo{
		ResCode: "",
		ResMsg:  "",
	}

	responseInfo := &common.ResponseInfo{
		PubResInfo:  pubResInfo,
		BusiResInfo: pubRespMsg.(map[string]interface{}),
	}

	responseByte, err := json.Marshal(responseInfo)
	if err != nil {
		fmt.Println("parse response info failed")
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().BodyChan <- responseByte

	ctx.Output(map[string]interface{}{
		//"exID":       string(line),
		"demMemID":   ctx.GetDataBox().Param("UserId"),
		"supMemID":   ctx.GetDataBox().Param("NodeMemberId"),
		"taskID":     strings.Replace(ctx.GetDataBox().Param("TaskId"), "|@|", ".", -1),
		"seqNo":      ctx.GetDataBox().Param("seqNo"),
		"dmpSeqNo":   ctx.GetDataBox().Param("fileNo"),
		"recordType": "2",
		"succCount":  "1",
		"flowStatus": "11",
		"usedTime":   11,
		"errCode":    "031014",
		//"stepInfoM":  stepInfoM,
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

func getPartnerUrl() (string, string, error) {
	//p, err := member.GetPartnerInfoById(memberId)
	//if err != nil {
	//	return "", "", err
	//}
	svrUrl := member.GetPartnersInfo().PartnerDetailList.PartnerDetailInfo[0].SvrURL
	memberId := member.GetPartnersInfo().PartnerDetailList.PartnerDetailInfo[0].MemberId
	return svrUrl, memberId, nil
}

