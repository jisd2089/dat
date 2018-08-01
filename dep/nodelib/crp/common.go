package crp

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
	"encoding/json"
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

func buildResponseFunc(ctx *Context) {
	fmt.Println("buildResponseFunc rule...")

	pubRespMsg := ctx.DataResponse.Bobject
	pubResInfo := &PubResInfo{
		ResCode: "",
		ResMsg: "",

	}

	responseInfo := &ResponseInfo{
		PubResInfo: pubResInfo,
		BusiResInfo: pubRespMsg.(map[string]interface{}),
	}

	responseByte, err := json.Marshal(responseInfo)
	if err != nil {
		fmt.Println("parse response info failed")
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().Callback(responseByte)


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

type CommonRequest struct {
	DemMemId string
	SupMemId string

	PreTimeOut      int
	AuthMode        string
	Reqsign         string
	JobId           string
	SerialNo        string
	IdType          string
	DmpSerialNo     string
	BusiSerialNo    string
	BusiReqInfo     map[string]interface{}
	BusiInfoStr     string
	BusiReqInfoHash string
	TimeStamp       string
	TaskIdStr       string
	Timeout         int
	//TaskIdInfo      []xml.TaskIdInfo
	ConnObjStr   string
	SvcStartTime int
	TranHash     string
	SignBlock1   string
	ServiceId    string
	UnitPrice    float64
}

type CommonResponse struct {
	SuccFlag bool
	ResCode  string
	ResMsg   string
	ErrCode  string
	ErrMsg   string
	//FlowStatus   busilog.FlowStatus
	SupMemId     string
	TaskId       string
	BusiSerialNo string
	DmpSerialNo  string
	TimeStamp    string
	PubResInfo   map[string]interface{}
	BusiResInfo  map[string]interface{}
	BusiInfoStr  string
	SvcStartTime int
	SignBlock1   string
	SignMemId2   string
	SignBlock2   string
	Chargflag    bool
}

type CommonRequestData struct {
	PubReqInfo PubReqInfo             `json:"pubReqInfo"`
	BusiInfo   map[string]interface{} `json:"busiInfo"`
}

type PubReqInfo struct {
	MemId     string `json:"memId"`
	SerialNo  string `json:"serialNo"`
	JobId     string `json:"jobId"`
	AuthMode  string `json:"authMode"`
	TimeStamp string `json:"timeStamp"`
	ReqSign   string `json:"reqSign"`
}

type PubAnsInfo struct {
	SerialNo     string `json:"serialNo"`
	BusiSerialNo string `json:"busiSerialNo"`
	ResCode      string `json:"resCode"`
	ResMsg       string `json:"resMsg"`
	TimeStamp    string `json:"timeStamp"`
}

type ResponseInfo struct {
	PubResInfo  *PubResInfo             `json:"PubResInfo"`
	BusiResInfo map[string]interface{} `json:"PubResInfo"`
}

type PubResInfo struct {
	ResCode    string `json:"resCode"`
	ResMsg     string `json:"resMsg"`
	Chargeflag string `json:"chargflag"`
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
