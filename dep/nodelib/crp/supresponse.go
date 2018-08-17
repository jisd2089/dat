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
	"strings"
	"strconv"
	"github.com/valyala/fasthttp"
	"time"
	logger "drcs/log"
	"encoding/json"
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
			"aesencrypt": {
				ParseFunc: aesEncryptParamFunc,
			},
			//"base64encode": {
			//	ParseFunc: base64EncodeFunc,
			//},
			"execquery": {
				ParseFunc: queryResponseFunc,
			},
			//"base64decode": {
			//	ParseFunc: base64DecodeFunc,
			//},
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

	//fmt.Println(string(requestDataByte))

	dataReq := &request.DataRequest{
		Rule:         "aesencrypt",
		Method:       "AESENCRYPT",
		TransferType: request.ENCRYPT,
		Reloadable:   true,
		Parameters:   requestDataByte,
	}

	dataReq.SetParam("encryptKey", EDUN_SECRET_KEY_TEST)
	dataReq.SetParam("iv", EDUN_SECRET_KEY_TEST)

	ctx.AddChanQueue(dataReq)
}

//func aesEncryptParamFunc(ctx *Context) {
//	logger.Info("aesEncryptParamFunc start")
//
//	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
//		logger.Error("[aesEncryptParamFunc] ase encrypt failed [%s]", ctx.DataResponse.ReturnMsg)
//		errEnd(ctx)
//		return
//	}
//
//	ctx.AddChanQueue(&request.DataRequest{
//		Rule:         "base64encode",
//		Method:       "BASE64ENCODE",
//		TransferType: request.ENCODE,
//		Reloadable:   true,
//		Parameters:   ctx.DataResponse.Body,
//	})
//}

func aesEncryptParamFunc(ctx *Context) {
	logger.Info("aesEncryptParamFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[aesEncryptParamFunc] base64 encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	//fmt.Println("aesEncryptParamFunc:", ctx.DataResponse.BodyStr)
	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")

	seqNo := SeqUtil{}.GenSeqNo()

	uriData := &URIData{}
	uriData.Appid = EDUN_APP_ID_TEST
	// seqNo 长度要求15位
	uriData.SeqNo = seqNo
	uriData.SecretId = EDUN_SECRET_ID_TEST
	uriData.ProductId = EDUN_PRODUCT_ID
	uriData.ReqData = ctx.DataResponse.BodyStr

	uriDataByte, err := json.Marshal(uriData)
	if err != nil {
		logger.Error("[aesEncryptParamFunc] json Marshal uriData failed [%s]", err.Error())
		errEnd(ctx)
		return
	}

	dataRequest := &request.DataRequest{
		Rule:         "execquery",
		Method:       "POSTBODY",
		Url:          EDUN_URL_TEST,
		TransferType: request.NONETYPE,
		Reloadable:   true,
		HeaderArgs:   header,
		Parameters:   uriDataByte,
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

	//fmt.Println("response body:", string(ctx.DataResponse.Body))

	responseData := &ResponseData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		logger.Error("[queryResponseFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	// TODO mock
	responseData.StatusCode = "100"
	responseData.Message = "null"
	responseData.RspData = "J0lXQtYtkBmZrkXFAE4QTWJUYEzLJcWyFHAx1VeJ6TI="
	// mock end

	if !strings.EqualFold(responseData.StatusCode, EDUN_SUCC) {
		logger.Error("[queryResponseFunc] edunwang execute query response [%s]", responseData.Message)

		pubRespMsg := &PubResProductMsg{}

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

	//if responseData.Message != "" {
	//	logger.Error("[queryResponseFunc] edunwang execute query response [%s]", responseData.Message)
	//	errEnd(ctx)
	//	return
	//}

	ctx.GetDataBox().SetParam("resCode", GetCenterCodeFromEdun(responseData.StatusCode))

	dataRequest := &request.DataRequest{
		Rule:         "aesdecrypt",
		Method:       "AESDECRYPT",
		TransferType: request.ENCRYPT,
		Reloadable:   true,
		PostData:     responseData.RspData,
	}

	dataRequest.SetParam("encryptKey", EDUN_SECRET_KEY_TEST)
	dataRequest.SetParam("iv", EDUN_SECRET_KEY_TEST)

	ctx.AddChanQueue(dataRequest)
}

//func base64DecodeFunc(ctx *Context) {
//	logger.Info("base64DecodeFunc start")
//
//	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
//		logger.Error("[base64DecodeFunc] base64 decode err [%s]", ctx.DataResponse.ReturnMsg)
//		errEnd(ctx)
//		return
//	}
//
//	dataRequest := &request.DataRequest{
//		Rule:         "aesdecrypt",
//		Method:       "AESDECRYPT",
//		TransferType: request.ENCRYPT,
//		Reloadable:   true,
//		Parameters:   ctx.DataResponse.Body,
//	}
//
//	dataRequest.SetParam("encryptKey", EDUN_SECRET_KEY_TEST)
//	dataRequest.SetParam("iv", EDUN_SECRET_KEY_TEST)
//
//	ctx.AddChanQueue(dataRequest)
//}

func aesDecryptFunc(ctx *Context) {
	logger.Info("aesDecryptFunc start")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[aesDecryptFunc] aes decrypt err [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	respData := &RspData{}
	// TODO mock
	//respData.Tag = "疑似仿冒包装"
	//respData.EvilScore = 77
	//fmt.Println("decrypt response:", string(ctx.DataResponse.Body))

	if err := json.Unmarshal(ctx.DataResponse.Body, respData); err != nil {
		logger.Error("[aesDecryptFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	pubRespMsg := &PubResProductMsg{}

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

	//fmt.Println(string(pubRespMsgByte))

	ctx.GetDataBox().BodyChan <- pubRespMsgByte

	procEndFunc(ctx)
}
