package fusion

import (
	. "drcs/core/databox"
	"strings"
	logger "drcs/log"
	"drcs/core/interaction/request"
	"time"
)

/**
    Author: luzequan
    Created: 2018-09-03 11:55:06
*/
func init() {
	CUSTOMER.Register()
}

var CUSTOMER = &DataBox{
	Name:        "customer_request",
	Description: "customer_request",
	RuleTree: &RuleTree{
		Root: customerRootFunc,

		Trunk: map[string]*Rule{
			"uploaddataset": {
				ParseFunc: uploadDataSetFunc,
			},
			"parseparam": {
				ParseFunc: parseReqParamFunc,
			},
			"depauth": {
				ParseFunc: depAuthFunc,
			},
			"predictcreditscore": {
				ParseFunc: predictCreditScoreFunc,
			},
			"predictcreditscorecard": {
				ParseFunc: predictCreditScoreCardFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func customerRootFunc(ctx *Context) {
	//logger.Info("customerRootFunc start ", ctx.GetDataBox().GetId())

	switch typ {
	case "upload":
		ctx.AddChanQueue(&request.DataRequest{
			Rule:         "uploaddataset",
			Method:       "GET",
			TransferType: request.NONETYPE,
			Reloadable:   true,
			ConnTimeout:  time.Duration(time.Second * 3000),
		})
	case "api":
		ctx.AddChanQueue(&request.DataRequest{
			Rule:         "parseparam",
			Method:       "GET",
			TransferType: request.NONETYPE,
			Reloadable:   true,
			ConnTimeout:  time.Duration(time.Second * 3000),
		})
	}
}

//
func uploadDataSetFunc(ctx *Context) {
	//logger.Info("uploadDataSetFunc start ", ctx.GetDataBox().GetId())

	// 业务逻辑：例如：txt convert to csv

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "end",
		Method:       "POSTFILE",
		TransferType: request.FASTHTTP,
		Url:          "http://127.0.0.1:8096/api/crp/sup",
		Reloadable:   true,
		ConnTimeout:  time.Duration(time.Second * 3000),
	})

}

func parseReqParamFunc(ctx *Context) {
	//logger.Info("parseReqParamFunc start ", ctx.GetDataBox().GetId())

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[parseReqParamFunc] ping redis failed: [%s] ", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	commonRequestData := &CommonRequestData{}
	err := json.Unmarshal(ctx.GetDataBox().HttpRequestBody, &commonRequestData)
	if err != nil {
		logger.Error("[parseReqParamFunc] unmarshal CommonRequestData err: [%s] ", err.Error())
		errEnd(ctx)
		return
	}

	reqDataJson, err := json.Marshal(commonRequestData.BusiInfo)
	if err != nil {
		logger.Error("[depAuthFunc] marshal request data err: [%s] ", err.Error())
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().HttpRequestBody = reqDataJson

	dataReq := &request.DataRequest{
		Rule:         "applybalance",
		TransferType: request.DEPAUTH,
		Method:       "APPKEY",
		Reloadable:   true,
		//Bobject:      commonRequestData.BusiInfo,
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

func depAuthFunc(ctx *Context) {
	//logger.Info("depAuthFunc start ", ctx.GetDataBox().GetId())

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[depAuthFunc] dep authentication failed: [%s] ", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	ctx.AddChanQueue(&request.DataRequest{
		Rule: "applybalance",
		//Rule:         "reduceredisquato",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
		//Parameters:   reqDataJson,
		ConnTimeout: time.Duration(time.Second * 3000),
	})
}

func predictCreditScoreFunc(ctx *Context) {

}

func predictCreditScoreCardFunc(ctx *Context) {

}
