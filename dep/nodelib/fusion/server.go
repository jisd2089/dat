package fusion

import (
	logger "drcs/log"
	"strconv"
	"time"
	"drcs/core/interaction/request"
	"strings"
	"github.com/valyala/fasthttp"
	. "drcs/core/databox"
)

/**
    Author: luzequan
    Created: 2018-09-03 11:55:13
*/
func init() {
	SERVER.Register()
}

var SERVER = &DataBox{
	Name:        "server_response",
	Description: "server_response",
	RuleTree: &RuleTree{
		Root: serverRootFunc,

		Trunk: map[string]*Rule{
			"execuploaddataset": {
				ParseFunc: execUploadDataSetFunc,
			},
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

func serverRootFunc(ctx *Context) {
	//logger.Info("serverRootFunc start ", ctx.GetDataBox().GetId())

	// switch

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "execuploaddataset",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func execUploadDataSetFunc(ctx *Context) {
	//logger.Info("execUploadDataSetFunc start ", ctx.GetDataBox().GetId())

	// 调用 泰融 上传接口
}

// api
func parseResponseParamFunc(ctx *Context) {
	//logger.Info("parseResponseParamFunc start ", ctx.GetDataBox().GetId())

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
	//logger.Info("rsaEncryptParamFunc start ", ctx.GetDataBox().GetId())

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

	ctx.AddChanQueue(&request.DataRequest{
		Rule:   "execquery",
		Method: "POSTARGS",
		Url:    "http://10.101.12.43:8088/api/sup/resp",
		//Url:          SMARTSAIL_URL,
		TransferType: request.FASTHTTP,
		Reloadable:   true,
		HeaderArgs:   header,
		PostArgs:     args,
	})
}

func querySmartResponseFunc(ctx *Context) {
	//logger.Info("querySmartResponseFunc start ", ctx.GetDataBox().GetId())

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

		ctx.GetDataBox().BodyChan <- pubRespMsgByte

		procEndFunc(ctx)
		return
	}

	ctx.GetDataBox().SetParam("resCode", GetCenterCodeFromSMARTSAIL(responseData.RespCode))

	dataRequest := &request.DataRequest{
		Rule:         "rsadecrypt",
		Method:       "RSADECRYPT",
		TransferType: request.ENCRYPT,
		Reloadable:   true,
		PostData:     responseData.RespDetail,
	}

	dataRequest.SetParam("encryptKey", SMARTSAIL_PRIVATE_KEY)

	ctx.AddChanQueue(dataRequest)
}

func mockQuerySmartResponseFunc(ctx *Context) {
	//logger.Info("mockQuerySmartResponseFunc start ", ctx.GetDataBox().GetId())

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[mockQuerySmartResponseFunc] execute query failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	ctx.DataResponse.Body = []byte(`{"code":0,"msg":"系统正常","data":"ySApoEWkw0dMfRIk8vGV4ufnnR9ojNHUsR0PSyuxD39WVP/XLujQm8W130BqUw/yAb1hodRf8PK7iy+OyXCAlQJ+y960nIsKcvwvP2oaAVfTbe/cu2J4s3eeO0GroghY0VhMSJfTP2VKcrOu6EpbaJHDZpQ83y3XCjmB1SH9KGSgjapVpEiON/nG4I5Nb4a4rCcsgntH6CyWjOsabvbYlx6Ix5HYhqGL96KCNPwRmGpce9bAlaK/5/UBIKocSvCYog1kDUl9g39eT68F+oPNmD0U7p8WyxDFoyUkcweXL9mp1yOfnpXUZdpVGosM+qrwsfNeVTCGydX0PAkXEq3jGg=="}`)

	responseData := &ResponseData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		logger.Error("[mockQuerySmartResponseFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	if responseData.RespCode != SMARTSAIL_SUCC {
		logger.Error("[mockQuerySmartResponseFunc] smartsail execute query response [%s]", responseData.RespMessage)

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

		ctx.GetDataBox().BodyChan <- pubRespMsgByte

		procEndFunc(ctx)
		return
	}

	ctx.GetDataBox().SetParam("resCode", GetCenterCodeFromSMARTSAIL(responseData.RespCode))

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

func rsaDecryptFunc(ctx *Context) {
	//logger.Info("rsaDecryptFunc start ", ctx.GetDataBox().GetId())

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

	// 请求真实供方 成功返回
	pubRespMsg := &PubResProductMsg{}

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
		logger.Error("[rsaDecryptFunc] json marshal PubResProductMsg err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().BodyChan <- pubRespMsgByte

	procEndFunc(ctx)
}