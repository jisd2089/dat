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
	"github.com/ouqiang/gocron/modules/logger"
	"strconv"
	"time"
)

func init() {
	SMARTRESPONSE.Register()
}

const (
	SMARTSAIL_CLIKEY = ""
)

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
	logger.Info("smartResponseRootFunc start")

	ctx.AddQueue(&request.DataRequest{
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

	dataReq := &request.DataRequest{
		Rule:         "rsaencrypt",
		Method:       "RSAENCRYPT",
		TransferType: request.ENCRYPT,
		Reloadable:   true,
		Parameters:   requestDataByte,
	}

	encryptKey := "0102030405060708"

	dataReq.SetParam("encryptKey", encryptKey)

	ctx.AddQueue(dataReq)
}

func rsaEncryptParamFunc(ctx *Context) {
	logger.Info("rsaEncryptParamFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[rsaEncryptParamFunc] rsa encrypt failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")

	args := make(map[string]string, 0)
	args["cliKey"] = SMARTSAIL_CLIKEY
	args["data"] = string(ctx.DataResponse.Body)

	ctx.AddQueue(&request.DataRequest{
		Rule:         "execquery",
		Method:       "POSTARGS",
		TransferType: request.FASTHTTP,
		Url:          "",
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

	dataRequest := &request.DataRequest{
		Rule:         "rsadecrypt",
		Method:       "RSADECRYPT",
		TransferType: request.ENCRYPT,
		Reloadable:   true,
		Parameters:   ctx.DataResponse.Body,
	}

	encryptKey := "0102030405060708"
	dataRequest.SetParam("encryptKey", encryptKey)

	ctx.AddQueue(dataRequest)
}

func rsaDecryptFunc(ctx *Context) {
	logger.Info("rsaDecryptFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[rsaDecryptFunc] rsa decrypt err [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	responseData := &ResponseData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		logger.Error("[rsaDecryptFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	// 请求真实供方 成功返回
	if strings.EqualFold(responseData.RespCode, "200") {
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
		return
	}

}
