package crp

/**
    Author: luzequan
    Created: 2018-08-01 17:30:03
*/
import (
	. "drcs/core/databox"
	"fmt"
	"drcs/core/interaction/request"
	. "drcs/dep/nodelib/crp/edunwang"
	"encoding/json"
	"strings"
)

func init() {
	SUPRESPONSE.Register()
}

var SUPRESPONSE = &DataBox{
	Name:        "sup_response",
	Description: "sup_response",
	RuleTree: &RuleTree{
		Root: supResponseRootFunc,

		Trunk: map[string]*Rule{
			"parseparam": {
				ParseFunc: parseRespParamFunc,
			},
			//"getorderinfo": {
			//	ParseFunc: depAuthFunc,
			//},
			"aesencrypt": {
				ParseFunc: aesEncryptParamFunc,
			},
			"base64encode": {
				ParseFunc: base64EncodeFunc,
			},
			"urlencode": {
				ParseFunc: urlEncodeFunc,
			},
			"execquery": {
				ParseFunc: callResponseFunc,
			},
			"aesdecrypt": {
				ParseFunc: aesDecryptFunc,
			},
			"buildresp": {
				ParseFunc: buildResponseFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func supResponseRootFunc(ctx *Context) {
	fmt.Println("supResponseRootFunc root...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "parseparam",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func parseRespParamFunc(ctx *Context) {
	fmt.Println("parseRespParamFunc rule...")

	reqBody := ctx.GetDataBox().HttpRequestBody

	commonRequestData := &CommonRequestData{}
	err := json.Unmarshal(reqBody, &commonRequestData)
	if err != nil {
		fmt.Println(err.Error())
		errEnd(ctx)
		return
	}
	fmt.Println(commonRequestData)

	ctx.GetDataBox().SetParam("demMemberId", commonRequestData.PubReqInfo.MemId)
	ctx.GetDataBox().SetParam("demMemberId", commonRequestData.PubReqInfo.MemId)
	ctx.GetDataBox().SetParam("demMemberId", commonRequestData.PubReqInfo.MemId)
	ctx.GetDataBox().SetParam("demMemberId", commonRequestData.PubReqInfo.MemId)
	ctx.GetDataBox().SetParam("demMemberId", commonRequestData.PubReqInfo.MemId)

	dataReq := &request.DataRequest{
		Rule:         "aesencrypt",
		Method:       "AESEncrypt",
		TransferType: request.ENCRYPT,
		Reloadable:   true,
		Bobject:      commonRequestData.BusiInfo,
	}

	dataReq.SetParam("memberId", commonRequestData.PubReqInfo.MemId)
	dataReq.SetParam("serialNo", commonRequestData.PubReqInfo.SerialNo)
	dataReq.SetParam("reqSign", commonRequestData.PubReqInfo.ReqSign)
	dataReq.SetParam("appkey", ctx.GetDataBox().Param("appkey"))

	ctx.AddQueue(dataReq)
}

func aesEncryptParamFunc(ctx *Context) {
	fmt.Println("aesEncryptParamFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("aes encrypt failed")
		errEnd(ctx)
		return
	}

	ctx.AddQueue(&request.DataRequest{
		Rule:         "base64encode",
		Method:       "Base64Encode",
		TransferType: request.ENCODE,
		Reloadable:   true,
		Parameters:   ctx.DataResponse.Body,
	})
}

func base64EncodeFunc(ctx *Context) {
	fmt.Println("base64EncodeFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("base encode failed")
		errEnd(ctx)
		return
	}

	dataRequest := &request.DataRequest{
		Rule:         "urlencode",
		Method:       "URLEncode",
		TransferType: request.ENCODE,
		Reloadable:   true,
	}

	dataRequest.SetParam("urlstr", ctx.DataResponse.BodyStr)

	ctx.AddQueue(dataRequest)
}

func urlEncodeFunc(ctx *Context) {
	fmt.Println("urlEncodeFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("url encode failed")
		errEnd(ctx)
		return
	}

	dataRequest := &request.DataRequest{
		Rule:         "execquery",
		Method:       "POST",
		Url:          "http://api.edunwang.com/test/black_check?appid=xxxx&secret_id=xxxx&seq_no=xxx&product_id=xxx&req_data=xxxx",
		TransferType: request.FASTHTTP,
		Reloadable:   true,
	}

	dataRequest.SetParam("appid", ctx.DataResponse.BodyStr)
	dataRequest.SetParam("secret_id", ctx.DataResponse.BodyStr)
	dataRequest.SetParam("seq_no", ctx.DataResponse.BodyStr)
	dataRequest.SetParam("product_id", ctx.DataResponse.BodyStr)
	dataRequest.SetParam("req_data", ctx.DataResponse.BodyStr)

	ctx.AddQueue(dataRequest)

}

func callResponseFunc(ctx *Context) {
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

func aesDecryptFunc(ctx *Context) {
	fmt.Println("aesDecryptFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("exec edunwang query failed")
		errEnd(ctx)
		return
	}

	respData := &RspData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, respData); err != nil {
		fmt.Println("convert respData to struct failed")
		errEnd(ctx)
		return
	}
	pubRespMsg := &PubResProductMsg_0_000_000{}
	pubRespMsg.DetailInfo.Tag = respData.Tag
	pubRespMsg.DetailInfo.EvilScore = respData.EvilScore

	ctx.AddQueue(&request.DataRequest{
		Rule:         "buildresp",
		Method:       "Get",
		TransferType: request.NONETYPE,
		Reloadable:   true,
		Bobject:      pubRespMsg,
	})

}