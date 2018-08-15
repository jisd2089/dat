package crp

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	. "drcs/dep/nodelib/crp/common"
	. "drcs/dep/nodelib/crp/smartsail"
	"strings"
	"sync"
	"time"
	"strconv"
	logger "drcs/log"
	"crypto/md5"
	"fmt"
	"encoding/json"
)

func init() {
	SMARTREQUEST.Register()
}

var slock sync.Mutex

var SMARTREQUEST = &DataBox{
	Name:        "smart_request",
	Description: "smart_request",
	RuleTree: &RuleTree{
		Root: smartRequestRootFunc,

		Trunk: map[string]*Rule{
			"parseparam": {
				ParseFunc: parseRequestParamFunc,
			},
			"depauth": {
				ParseFunc: depAuthFunc,
			},
			"getorderinfo": {
				ParseFunc: depAuthFunc,
			},
			"applybalance": {
				ParseFunc: applyBalanceFunc,
			},
			"updateredisquato": {
				ParseFunc: updateRedisQuatoFunc,
			},
			"reduceredisquato": {
				ParseFunc: reduceRedisQuatoFunc,
			},
			"singlequery": {
				ParseFunc: singleQueryFunc,
			},
			"staticquery": {
				ParseFunc: staticQueryFunc,
			},
			"queryresponse": {
				ParseFunc: callSmartResponseFunc,
			},
			"aesdecrypt": {
				ParseFunc: aesDecryptFunc,
			},
			"buildresp": {
				ParseFunc: callResponseFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func smartRequestRootFunc(ctx *Context) {
	logger.Info("smartRequestRootFunc start")

	start := int(time.Now().UnixNano() / 1e6)

	ctx.GetDataBox().SetParam("startTime", strconv.Itoa(start))

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "parseparam",
		Method:       "PING",
		TransferType: request.REDIS,
		Reloadable:   true,
		CommandParams: ctx.GetDataBox().Params,
	})

	// Test
	//pubRespMsg := &PubResProductMsg{}
	//
	//pubAnsInfo := &PubAnsInfo{}
	//pubAnsInfo.ResCode = "000000"
	//pubAnsInfo.ResMsg = "成功"
	//pubRespMsg.PubAnsInfo = pubAnsInfo
	//
	//bodyByte, err := json.Marshal(pubRespMsg)
	//
	//if err != nil {
	//	logger.Error("[callSmartResponseFunc] unmarshal response body to PubResProductMsg_0_000_000 err: [%s] ", err.Error())
	//	return
	//}
	//
	//ctx.GetDataBox().BodyChan <- bodyByte
	//
	//procEndFunc(ctx)
}

func parseRequestParamFunc(ctx *Context) {
	logger.Info("parseRequestParamFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[parseRequestParamFunc] ping redis failed: [%s] ", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	reqBody := ctx.GetDataBox().HttpRequestBody

	commonRequestData := &CommonRequestData{}
	err := json.Unmarshal(reqBody, &commonRequestData)
	if err != nil {
		logger.Error("[parseRequestParamFunc] unmarshal CommonRequestData err: [%s] ", err.Error())
		errEnd(ctx)
		return
	}

	dataReq := &request.DataRequest{
		Rule:         "depauth",
		TransferType: request.DEPAUTH,
		Method:       "APPKEY",
		Reloadable:   true,
		Bobject:      commonRequestData.BusiInfo,
	}

	dataReq.SetParam("memberId", commonRequestData.PubReqInfo.MemId)
	dataReq.SetParam("serialNo", commonRequestData.PubReqInfo.SerialNo)
	dataReq.SetParam("reqSign", commonRequestData.PubReqInfo.ReqSign)
	dataReq.SetParam("pubkey", ctx.GetDataBox().Param("pubkey"))
	dataReq.SetParam("jobId", commonRequestData.PubReqInfo.JobId)

	ctx.GetDataBox().SetParam("demMemberId", commonRequestData.PubReqInfo.MemId)
	ctx.GetDataBox().SetParam("jobId", commonRequestData.PubReqInfo.JobId)
	ctx.GetDataBox().SetParam("serialNo", commonRequestData.PubReqInfo.SerialNo)

	ctx.AddChanQueue(dataReq)
}

func callSmartResponseFunc(ctx *Context) {
	logger.Info("callSmartResponseFunc start")

	pubRespMsg := &PubResProductMsg{}
	// TODO mock
	//pubAnsInfo := &PubAnsInfo{}
	//pubAnsInfo.ResCode = "000000"
	//pubAnsInfo.ResMsg = "成功"
	//pubRespMsg.PubAnsInfo = pubAnsInfo
	//ctx.DataResponse.Body, _ = json.Marshal(pubRespMsg)

	if err := json.Unmarshal(ctx.DataResponse.Body, pubRespMsg); err != nil {
		logger.Error("[callSmartResponseFunc] unmarshal response body to PubResProductMsg_0_000_000 err: [%s] ", err.Error())
		returnBalanceFunc(ctx)
		return
	}


	ctx.GetDataBox().BodyChan <- ctx.DataResponse.Body

	// 不收费处理逻辑
	if strings.EqualFold(pubRespMsg.PubAnsInfo.ResCode, CenterCodeReqFailNoCharge) {
		ctx.AddChanQueue(&request.DataRequest{
			Rule:         "returnbalance",
			Method:       "GET",
			TransferType: request.NONETYPE,
			Reloadable:   true,
		})
		return
	}

	demMemberId := ctx.GetDataBox().Param("demMemberId")
	busiSerialNo := ctx.GetDataBox().Param("busiSerialNo")
	start := ctx.GetDataBox().Param("startTime")
	startTime, err := strconv.Atoi(start)
	if err != nil {
		logger.Error("[callResponseFunc] convert startTime string to int err: [%s] ", err.Error())
		returnBalanceFunc(ctx)
		return
	}
	endTime := int(time.Now().UnixNano() / 1e6)

	h := md5.New()
	h.Write([]byte(demMemberId))
	busiInfoStr := fmt.Sprintf("%x", h.Sum(nil))

	//msg := "3" + demMemberId + "1"
	//priKey, _ := security.GetPrivateKey()
	//signInfo := cncrypt.Sign(priKey, []byte(msg))
	//stepInfoM := []map[string]interface{}{}
	//stepInfo1 := map[string]interface{}{"no": 1, "memID": demMemberId, "stepStatus": security.StepStatusSucc, "signature": ""}
	//stepInfo2 := map[string]interface{}{"no": 2, "memID": "", "stepStatus": security.StepStatusSucc, "signature": ""}
	//stepInfo3 := map[string]interface{}{"no": 3, "memID": demMemberId, "stepStatus": security.StepStatusSucc, "signature": signInfo}
	//stepInfoM = append(stepInfoM, stepInfo1)
	//stepInfoM = append(stepInfoM, stepInfo2)
	//stepInfoM = append(stepInfoM, stepInfo3)

	ctx.Output(map[string]interface{}{
		"exID":       busiInfoStr,
		"demMemID":   demMemberId,
		"supMemID":   ctx.GetDataBox().Param("supMemId"),
		"taskID":     strings.Replace(ctx.GetDataBox().Param("taskIdStr"), "|@|", ".", -1),
		"seqNo":      busiSerialNo,
		"recordType": RecordTypeSingle,
		"succCount":  "1",
		"flowStatus": FlowStatusDemSucc,
		"usedTime":   endTime - startTime,
		"errCode":    ErrCodeSucc,
		//"stepInfoM":  stepInfoM,
		//"dmpSeqNo":   "",
	})

	procEndFunc(ctx)
}
