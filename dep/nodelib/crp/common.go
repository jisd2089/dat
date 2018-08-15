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
	"strings"
	logger "drcs/log"
	"strconv"
	"time"
	"encoding/json"
)

/**
    Author: luzequan
    Created: 2018-05-15 19:30:06
*/

//var json = jsoniter.ConfigCompatibleWithStandardLibrary

func procEndFunc(ctx *Context) {
	//logger.Info("end start")

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
}

func errEnd(ctx *Context) {
	logger.Error(" return for abnormal reason ")

	pubResProductMsg_Error := &common.PubResProductMsg_Error{}
	pubAnsInfo := &common.PubAnsInfo{}
	pubAnsInfo.ResCode = common.CenterCodeReqFailNoCharge
	pubAnsInfo.ResMsg = common.GetCenterCodeText(common.CenterCodeReqFailNoCharge)
	pubAnsInfo.SerialNo = ctx.GetDataBox().Param("serialNo")
	pubAnsInfo.BusiSerialNo = ctx.GetDataBox().Param("busiSerialNo")
	pubAnsInfo.TimeStamp = strconv.Itoa(int(time.Now().UnixNano() / 1e6))

	pubResProductMsg_Error.PubAnsInfo = pubAnsInfo

	responseByte, err := json.Marshal(pubResProductMsg_Error)
	if err != nil {
		responseByte = []byte("response error")
	}

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

	ctx.AddChanQueue(&request.DataRequest{
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

type SeqUtil struct {}

// seq no 长度15位
func (s SeqUtil) GenSeqNo() string {
	now := time.Now()
	serialNo := fmt.Sprintf("%02d%02d%02d%09d", now.Hour(), now.Minute(), now.Second(), now.Nanosecond())

	return serialNo
}

