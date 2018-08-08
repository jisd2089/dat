package crp

/**
    Author: luzequan
    Created: 2018-08-01 17:30:03
*/
import (
	. "drcs/core/databox"
	"drcs/core/interaction/request"
	. "drcs/dep/nodelib/crp/edunwang"
	. "drcs/dep/nodelib/crp/common"
	"encoding/json"
	"strings"
	"strconv"
	"github.com/valyala/fasthttp"
	"time"
	"github.com/ouqiang/gocron/modules/logger"
)

func init() {
	SUPRESPONSE.Register()
}

const (
	EDUN_URL        = "http://api.edunwang.com/test/black_check"
	EDUN_SECRET_KEY = "46d4ead46317428b" // 正式环境
	EDUN_APP_ID     = "422833408034"
	EDUN_SECRET_ID  = "302fab9c7acc4209a328e81c3354"
)

var SUPRESPONSE = &DataBox{
	Name:        "sup_response",
	Description: "sup_response",
	RuleTree: &RuleTree{
		Root: supResponseRootFunc,

		Trunk: map[string]*Rule{
			"parseparam": {
				ParseFunc: parseRespParamFunc,
			},
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

func supResponseRootFunc(ctx *Context) {
	logger.Info("supResponseRootFunc start")

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "parseparam",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func parseRespParamFunc(ctx *Context) {
	logger.Info("parseRespParamFunc start")

	reqBody := ctx.GetDataBox().HttpRequestBody

	busiInfo := map[string]interface{}{}
	err := json.Unmarshal(reqBody, &busiInfo)
	if err != nil {
		logger.Error("[parseRespParamFunc] unmarshal request body err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	requestData := &RequestData{}
	idNum, ok := busiInfo["identityNumber"]
	if !ok {
		logger.Error("[parseRespParamFunc] request data param [%s] is nil", "identityNumber")
		errEnd(ctx)
		return
	}
	requestData.IdNum = idNum.(string)
	name, ok := busiInfo["fullName"]
	if !ok {
		logger.Error("[parseRespParamFunc] request data param [%s] is nil", "fullName")
		errEnd(ctx)
		return
	}
	requestData.Name = name.(string)
	phoneNumber, ok := busiInfo["phoneNumber"]
	if !ok {
		logger.Error("[parseRespParamFunc] request data param [%s] is nil", "phoneNumber")
		errEnd(ctx)
		return
	}
	requestData.PhoneNum = phoneNumber.(string)
	timestampstr, ok := busiInfo["timestamp"]
	if !ok {
		logger.Error("[parseRespParamFunc] request data param [%s] is nil", "timestamp")
		errEnd(ctx)
		return
	}
	timestamp, err := strconv.Atoi(timestampstr.(string))
	if err != nil {
		logger.Error("[parseRespParamFunc] convert timestamp string to int err [%v]", timestampstr)
		errEnd(ctx)
		return
	}
	requestData.TimeStamp = timestamp

	requestDataByte, err := json.Marshal(requestData)
	if err != nil {
		logger.Error("[parseRespParamFunc] json marshal request data err [%v]", err.Error())
		errEnd(ctx)
		return
	}

	dataReq := &request.DataRequest{
		Rule:         "aesencrypt",
		Method:       "AESENCRYPT",
		TransferType: request.ENCRYPT,
		Reloadable:   true,
		Parameters:   requestDataByte,
	}

	//encryptKey := "0102030405060708" // TEST
	encryptKey := EDUN_SECRET_KEY

	dataReq.SetParam("encryptKey", encryptKey)

	ctx.AddChanQueue(dataReq)
}

func aesEncryptParamFunc(ctx *Context) {
	logger.Info("aesEncryptParamFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[aesEncryptParamFunc] ase encrypt failed [%s]", ctx.DataResponse.ReturnMsg)
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

func base64EncodeFunc(ctx *Context) {
	logger.Info("base64EncodeFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[base64EncodeFunc] base64 encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	dataRequest := &request.DataRequest{
		Rule:         "urlencode",
		Method:       "URLENCODE",
		TransferType: request.ENCODE,
		Reloadable:   true,
	}

	dataRequest.SetParam("urlstr", ctx.DataResponse.BodyStr)

	ctx.AddChanQueue(dataRequest)
}

func urlEncodeFunc(ctx *Context) {
	logger.Info("urlEncodeFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[urlEncodeFunc] url encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")

	args := make(map[string]string, 0)
	args["appid"] = EDUN_APP_ID
	args["seq_no"] = ""
	args["secret_id"] = EDUN_SECRET_ID
	args["product_id"] = ""
	args["req_data"] = ctx.DataResponse.BodyStr

	dataRequest := &request.DataRequest{
		Rule:         "execquery",
		Method:       "POSTARGS",
		Url:          EDUN_URL,
		//Url:          "http://api.edunwang.com/test/black_check?appid=xxxx&secret_id=xxxx&seq_no=xxx&product_id=xxx&req_data=xxxx",
		TransferType: request.NONETYPE,
		Reloadable:   true,
		HeaderArgs:   header,
		PostArgs:     args,
	}

	ctx.AddChanQueue(dataRequest)
}

func queryResponseFunc(ctx *Context) {
	logger.Info("queryResponseFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[queryResponseFunc] execute query failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	responseData := &ResponseData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		logger.Error("[queryResponseFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	// TODO mock
	responseData.StatusCode = "100"
	responseData.Message = "null"
	responseData.RspData = ""

	if !strings.EqualFold(responseData.StatusCode, EDUN_SUCC) {
		logger.Error("[queryResponseFunc] edunwang execute query response [%s]", responseData.Message)

		pubRespMsg := &PubResProductMsg_0_000_000{}

		pubAnsInfo := &PubAnsInfo{}
		pubAnsInfo.ResCode = GetCenterCodeFromEdun(responseData.StatusCode)
		pubAnsInfo.ResMsg = responseData.Message
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

	if !strings.EqualFold(responseData.Message, "null") {
		logger.Error("[queryResponseFunc] edunwang execute query response [%s]", responseData.Message)
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().SetParam("resCode", GetCenterCodeFromEdun(responseData.StatusCode))

	dataRequest := &request.DataRequest{
		Rule:         "aesdecrypt",
		Method:       "AESDECRYPT",
		TransferType: request.NONETYPE,
		Reloadable:   true,
		Parameters:   []byte(responseData.RspData),
	}

	ctx.AddChanQueue(dataRequest)
}

func aesDecryptFunc(ctx *Context) {
	logger.Info("aesDecryptFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[aesDecryptFunc] aes decrypt err [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	respData := &RspData{}
	// TODO mock
	respData.Tag = "疑似仿冒包装"
	respData.EvilScore = 77

	if err := json.Unmarshal(ctx.DataResponse.Body, respData); err != nil {
		logger.Error("[aesDecryptFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	pubRespMsg := &PubResProductMsg_0_000_000{}

	pubAnsInfo := &PubAnsInfo{}
	pubAnsInfo.ResCode = ctx.GetDataBox().Param("resCode")
	pubAnsInfo.ResMsg = "成功"
	pubAnsInfo.SerialNo = ctx.GetDataBox().Param("serialNo")
	pubAnsInfo.BusiSerialNo = ctx.GetDataBox().Param("busiSerialNo")
	pubAnsInfo.TimeStamp = strconv.Itoa(int(time.Now().UnixNano() / 1e6))

	pubRespMsg.PubAnsInfo = pubAnsInfo
	pubRespMsg.DetailInfo.Tag = respData.Tag
	pubRespMsg.DetailInfo.EvilScore = respData.EvilScore

	pubRespMsgByte, err := json.Marshal(pubRespMsg)
	if err != nil {
		logger.Error("[aesDecryptFunc] json marshal PubResProductMsg_0_000_000 err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().BodyChan <- pubRespMsgByte

	procEndFunc(ctx)
}
