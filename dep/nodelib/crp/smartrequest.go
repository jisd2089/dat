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
	"fmt"
	"strings"
	"encoding/json"
	"sync"
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
			//"getpolicy": {
			//	ParseFunc: getOrderRoutePolicyFunc,
			//},
			//"aesencrypt": {
			//	ParseFunc: aesEncryptParamFunc,
			//},
			//"base64encode": {
			//	ParseFunc: base64EncodeFunc,
			//},
			//"urlencode": {
			//	ParseFunc: urlEncodeFunc,
			//},
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
	fmt.Println("smartRequest Root ...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "parseparam",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func parseRequestParamFunc(ctx *Context) {
	fmt.Println("parseRequestParamFunc rule...")

	reqBody := ctx.GetDataBox().HttpRequestBody

	commonRequestData := &CommonRequestData{}
	err := json.Unmarshal(reqBody, &commonRequestData)
	if err != nil {
		fmt.Println(err.Error())
		errEnd(ctx)
		return
	}
	fmt.Println(commonRequestData)

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

	ctx.GetDataBox().SetParam("jobId", commonRequestData.PubReqInfo.JobId)

	ctx.AddQueue(dataReq)
}

func callSmartResponseFunc(ctx *Context) {
	fmt.Println("callSmartResponseFunc rule...")

	respData := &RespDetail{}
	if err := json.Unmarshal(ctx.DataResponse.Body, respData); err != nil {
		errEnd(ctx)
		return
	}

	pubRespMsg := &PubResProductMsg_0_000_000{}
	//pubRespMsg.DetailInfo.Tag = respData.Tag
	//pubRespMsg.DetailInfo.EvilScore = respData.EvilScore

	responseByte, err := json.Marshal(pubRespMsg)
	if err != nil {
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

	errEnd(ctx)
}
