package crp

/**
    Author: luzequan
    Created: 2018-08-03 17:30:03
*/
import (
	. "drcs/core/databox"
	"drcs/core/interaction/request"
	. "drcs/dep/nodelib/crp/smartsail"
	. "drcs/dep/nodelib/crp/common"
	"strings"
	"github.com/valyala/fasthttp"
	logger "drcs/log"
	"strconv"
	"time"
)

func init() {
	SMARTRESPONSE.Register()
}

var SMARTRESPONSE = &DataBox{
	Name:        "smart_response",
	Description: "smart_response",
	RuleTree: &RuleTree{
		Root: smartResponseRootFunc,

		Trunk: map[string]*Rule{
			"parseparam": {
				ParseFunc: parseResponseParamFunc,
			},
			"rsaencrypt": {
				ParseFunc: rsaEncryptParamFunc,
			},
			"execquery": {
				ParseFunc: querySmartResponseFunc,
			},
			"mockexecquery": {
				ParseFunc: mockQuerySmartResponseFunc,
			},
			"rsadecrypt": {
				ParseFunc: rsaDecryptFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func smartResponseRootFunc(ctx *Context) {
	logger.Info("smartResponseRootFunc start ", ctx.GetDataBox().GetId())

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "parseparam",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func parseResponseParamFunc(ctx *Context) {
	logger.Info("parseResponseParamFunc start ", ctx.GetDataBox().GetId())

	reqBody := ctx.GetDataBox().HttpRequestBody

	busiInfo := map[string]interface{}{}
	err := json.Unmarshal(reqBody, &busiInfo)
	if err != nil {
		logger.Error("[parseResponseParamFunc] unmarshal request body err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	requestData := &RequestData{}
	name, ok := busiInfo["fullName"]
	if !ok {
		logger.Error("[parseResponseParamFunc] request data param [%s] is nil", "fullName")
		errEnd(ctx)
		return
	}
	requestData.Name = name.(string)
	phoneNumber, ok := busiInfo["phoneNumber"]
	if !ok {
		logger.Error("[parseResponseParamFunc] request data param [%s] is nil", "phoneNumber")
		errEnd(ctx)
		return
	}
	requestData.Phone = phoneNumber.(string)
	starttime, ok := busiInfo["starttime"]
	if !ok {
		logger.Error("[parseResponseParamFunc] request data param [%s] is nil", "starttime")
		errEnd(ctx)
		return
	}
	requestData.StartTime = starttime.(string)

	requestDataByte, err := json.Marshal(requestData)
	if err != nil {
		logger.Error("[parseResponseParamFunc] json marshal request data err [%v]", err.Error())
		errEnd(ctx)
		return
	}

	dataReq := &request.DataRequest{
		Rule:         "rsaencrypt",
		Method:       "RSAENCRYPT",
		TransferType: request.ENCRYPT,
		Reloadable:   true,
		Parameters:   requestDataByte,
	}

	dataReq.SetParam("encryptKey", SMARTSAIL_PUBLIC_KEY)

	ctx.AddChanQueue(dataReq)
}

func rsaEncryptParamFunc(ctx *Context) {
	logger.Info("rsaEncryptParamFunc start ", ctx.GetDataBox().GetId())

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[rsaEncryptParamFunc] base64 encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")

	requestMsg := &RequestMsg{}
	requestMsg.CliKey = SMARTSAIL_CLIKEY
	requestMsg.RequestData = ctx.DataResponse.BodyStr

	requestMsgByte, err := json.Marshal(requestMsg)
	if err != nil {
		logger.Error("[rsaEncryptParamFunc] json Marshal uriData failed [%s]", err.Error())
		errEnd(ctx)
		return
	}

	args := make(map[string]string, 0)
	args["data"] = string(requestMsgByte)

	var (
		nextRule     string
		transferType string
	)

	switch ctx.GetDataBox().Param("guardFlag") {
	case "00":
		nextRule = "execquery"
		transferType = request.FASTHTTP
	case "01":
		nextRule = "mockexecquery"
		transferType = request.NONETYPE
	default:
		nextRule = "execquery"
		transferType = request.FASTHTTP
	}

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         nextRule,
		Method:       "POSTARGS",
		Url:          SMARTSAIL_URL,
		TransferType: transferType,
		Reloadable:   true,
		HeaderArgs:   header,
		PostArgs:     args,
		ConnTimeout:  time.Duration(time.Minute * 30),
		//Url:    "http://10.101.12.43:8088/api/sup/resp",
	})
}

func querySmartResponseFunc(ctx *Context) {
	logger.Info("querySmartResponseFunc start ", ctx.GetDataBox().GetId())

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[querySmartResponseFunc] execute query failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	responseData := &ResponseData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		logger.Error("[querySmartResponseFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	if responseData.RespCode != SMARTSAIL_SUCC {
		logger.Error("[querySmartResponseFunc] smartsail execute query response [%s]", responseData.RespMessage)

		pubRespMsg := &PubResProductMsg{}

		pubAnsInfo := &PubAnsInfo{}
		pubAnsInfo.ResCode = GetCenterCodeFromSMARTSAIL(responseData.RespCode)
		pubAnsInfo.ResMsg = responseData.RespMessage
		pubAnsInfo.SerialNo = ctx.GetDataBox().Param("serialNo")
		pubAnsInfo.BusiSerialNo = ctx.GetDataBox().Param("busiSerialNo")
		pubAnsInfo.TimeStamp = strconv.Itoa(int(time.Now().UnixNano() / 1e6))

		pubRespMsg.PubAnsInfo = pubAnsInfo

		pubRespMsgByte, err := json.Marshal(pubRespMsg)
		if err != nil {
			errEnd(ctx)
			return
		}

		select {
		case <-ctx.GetDataBox().StopChan:
		case ctx.GetDataBox().BodyChan <- pubRespMsgByte:
		}

		//ctx.GetDataBox().BodyChan <- pubRespMsgByte

		procEndFunc(ctx)
		return
	}

	//ctx.GetDataBox().SetParam("resCode", GetCenterCodeFromSMARTSAIL(responseData.RespCode))

	dataRequest := &request.DataRequest{
		Rule:         "rsadecrypt",
		Method:       "RSADECRYPT",
		TransferType: request.ENCRYPT,
		Reloadable:   true,
		PostData:     responseData.RespDetail,
		ConnTimeout:  time.Duration(time.Minute * 30),
	}

	dataRequest.SetParam("encryptKey", SMARTSAIL_PRIVATE_KEY)

	ctx.AddChanQueue(dataRequest)
}

func mockQuerySmartResponseFunc(ctx *Context) {
	logger.Info("mockQuerySmartResponseFunc start ", ctx.GetDataBox().GetId())

	//ctx.DataResponse.Body = []byte(`{"code":0,"msg":"系统正常","data":"ySApoEWkw0dMfRIk8vGV4ufnnR9ojNHUsR0PSyuxD39WVP/XLujQm8W130BqUw/yAb1hodRf8PK7iy+OyXCAlQJ+y960nIsKcvwvP2oaAVfTbe/cu2J4s3eeO0GroghY0VhMSJfTP2VKcrOu6EpbaJHDZpQ83y3XCjmB1SH9KGSgjapVpEiON/nG4I5Nb4a4rCcsgntH6CyWjOsabvbYlx6Ix5HYhqGL96KCNPwRmGpce9bAlaK/5/UBIKocSvCYog1kDUl9g39eT68F+oPNmD0U7p8WyxDFoyUkcweXL9mp1yOfnpXUZdpVGosM+qrwsfNeVTCGydX0PAkXEq3jGg=="}`)
	ctx.DataResponse.Body = []byte(`{
	"code": 200,
	"msg": null,
	"reqTime": "2018-09-06 14:25:20",
	"code_message": "successful result",
	"var_detail": [{
		"mobile": "13868185986",
		"consume_12m_cnt": 2,
		"earliest_cons_m": 28,
		"latest_record_m": 6,
		"consvariety": 0.7401520569,
		"cons_type_amtrate": 0.5015,
		"cons_type_cntrate": 0.3333,
		"consume_12m_amt": 28.7,
		"month_num": 2,
		"mon_max_record": 19.9,
		"max_interval_month": 11,
		"cos_stab": 0.3867595819,
		"label_24_num": 4,
		"close_12m_mean_money": 8.8,
		"level_12m_consume": 1,
		"child_6m_money": 0.0,
		"child_6m_num": 0,
		"hosehold_6m_money": 0.0,
		"hosehold_6m_num": 0,
		"virtual_6m_money": 0.0,
		"type_12m_sum": 2,
		"car_6m_money": 0.0,
		"car_6m_num": 0,
		"diamond_6m_money": 0.0,
		"diamond_6m_num": 0,
		"sports_6m_money": 0.0,
		"sports_6m_num": 0,
		"entertainment_6m_money": 0.0,
		"entertainment_6m_num": 0,
		"digital_6m_money": 0.0,
		"digital_6m_num": 0,
		"virtual_6m_num": 0
	}],
	"error_msg": null
}`)

	responseData := &ResponseDecryptData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		logger.Error("[mockQuerySmartResponseFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	pubRespMsg := &PubResProductMsg{}

	pubAnsInfo := &PubAnsInfo{}
	pubAnsInfo.ResCode = CenterCodeMockSucc
	pubAnsInfo.ResMsg = "挡板返回成功"
	pubAnsInfo.SerialNo = ctx.GetDataBox().Param("serialNo")
	pubAnsInfo.BusiSerialNo = ctx.GetDataBox().Param("busiSerialNo")
	pubAnsInfo.TimeStamp = strconv.Itoa(int(time.Now().UnixNano() / 1e6))

	pubRespMsg.PubAnsInfo = pubAnsInfo
	pubRespMsg.DetailInfo = responseData.RespDetail

	pubRespMsgByte, err := json.Marshal(pubRespMsg)
	if err != nil {
		errEnd(ctx)
		return
	}

	select {
	case <-ctx.GetDataBox().StopChan:
	case ctx.GetDataBox().BodyChan <- pubRespMsgByte:
	}
	//ctx.GetDataBox().BodyChan <- pubRespMsgByte

	procEndFunc(ctx)
}

func rsaDecryptFunc(ctx *Context) {
	logger.Info("rsaDecryptFunc start ", ctx.GetDataBox().GetId())

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[rsaDecryptFunc] rsa decrypt err [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	responseData := &ResponseDecryptData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		logger.Error("[rsaDecryptFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	if responseData.RespCode != SMARTSAIL_RESP_SUCC {
		logger.Error("[rsaDecryptFunc] smartsail response no charge [%s]", responseData.RespMessage)

		pubRespMsg := &PubResProductMsg{}

		pubAnsInfo := &PubAnsInfo{}
		pubAnsInfo.ResCode = GetCenterCodeFromSMARTSAIL(responseData.RespCode)
		pubAnsInfo.ResMsg = responseData.RespMessage
		pubAnsInfo.SerialNo = ctx.GetDataBox().Param("serialNo")
		pubAnsInfo.BusiSerialNo = ctx.GetDataBox().Param("busiSerialNo")
		pubAnsInfo.TimeStamp = strconv.Itoa(int(time.Now().UnixNano() / 1e6))

		pubRespMsg.PubAnsInfo = pubAnsInfo

		pubRespMsgByte, err := json.Marshal(pubRespMsg)
		if err != nil {
			errEnd(ctx)
			return
		}

		select {
		case <-ctx.GetDataBox().StopChan:
		case ctx.GetDataBox().BodyChan <- pubRespMsgByte:
		}
		//ctx.GetDataBox().BodyChan <- pubRespMsgByte

		procEndFunc(ctx)
		return
	}

	// 请求真实供方 成功返回
	pubRespMsg := &PubResProductMsg{}

	pubAnsInfo := &PubAnsInfo{}
	//pubAnsInfo.ResCode = ctx.GetDataBox().Param("resCode")
	pubAnsInfo.ResCode = GetCenterCodeFromSMARTSAIL(responseData.RespCode)
	pubAnsInfo.ResMsg = "成功"
	pubAnsInfo.SerialNo = ctx.GetDataBox().Param("serialNo")
	pubAnsInfo.BusiSerialNo = ctx.GetDataBox().Param("busiSerialNo")
	pubAnsInfo.TimeStamp = strconv.Itoa(int(time.Now().UnixNano() / 1e6))

	pubRespMsg.PubAnsInfo = pubAnsInfo
	pubRespMsg.DetailInfo = responseData.RespDetail

	pubRespMsgByte, err := json.Marshal(pubRespMsg)
	if err != nil {
		logger.Error("[rsaDecryptFunc] json marshal PubResProductMsg err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	select {
	case <-ctx.GetDataBox().StopChan:
	case ctx.GetDataBox().BodyChan <- pubRespMsgByte:
	}
	//ctx.GetDataBox().BodyChan <- pubRespMsgByte

	procEndFunc(ctx)
}
