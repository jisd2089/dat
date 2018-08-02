package crp

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"fmt"
	"strings"
	"encoding/json"
	. "drcs/dep/nodelib/crp/edunwang"
	"drcs/common/balance"
	"sync"
	"strconv"
	"os"
	"drcs/dep/or"
	"drcs/dep/order"
)

func init() {
	DEMREQUEST.Register()
}

var lock sync.Mutex

var DEMREQUEST = &DataBox{
	Name:        "dem_request",
	Description: "dem_request",
	RuleTree: &RuleTree{
		Root: demrequestRootFunc,

		Trunk: map[string]*Rule{
			"parseparam": {
				ParseFunc: parseReqParamFunc,
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
				ParseFunc: queryResponseFunc,
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

func demrequestRootFunc(ctx *Context) {
	fmt.Println("demrequest Root ...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "parseparam",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func parseReqParamFunc(ctx *Context) {
	fmt.Println("parseReqParamFunc rule...")

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
	dataReq.SetParam("appkey", ctx.GetDataBox().Param("appkey"))
	dataReq.SetParam("jobId", commonRequestData.PubReqInfo.JobId)

	ctx.AddQueue(dataReq)
}

func depAuthFunc(ctx *Context) {
	fmt.Println("depAuthFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("Authentication failed")
		errEnd(ctx)
		return
	}

	reqData := ctx.DataResponse.Bobject

	reqDataJson, err := json.Marshal(reqData)
	if err != nil {
		fmt.Println("parse reqData failed")
		errEnd(ctx)
		return
	}

	ctx.AddQueue(&request.DataRequest{
		Rule:         "applybalance",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
		Parameters:   reqDataJson,
	})
}

func applyBalanceFunc(ctx *Context) {
	fmt.Println("applyBalance rule...")

	jobId := ctx.GetDataBox().Param("jobId")
	orderRoutePolicy := or.OrderRoutePolicyMap[jobId]
	supMemberId := orderRoutePolicy.Calllist[0]
	taskId := orderRoutePolicy.MemTaskIdMap[supMemberId]

	orderData := order.GetOrderInfoMap()[jobId]

	orderDetailInfo := orderData.TaskInfoMapById[taskId]

	unitPriceStr := orderDetailInfo.ValuationPrice

	ctx.GetDataBox().SetParam("unitPrice", unitPriceStr)

	memberId := ctx.GetDataBox().Param("memberId")
	//unitPriceStr := ctx.GetDataBox().Param("unitPrice")
	balanceUrl := ctx.GetDataBox().Param("balanceUrl")

	unitPrice, err := strconv.ParseFloat(unitPriceStr, 64)
	if err != nil {
		fmt.Println("apply balance failed", err.Error())
		errEnd(ctx)
		return
	}

	lock.Lock()
	defer lock.Unlock()
	hasBalance := balance.Hasbalance(memberId, unitPrice)
	if !hasBalance {
		if err := balance.ApplyBalance(memberId, unitPrice, 100, balanceUrl); err != nil {
			fmt.Println("apply balance failed", err.Error())
			errEnd(ctx)
			return
		}
	}

	dataRequest := &request.DataRequest{
		Rule:         "updateredisquato",
		Method:       "HIncrBy",
		TransferType: request.REDIS,
		Reloadable:   true,
		Parameters:   ctx.DataResponse.Body,
	}

	dataRequest.SetParam("key", strconv.Itoa(os.Getpid()))
	dataRequest.SetParam("field", memberId)
	dataRequest.SetParam("incr", unitPriceStr)

	ctx.AddQueue(dataRequest)
}

func updateRedisQuatoFunc(ctx *Context) {
	fmt.Println("updateRedisQuatoFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("update redis quato failed")
		errEnd(ctx)
		return
	}

	memberId := ctx.GetDataBox().Param("memberId")
	unitPriceStr := ctx.GetDataBox().Param("unitPrice")

	unitPrice, err := strconv.ParseFloat(unitPriceStr, 64)
	if err != nil {
		fmt.Println("apply balance failed", err.Error())
		errEnd(ctx)
		return
	}

	if err := balance.UpdateBalance(memberId, -unitPrice); err != nil {
		fmt.Println("apply balance failed", err.Error())
		errEnd(ctx)
		return
	}

	dataRequest := &request.DataRequest{
		Rule:         "reduceredisquato",
		Method:       "HDecrBy",
		TransferType: request.REDIS,
		Reloadable:   true,
		Parameters:   ctx.DataResponse.Body,
	}

	dataRequest.SetParam("key", strconv.Itoa(os.Getpid()))
	dataRequest.SetParam("field", memberId)
	dataRequest.SetParam("incr", unitPriceStr)

	ctx.AddQueue(dataRequest)
}

func reduceRedisQuatoFunc(ctx *Context) {
	fmt.Println("reduceRedisQuatoFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("update redis quato failed")
		errEnd(ctx)
		return
	}

	jobId := ctx.GetDataBox().Param("jobId")

	// 根据jobid获取orderroute map
	orPolicyMap, ok := or.OrderRoutePolicyMap[jobId]
	if !ok {
		errEnd(ctx)
		return
	}

	var nextRule string

	switch orPolicyMap.RouteMethod {
	case 0:
		nextRule = "singlequery"
	case 1:
		nextRule = "staticquery"
	}

	ctx.AddQueue(&request.DataRequest{
		Rule:         nextRule,
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
		Parameters:   ctx.DataResponse.Body,
	})
}

func singleQueryFunc(ctx *Context) {
	fmt.Println("singleQueryFunc rule...")

	dataRequest := &request.DataRequest{
		Rule:         "queryresponse",
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

func staticQueryFunc(ctx *Context) {
	fmt.Println("staticQueryFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		ctx.AddQueue(&request.DataRequest{
			Rule:         "queryresponse",
			Method:       "POST",
			Url:          "http://api.edunwang.com/test/black_check?appid=xxxx&secret_id=xxxx&seq_no=xxx&product_id=xxx&req_data=xxxx",
			TransferType: request.NONETYPE,
			Reloadable:   true,
		})
		return
	}

	dataRequest := &request.DataRequest{
		Rule:         "staticquery",
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

func queryResponseFunc(ctx *Context) {
	fmt.Println("queryResponseFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("exec edunwang query failed")
		errEnd(ctx)
		return
	}

	responseData := &ResponseData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		fmt.Println("edunwang parse response failed")
		errEnd(ctx)
		return
	}

	if !strings.EqualFold(responseData.StatusCode, "100") {
		fmt.Println("edunwang query response failed")
		errEnd(ctx)
		return
	}

	if !strings.EqualFold(responseData.Message, "null") {
		fmt.Println("edunwang query response err msg", responseData.Message)
		errEnd(ctx)
		return
	}

	dataRequest := &request.DataRequest{
		Rule:         "urlencode",
		Method:       "AESDecrypt",
		TransferType: request.ENCODE,
		Reloadable:   true,
		Parameters:   []byte(responseData.RspData),
	}

	dataRequest.SetParam("urlstr", ctx.DataResponse.BodyStr)

	ctx.AddQueue(dataRequest)
}
