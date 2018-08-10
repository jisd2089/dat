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
	"encoding/json"
	"strings"
	"github.com/valyala/fasthttp"
	logger "drcs/log"
	"strconv"
	"time"
	"fmt"
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
			"base64encode": {
				ParseFunc: base64EncodeParamFunc,
			},
			"execquery": {
				ParseFunc: querySmartResponseFunc,
			},
			"base64decode": {
				ParseFunc: base64DecodeRespFunc,
			},
			"rsadecrypt": {
				ParseFunc: rsaDecryptFunc,
			},
			//"base64decode": {
			//	ParseFunc: base64DecodeParamFunc,
			//},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func smartResponseRootFunc(ctx *Context) {
	logger.Info("smartResponseRootFunc start")

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "parseparam",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func parseResponseParamFunc(ctx *Context) {
	logger.Info("parseResponseParamFunc start")

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

	fmt.Println("request data:", string(requestDataByte))

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
	logger.Info("rsaEncryptParamFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[rsaEncryptParamFunc] rsa encrypt failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "base64encode",
		Method:       "BASE64ENCODE",
		TransferType: request.ENCODE,
		Reloadable:   true,
		Parameters:   ctx.DataResponse.Body,
	})
}

func base64EncodeParamFunc(ctx *Context) {
	logger.Info("base64EncodeFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[base64EncodeParamFunc] base64 encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	fmt.Println("base64:", ctx.DataResponse.BodyStr)

	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")

	requestMsg := &RequestMsg{}
	requestMsg.CliKey = SMARTSAIL_CLIKEY
	requestMsg.RequestData = ctx.DataResponse.BodyStr

	requestMsgByte, err := json.Marshal(requestMsg)
	if err != nil {
		logger.Error("[base64EncodeParamFunc] json Marshal uriData failed [%s]", err.Error())
		errEnd(ctx)
		return
	}

	args := make(map[string]string, 0)
	args["data"] = string(requestMsgByte)

	fmt.Println("data:", string(requestMsgByte))

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "execquery",
		Method:       "POSTARGS",
		Url:          SMARTSAIL_URL,
		TransferType: request.NONETYPE,
		Reloadable:   true,
		HeaderArgs:   header,
		PostArgs:     args,
	})
}

func querySmartResponseFunc(ctx *Context) {
	logger.Info("querySmartResponseFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[querySmartResponseFunc] execute query failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	ctx.DataResponse.Body = []byte(`{"code":0,"msg":"系统正常","data":"ySApoEWkw0dMfRIk8vGV4ufnnR9ojNHUsR0PSyuxD39WVP/XLujQm8W130BqUw/yAb1hodRf8PK7iy+OyXCAlQJ+y960nIsKcvwvP2oaAVfTbe/cu2J4s3eeO0GroghY0VhMSJfTP2VKcrOu6EpbaJHDZpQ83y3XCjmB1SH9KGSgjapVpEiON/nG4I5Nb4a4rCcsgntH6CyWjOsabvbYlx6Ix5HYhqGL96KCNPwRmGpce9bAlaK/5/UBIKocSvCYog1kDUl9g39eT68F+oPNmD0U7p8WyxDFoyUkcweXL9mp1yOfnpXUZdpVGosM+qrwsfNeVTCGydX0PAkXEq3jGg=="}`)

	fmt.Println("response body:", string(ctx.DataResponse.Body))

	responseData := &ResponseData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		logger.Error("[querySmartResponseFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	if responseData.RespCode != SMARTSAIL_SUCC {
		logger.Error("[querySmartResponseFunc] smartsail execute query response [%s]", responseData.RespMessage)

		pubRespMsg := &PubResProductMsg_0_000_000{}

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

		ctx.GetDataBox().BodyChan <- pubRespMsgByte

		procEndFunc(ctx)
		return
	}

	ctx.GetDataBox().SetParam("resCode", GetCenterCodeFromSMARTSAIL(responseData.RespCode))

	dataRequest := &request.DataRequest{
		Rule:         "base64decode",
		Method:       "BASE64DECODE",
		TransferType: request.ENCODE,
		Reloadable:   true,
		PostData:     responseData.RespDetail,
	}

	ctx.AddChanQueue(dataRequest)
}

func base64DecodeRespFunc(ctx *Context) {
	logger.Info("base64DecodeRespFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[base64DecodeRespFunc] base64 decode err [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	fmt.Println("base64:", string(ctx.DataResponse.Body))

	dataRequest := &request.DataRequest{
		Rule:         "rsadecrypt",
		Method:       "RSADECRYPT",
		TransferType: request.ENCRYPT,
		Reloadable:   true,
		Parameters:   ctx.DataResponse.Body,
	}

	dataRequest.SetParam("encryptKey", SMARTSAIL_PRIVATE_KEY)

	ctx.AddChanQueue(dataRequest)
}

func rsaDecryptFunc(ctx *Context) {
	logger.Info("rsaDecryptFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[rsaDecryptFunc] rsa decrypt err [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	fmt.Println("response body:", string(ctx.DataResponse.Body))

	responseData := &ResponseDecryptData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		logger.Error("[rsaDecryptFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	// 请求真实供方 成功返回
	pubRespMsg := &PubResProductMsg_0_000_000{}

	pubAnsInfo := &PubAnsInfo{}
	pubAnsInfo.ResCode = ctx.GetDataBox().Param("resCode")
	pubAnsInfo.ResMsg = "成功"
	pubAnsInfo.SerialNo = ctx.GetDataBox().Param("serialNo")
	pubAnsInfo.BusiSerialNo = ctx.GetDataBox().Param("busiSerialNo")
	pubAnsInfo.TimeStamp = strconv.Itoa(int(time.Now().UnixNano() / 1e6))

	pubRespMsg.PubAnsInfo = pubAnsInfo
	pubRespMsg.DetailInfo = responseData.RespDetail

	pubRespMsgByte, err := json.Marshal(pubRespMsg)
	if err != nil {
		logger.Error("[rsaDecryptFunc] json marshal PubResProductMsg_0_000_000 err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().BodyChan <- pubRespMsgByte

	procEndFunc(ctx)
}
