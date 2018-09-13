package fusion

import (
	. "drcs/core/databox"
	"strings"
	logger "drcs/log"
	"drcs/core/interaction/request"
	. "drcs/dep/nodelib/fusion/common"
	"time"
	"fmt"
	"github.com/valyala/fasthttp"
	"drcs/dep/or"
	"drcs/dep/member"
	"drcs/dep/order"
	"crypto/md5"
	"drcs/dep/security"
	"drcs/common/cncrypt"
	"strconv"
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
			"uploaddataSuccessSet": {
				ParseFunc: uploadDataSetSuccessFunc,
			},
			"parseparam": {
				ParseFunc: parseReqParamFunc,
			},
			"suppredict": {
				ParseFunc: supPredictFunc,
			},
			"suppredictresponse": {
				ParseFunc: supPredictResponseFunc,
			},
			"depauth": {
				ParseFunc: depAuthFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func customerRootFunc(ctx *Context) {
	//logger.Info("customerRootFunc start ", ctx.GetDataBox().GetId())

	switch ctx.GetDataBox().Param("processType") {
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

func uploadDataSetFunc(ctx *Context) {
	//logger.Info("uploadDataSetFunc start ", ctx.GetDataBox().GetId())

	// 业务逻辑
	jobId := ctx.GetDataBox().Param("jobId")

	// 根据jobid获取order map
	orderInfoMap, ok := order.GetOrderInfoMap()[jobId]
	if !ok {
		logger.Error("[uploadDataSetFunc] get order list by jobId [%s] failed", jobId)
		return
	}

	odl := orderInfoMap.OrderDetailList[0]
	if odl == nil {
		logger.Error("[uploadDataSetFunc] get order detail info failed")
		return
	}

	memberDetailInfo, err := member.GetPartnerInfoById(odl.DemMemId)
	if err != nil {
		logger.Error("[uploadDataSetFunc] get partner info by memberid [%s] error: [%s]", odl.DemMemId, err.Error())
		return
	}

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "uploaddataSuccessSet",
		Method:       "POSTFILE",
		TransferType: request.FASTHTTP,
		Url:          memberDetailInfo.SvrURL,
		//Url:          "http://127.0.0.1:8096/api/drcs/serverAcceptfile",
		Reloadable:   true,
		ConnTimeout:  time.Duration(time.Second * 3000),
		PostData:     ctx.GetDataBox().DataFilePath,
	})
}

func parseReqParamFunc(ctx *Context) {
	logger.Info("parseReqParamFunc start ", ctx.GetDataBox().GetId())

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
		Rule:         "suppredict",
		TransferType: request.NONETYPE,
		Method:       "GET",
		Reloadable:   true,
	}

	ctx.GetDataBox().SetParam("demMemberId", commonRequestData.PubReqInfo.MemId)
	ctx.GetDataBox().SetParam("jobId", commonRequestData.PubReqInfo.JobId)
	ctx.GetDataBox().SetParam("serialNo", commonRequestData.PubReqInfo.SerialNo)

	ctx.AddChanQueue(dataReq)
}

func supPredictFunc(ctx *Context) {

	jobId := ctx.GetDataBox().Param("jobId")

	// 根据jobid获取orderroute map
	orPolicyMap, ok := or.OrderRoutePolicyMap[jobId]
	if !ok {
		logger.Error("[supPredictFunc] get order route policy by jobId failed")
		return
	}

	supMemberId := orPolicyMap.Calllist[0]
	memberDetailInfo, err := member.GetPartnerInfoById(supMemberId)
	if err != nil {
		logger.Error("[supPredictFunc] get partner info by memberid [%s] error: [%s]", supMemberId, err.Error())
		return
	}

	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")
	header.Set("prdtIdCd", ctx.GetDataBox().Param("prdtIdCd"))
	header.Set("serialNo", ctx.GetDataBox().Param("serialNo"))
	header.Set("busiSerialNo", ctx.GetDataBox().Param("busiSerialNo"))
	header.Set("jobId", ctx.GetDataBox().Param("jobId"))

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "suppredictresponse",
		Method:       "POSTBODY",
		TransferType: request.NONETYPE,
		//Url:          "http://127.0.0.1:8096/api/drcs/serverPredictCreditScore",
		Url:          memberDetailInfo.SvrURL,
		Reloadable:   true,
		HeaderArgs:   header,
		Parameters:   ctx.GetDataBox().HttpRequestBody,
		ConnTimeout:  time.Duration(time.Second * 3000),
	})
}

func supPredictResponseFunc(ctx *Context) {

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[supPredictResponseFunc] resultmsg encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().BodyChan <- ctx.DataResponse.Body


	demMemberId := ctx.GetDataBox().Param("demMemberId")
	busiSerialNo := ctx.GetDataBox().Param("busiSerialNo")
	start := ctx.GetDataBox().Param("startTime")
	startTime, err := strconv.Atoi(start)
	if err != nil {
		logger.Error("[supPredictResponseFunc] convert startTime string to int err: [%s] ", err.Error())
		errEnd(ctx)
		return
	}
	endTime := int(time.Now().UnixNano() / 1e6)

	h := md5.New()
	h.Write([]byte(demMemberId))
	busiInfoStr := fmt.Sprintf("%x", h.Sum(nil))

	msg := "3" + demMemberId + "1"
	priKey, _ := security.GetPrivateKey()
	signInfo := cncrypt.Sign(priKey, []byte(msg))
	stepInfoM := []map[string]interface{}{}
	stepInfo1 := map[string]interface{}{"no": 1, "memID": demMemberId, "stepStatus": string(security.StepStatusSucc), "signature": ""}
	stepInfo2 := map[string]interface{}{"no": 2, "memID": "", "stepStatus": string(security.StepStatusSucc), "signature": ""}
	stepInfo3 := map[string]interface{}{"no": 3, "memID": demMemberId, "stepStatus": string(security.StepStatusSucc), "signature": signInfo}
	stepInfoM = append(stepInfoM, stepInfo1)
	stepInfoM = append(stepInfoM, stepInfo2)
	stepInfoM = append(stepInfoM, stepInfo3)

	ctx.Output(map[string]interface{}{
		"exID":       busiInfoStr,
		"demMemID":   demMemberId,
		"supMemID":   ctx.GetDataBox().Param("supMemId"),
		"taskID":     strings.Replace(ctx.GetDataBox().Param("taskIdStr"), "|@|", ".", -1),
		"seqNo":      busiSerialNo,
		"recordType": string(RecordTypeSingle),
		"succCount":  "1",
		"flowStatus": string(FlowStatusDemSucc),
		"usedTime":   endTime - startTime,
		"errCode":    ErrCodeSucc,
		"stepInfoM":  stepInfoM,
		//"dmpSeqNo":   "",
	})

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

	procEndFunc(ctx)
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
