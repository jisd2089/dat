package crp

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	. "drcs/dep/nodelib/crp/common"
	. "drcs/dep/nodelib/crp/edunwang"
	"fmt"

	"drcs/common/balance"
	"sync"
	"strconv"
	"os"
	"drcs/dep/or"
	"drcs/dep/order"
	"drcs/dep/util"
	"time"
	"crypto/md5"
	"github.com/valyala/fasthttp"
	logger "drcs/log"
	"drcs/dep/security"
	"drcs/common/cncrypt"
	"strings"
	"encoding/json"
	"drcs/dep/member"
)

func init() {
	DEMREQUEST.Register()
}

var (
	lock sync.Mutex
)

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
			"queryedunresponse": {
				ParseFunc: callResponseFunc,
			},
			"querysmartresponse": {
				ParseFunc: callSmartResponseFunc,
			},
			"returnbalance": {
				ParseFunc: returnBalanceFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func demrequestRootFunc(ctx *Context) {
	logger.Info("demrequestRootFunc start")

	start := int(time.Now().UnixNano() / 1e6)

	ctx.GetDataBox().SetParam("startTime", strconv.Itoa(start))

	ctx.AddChanQueue(&request.DataRequest{
		Rule:          "parseparam",
		Method:        "PING",
		TransferType:  request.REDIS,
		Reloadable:    true,
		CommandParams: ctx.GetDataBox().Params,
	})
}

/**
`{
		"pubReqInfo": {
			"timeStamp": "1469613279966",
			"jobId": "JON20180516000000431",
			"reqSign": "5f4d604a00df289b6b90b66e4d0e1be9d43cd236fc018197dd27e01a0f7e8a3c",
			"serialNo": "2201611161916567677531846",
			"memId": "0000162",
			"authMode": "00"
		},
		"busiInfo": {
			"fullName": "高尚",
			"identityNumber": "330123197507134199",
			"phoneNumber": "13211109876",
			"timestamp": "1531479822"
		}
	}`
 */
func parseReqParamFunc(ctx *Context) {
	logger.Info("parseReqParamFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[parseReqParamFunc] ping redis failed: [%s] ", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	reqBody := ctx.GetDataBox().HttpRequestBody

	commonRequestData := &CommonRequestData{}
	err := json.Unmarshal(reqBody, &commonRequestData)
	if err != nil {
		logger.Error("[parseReqParamFunc] unmarshal CommonRequestData err: [%s] ", err.Error())
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

func depAuthFunc(ctx *Context) {
	logger.Info("depAuthFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[depAuthFunc] dep authentication failed: [%s] ", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	reqData := ctx.DataResponse.Bobject

	reqDataJson, err := json.Marshal(reqData)
	if err != nil {
		logger.Error("[depAuthFunc] marshal request data err: [%s] ", err.Error())
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().HttpRequestBody = reqDataJson

	ctx.AddChanQueue(&request.DataRequest{
		Rule: "applybalance",
		//Rule:         "reduceredisquato",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
		//Parameters:   reqDataJson,
		ConnTimeout: time.Duration(time.Second * 300),
	})
}

func applyBalanceFunc(ctx *Context) {
	logger.Info("applyBalanceFunc start")

	// TODO
	jobId := ctx.GetDataBox().Param("jobId")
	orderRoutePolicy := or.OrderRoutePolicyMap[jobId]
	supMemberId := orderRoutePolicy.Calllist[0]
	taskId := orderRoutePolicy.MemTaskIdMap[supMemberId]

	orderData := order.GetOrderInfoMap()[jobId]

	orderDetailInfo := orderData.TaskInfoMapById[taskId]

	unitPriceStr := orderDetailInfo.ValuationPrice

	ctx.GetDataBox().SetParam("unitPrice", unitPriceStr)

	memberId := ctx.GetDataBox().Param("demMemberId")
	//unitPriceStr := ctx.GetDataBox().Param("unitPrice")
	balanceUrl := ctx.GetDataBox().Param("balanceUrl")

	unitPrice, err := strconv.ParseFloat(unitPriceStr, 64)
	if err != nil {
		logger.Error("[applyBalanceFunc] parse unit price string to float64 err: [%s] ", err.Error())
		errEnd(ctx)
		return
	}

	var (
		transType string
		amount    string
	)

	lock.Lock()
	defer lock.Unlock()
	hasBalance := balance.Hasbalance(memberId, unitPrice)
	if !hasBalance {
		applyAmount, err := balance.ApplyBalance(memberId, unitPrice, 100, balanceUrl)
		if err != nil {
			logger.Error("[applyBalanceFunc] apply balance from balance center failed: [%s] ", err.Error())
			errEnd(ctx)
			return
		}

		transType = request.REDIS

		amount = strconv.FormatFloat(applyAmount*1000, 'E', -1, 64)

	} else {
		transType = request.NONETYPE
	}

	dataRequest := &request.DataRequest{
		Rule:         "updateredisquato",
		Method:       "HIncrBy",
		TransferType: transType,
		Reloadable:   true,
		//Parameters:    ctx.DataResponse.Body,
		CommandParams: ctx.GetDataBox().Params,
	}

	dataRequest.SetParam("key", strconv.Itoa(os.Getpid()))
	dataRequest.SetParam("field", memberId)
	dataRequest.SetParam("amount", amount)

	ctx.AddChanQueue(dataRequest)
}

func updateRedisQuatoFunc(ctx *Context) {
	logger.Info("updateRedisQuatoFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[updateRedisQuatoFunc] update quato redis value failed: [%s] ", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	memberId := ctx.GetDataBox().Param("demMemberId")
	unitPriceStr := ctx.GetDataBox().Param("unitPrice")

	unitPrice, err := strconv.ParseFloat(unitPriceStr, 64)
	if err != nil {
		logger.Error("[updateRedisQuatoFunc] parse unit price string to float64 error: [%s] ", err.Error())
		errEnd(ctx)
		return
	}

	if err := balance.UpdateBalance(memberId, -unitPrice); err != nil {
		logger.Error("[updateRedisQuatoFunc] update balance error: [%s]", err.Error())
		errEnd(ctx)
		return
	}

	dataRequest := &request.DataRequest{
		Rule:         "reduceredisquato",
		Method:       "HDecrBy",
		TransferType: request.REDIS,
		Reloadable:   true,
		//Parameters:    ctx.DataResponse.Body,
		CommandParams: ctx.GetDataBox().Params,
	}

	amount := strconv.FormatFloat(unitPrice*1000, 'E', -1, 64)

	dataRequest.SetParam("key", strconv.Itoa(os.Getpid()))
	dataRequest.SetParam("field", memberId)
	dataRequest.SetParam("amount", amount)

	ctx.AddChanQueue(dataRequest)
}

func reduceRedisQuatoFunc(ctx *Context) {
	logger.Info("reduceRedisQuatoFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[updateRedisQuatoFunc] update quato redis value failed: [%s]", ctx.DataResponse.ReturnMsg)
		returnBalanceFunc(ctx)
		return
	}

	jobId := ctx.GetDataBox().Param("jobId")

	// 根据jobid获取orderroute map
	orPolicyMap, ok := or.OrderRoutePolicyMap[jobId]
	if !ok {
		logger.Error("[updateRedisQuatoFunc] get order route policy by jobId failed")
		returnBalanceFunc(ctx)
		return
	}

	supMemId := orPolicyMap.Calllist[0]
	ctx.GetDataBox().SetParam("supMemId", supMemId)

	taskIdStr := orPolicyMap.MemTaskIdMap[supMemId]
	ctx.GetDataBox().SetParam("taskIdStr", taskIdStr)

	demMemberId := ctx.GetDataBox().Param("demMemberId")

	seqUtil := &util.SeqUtil{}
	busiSerialNo := seqUtil.GenBusiSerialNo(demMemberId)
	ctx.GetDataBox().SetParam("busiSerialNo", busiSerialNo)

	var nextRule string
	switch orPolicyMap.RouteMethod {
	case 1:
		nextRule = "singlequery"
	case 2:
		nextRule = "staticquery"
	}

	ctx.AddChanQueue(&request.DataRequest{
		Rule:          nextRule,
		Method:        "GET",
		TransferType:  request.NONETYPE,
		Reloadable:    true,
		CommandParams: orPolicyMap.Calllist,
		PreRule:       "reduceredisquato",
	})
}

func singleQueryFunc(ctx *Context) {
	logger.Info("singleQueryFunc start")

	supMemberId := ctx.DataResponse.BodyStrs[0]
	memberDetailInfo, err := member.GetPartnerInfoById(supMemberId)
	if err != nil {
		logger.Error("[singleQueryFunc] get partner info by memberid [%s] error: [%s]", supMemberId, err.Error())
		returnBalanceFunc(ctx)
		return
	}

	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")
	header.Set("prdtIdCd", ctx.GetDataBox().Param("prdtIdCd"))
	header.Set("serialNo", ctx.GetDataBox().Param("serialNo"))
	header.Set("busiSerialNo", ctx.GetDataBox().Param("busiSerialNo"))

	dataRequest := &request.DataRequest{
		Rule:   "queryedunresponse",
		Method: "POSTBODY",
		//Url:    "http://127.0.0.1:8096/api/crp/sup", "http://10.101.12.43:8097/api/crp/sup"
		Url:          memberDetailInfo.SvrURL,
		TransferType: request.NONETYPE,
		Reloadable:   true,
		HeaderArgs:   header,
		Parameters:   ctx.GetDataBox().HttpRequestBody,
		ConnTimeout:  time.Duration(time.Second * 300),
	}

	ctx.AddChanQueue(dataRequest)
}

func staticQueryFunc(ctx *Context) {
	logger.Info("staticQueryFunc start")

	callList := ctx.DataResponse.BodyStrs

	// 进行static 第一次查询
	if strings.EqualFold(ctx.DataResponse.PreRule, "reduceredisquato") {

		ctx.GetDataBox().RequestIndex = 0

		if err := execQuery(ctx, callList[ctx.GetDataBox().RequestIndex]); err != nil {
			returnBalanceFunc(ctx)
			return
		}

	} else {
		if ctx.DataResponse.StatusCode == 200 && strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
			pubRespMsg := &PubResProductMsg{}
			if err := json.Unmarshal(ctx.DataResponse.Body, pubRespMsg); err != nil {
				logger.Error("[callResponseFunc] unmarshal response body to PubResProductMsg_0_000_000 err: [%s] ", err.Error())
				returnBalanceFunc(ctx)
				return
			}

			// TODO mock
			pubAnsInfo := &PubAnsInfo{}
			pubAnsInfo.ResCode = "000000"
			pubAnsInfo.ResMsg = "成功"
			pubRespMsg.PubAnsInfo = pubAnsInfo
			pubRespMsg.DetailInfo.Tag = "疑似仿冒包装"
			pubRespMsg.DetailInfo.EvilScore = 77
			//ctx.DataResponse.Body, _ = json.Marshal(pubRespMsg)
			//fmt.Println(string(ctx.DataResponse.Body))
			// TODO mock-end

			if strings.EqualFold(pubRespMsg.PubAnsInfo.ResCode, CenterCodeSucc) {
				ctx.AddChanQueue(&request.DataRequest{
					Rule:         "queryedunresponse",
					Method:       "POSTBODY",
					TransferType: request.NONETYPE,
					Reloadable:   true,
					Parameters:   ctx.DataResponse.Body,
				})
				return

			} else {

				ctx.GetDataBox().RequestIndex ++

				if ctx.GetDataBox().RequestIndex > len(callList)-1 {
					returnBalanceFunc(ctx)
					return
				}

				if err := execQuery(ctx, callList[ctx.GetDataBox().RequestIndex]); err != nil {
					returnBalanceFunc(ctx)
					return
				}
			}
		}
	}
}

func execQuery(ctx *Context, supMemberId string) error {
	memberDetailInfo, err := member.GetPartnerInfoById(supMemberId)
	if err != nil {
		logger.Error("[singleQueryFunc] get partner info by memberid [%s] error: [%s]", supMemberId, err.Error())
		return err
	}
	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")
	header.Set("prdtIdCd", ctx.GetDataBox().Param("prdtIdCd"))
	header.Set("serialNo", ctx.GetDataBox().Param("serialNo"))
	header.Set("busiSerialNo", ctx.GetDataBox().Param("busiSerialNo"))

	dataRequest := &request.DataRequest{
		Rule:         "staticquery",
		Method:       "POSTBODY",
		Url:          memberDetailInfo.SvrURL,
		TransferType: request.NONETYPE,
		Reloadable:   true,
		HeaderArgs:   header,
		Parameters:   ctx.GetDataBox().HttpRequestBody,
		ConnTimeout:  time.Duration(time.Second * 300),
		PreRule:      "staticquery",
	}

	ctx.AddChanQueue(dataRequest)

	return nil
}

func callResponseFunc(ctx *Context) {
	logger.Info("callResponseFunc start")

	pubRespMsg := &PubResProductMsg{}
	// TODO mock
	pubAnsInfo := &PubAnsInfo{}
	pubAnsInfo.ResCode = "000000"
	pubAnsInfo.ResMsg = "成功"
	pubRespMsg.PubAnsInfo = pubAnsInfo
	pubRespMsg.DetailInfo.Tag = "疑似仿冒包装"
	pubRespMsg.DetailInfo.EvilScore = 77
	ctx.DataResponse.Body, _ = json.Marshal(pubRespMsg)
	//fmt.Println(string(ctx.DataResponse.Body))
	// TODO mock-end

	//if err := json.Unmarshal(ctx.DataResponse.Body, pubRespMsg); err != nil {
	//	logger.Error("[callResponseFunc] unmarshal response body to PubResProductMsg_0_000_000 err: [%s] ", err.Error())
	//	returnBalanceFunc(ctx)
	//	return
	//}

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

	msg := "3" + demMemberId + "1"
	priKey, _ := security.GetPrivateKey()
	signInfo := cncrypt.Sign(priKey, []byte(msg))
	stepInfoM := []map[string]interface{}{}
	stepInfo1 := map[string]interface{}{"no": 1, "memID": demMemberId, "stepStatus": security.StepStatusSucc, "signature": ""}
	stepInfo2 := map[string]interface{}{"no": 2, "memID": "", "stepStatus": security.StepStatusSucc, "signature": ""}
	stepInfo3 := map[string]interface{}{"no": 3, "memID": demMemberId, "stepStatus": security.StepStatusSucc, "signature": signInfo}
	stepInfoM = append(stepInfoM, stepInfo1)
	stepInfoM = append(stepInfoM, stepInfo2)
	stepInfoM = append(stepInfoM, stepInfo3)

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
		"stepInfoM":  stepInfoM,
		//"dmpSeqNo":   "",
	})

	procEndFunc(ctx)
}

func returnBalanceFunc(ctx *Context) {
	logger.Info("returnBalanceFunc start")

	memberId := ctx.GetDataBox().Param("demMemberId")
	unitPriceStr := ctx.GetDataBox().Param("unitPrice")

	unitPrice, err := strconv.ParseFloat(unitPriceStr, 64)
	if err != nil {
		logger.Error("[returnBalanceFunc] convert balance [%s] string to float64 err: [%s] ", unitPriceStr, err.Error())
		errEnd(ctx)
		return
	}

	if err := balance.UpdateBalance(memberId, unitPrice); err != nil {
		logger.Error("[returnBalanceFunc] update accountId [%s] balance [%f] string to float64 err: [%s] ", memberId, unitPrice, err.Error())
		errEnd(ctx)
		return
	}

	dataRequest := &request.DataRequest{
		Rule:          "end",
		Method:        "HIncrBy",
		TransferType:  request.REDIS,
		Reloadable:    true,
		CommandParams: ctx.GetDataBox().Params,
	}

	amount := strconv.FormatFloat(unitPrice*1000, 'E', -1, 64)

	dataRequest.SetParam("key", strconv.Itoa(os.Getpid()))
	dataRequest.SetParam("field", memberId)
	dataRequest.SetParam("amount", amount)

	ctx.AddChanQueue(dataRequest)
}
