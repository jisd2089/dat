package fusion

import (
	. "drcs/core/databox"
	"strings"
	logger "drcs/log"
	"drcs/core/interaction/request"
	"time"
	"fmt"
	"os"
	"github.com/valyala/fasthttp"
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
			"uploadsuccess": {
				ParseFunc: uploadSuccessFunc,
			}, "uploaddataSuccessSet": {
				ParseFunc: uploadDataSetSuccessFunc,
			},
			"parseparam": {
				ParseFunc: parseReqParamFunc,
			},
			"supPredictCreditScore": {
				ParseFunc: supPredictCreditScoreFunc,
			},
			"supPredictCreditScoreSuccess": {
				ParseFunc: supPredictCreditScoreSuccessFunc,
			},
			"supPredictCreditScoreCard": {
				ParseFunc: supPredictCreditScoreCardFunc,
			},
			"supPredictCreditScoreCardSuccess": {
				ParseFunc: supPredictCreditScoreCardSuccessFunc,
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
			"predictresponse": {
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

	processType := ctx.GetDataBox().Param("processType")

	fmt.Println(processType)
	switch processType {
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
			Rule:         "supPredictCreditScore",
			Method:       "GET",
			TransferType: request.NONETYPE,
			Reloadable:   true,
			ConnTimeout:  time.Duration(time.Second * 3000),
		})
	case "apiCard":
		ctx.AddChanQueue(&request.DataRequest{
			Rule:         "supPredictCreditScoreCard",
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
		Rule:         "uploaddataSuccessSet",
		Method:       "POSTFILE",
		TransferType: request.FASTHTTP,
		Url:          "http://127.0.0.1:8096/api/drcs/serverAcceptfile",
		Reloadable:   true,
		ConnTimeout:  time.Duration(time.Second * 3000),
		PostData:     ctx.GetDataBox().DataFilePath,
	})

}

func supPredictCreditScoreFunc(ctx *Context) {

	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "supPredictCreditScoreSuccess",
		Method:       "POSTBODY",
		TransferType: request.FASTHTTP,
		Url:          "http://127.0.0.1:8096/api/drcs/serverPredictCreditScore",
		Reloadable:   true,
		HeaderArgs:   header,
		Parameters:   ctx.GetDataBox().HttpRequestBody,
		ConnTimeout:  time.Duration(time.Second * 3000),
	})

}

func supPredictCreditScoreCardFunc(ctx *Context) {

	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "supPredictCreditScoreCardSuccess",
		Method:       "POSTBODY",
		TransferType: request.FASTHTTP,
		Url:          "http://127.0.0.1:8096/api/drcs/serverPredictCreditScoreCard",
		Reloadable:   true,
		HeaderArgs:   header,
		Parameters:   ctx.GetDataBox().HttpRequestBody,
		ConnTimeout:  time.Duration(time.Second * 3000),
	})

}

func supPredictCreditScoreSuccessFunc(ctx *Context) {
	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[supPredictCreditScoreSuccessFunc] resultmsg encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}
	fmt.Println("supPredictCreditScoreSuccessFunc", string(ctx.DataResponse.Body))

	ctx.GetDataBox().BodyChan <- ctx.DataResponse.Body

	procEndFunc(ctx)
}

func supPredictCreditScoreCardSuccessFunc(ctx *Context) {
	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[supPredictCreditScoreCardSuccessFunc] resultmsg encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}
	fmt.Println("supPredictCreditScoreCardSuccessFunc", string(ctx.DataResponse.Body))

	ctx.GetDataBox().BodyChan <- ctx.DataResponse.Body

	procEndFunc(ctx)
}

func uploadDataSetSuccessFunc(ctx *Context) {
	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[uploadDataSetSuccessFunc] resultmsg encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}
	fmt.Println(string(ctx.DataResponse.Body))

	ctx.GetDataBox().BodyChan <- ctx.DataResponse.Body
	fmt.Println(ctx.GetDataBox().DataFilePath)
	os.RemoveAll(ctx.GetDataBox().DataFilePath)

	procEndFunc(ctx)
}

func uploadSuccessFunc(ctx *Context) {

	// 统计数据集条数

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "end",
		Method:       "GET",
		TransferType: request.NONETYPE,
		//Url:          "http://127.0.0.1:8096/api/crp/sup",
		Reloadable: true,
		//ConnTimeout:  time.Duration(time.Second * 3000),
		//PostData: ctx.GetDataBox().DataFilePath,
	})
}

func parseReqParamFunc(ctx *Context) {
	//logger.Info("parseReqParamFunc start ", ctx.GetDataBox().GetId())

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[parseReqParamFunc] ping redis failed: [%s] ", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	//commonRequestData := &CommonRequestData{}
	//err := json.Unmarshal(ctx.GetDataBox().HttpRequestBody, &commonRequestData)
	//if err != nil {
	//	logger.Error("[parseReqParamFunc] unmarshal CommonRequestData err: [%s] ", err.Error())
	//	errEnd(ctx)
	//	return
	//}
	//
	//reqDataJson, err := json.Marshal(commonRequestData.BusiInfo)
	//if err != nil {
	//	logger.Error("[depAuthFunc] marshal request data err: [%s] ", err.Error())
	//	errEnd(ctx)
	//	return
	//}
	//
	//ctx.GetDataBox().HttpRequestBody = reqDataJson

	dataReq := &request.DataRequest{
		Rule:         "applybalance",
		TransferType: request.DEPAUTH,
		Method:       "APPKEY",
		Reloadable:   true,
		//Bobject:      commonRequestData.BusiInfo,
	}

	//dataReq.SetParam("memberId", commonRequestData.PubReqInfo.MemId)
	//dataReq.SetParam("serialNo", commonRequestData.PubReqInfo.SerialNo)
	//dataReq.SetParam("reqSign", commonRequestData.PubReqInfo.ReqSign)
	//dataReq.SetParam("pubkey", ctx.GetDataBox().Param("pubkey"))
	//dataReq.SetParam("jobId", commonRequestData.PubReqInfo.JobId)
	//
	//ctx.GetDataBox().SetParam("demMemberId", commonRequestData.PubReqInfo.MemId)
	//ctx.GetDataBox().SetParam("jobId", commonRequestData.PubReqInfo.JobId)
	//ctx.GetDataBox().SetParam("serialNo", commonRequestData.PubReqInfo.SerialNo)

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
