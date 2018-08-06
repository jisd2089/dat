package crp

/**
    Author: luzequan
    Created: 2018-08-03 17:30:03
*/
import (
	. "drcs/core/databox"
	"fmt"
	"drcs/core/interaction/request"
	."drcs/dep/nodelib/crp/smartsail"
	"encoding/json"
	"strings"
	"github.com/valyala/fasthttp"
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
	fmt.Println("smartResponseRootFunc root...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "parseparam",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func parseResponseParamFunc(ctx *Context) {
	fmt.Println("parseResponseParamFunc rule...")

	reqBody := ctx.GetDataBox().HttpRequestBody

	busiInfo := map[string]interface{}{}
	err := json.Unmarshal(reqBody, &busiInfo)
	if err != nil {
		fmt.Println(err.Error())
		errEnd(ctx)
		return
	}

	fmt.Println(busiInfo)

	requestData := &RequestData{}
	name, ok := busiInfo["fullName"]
	if !ok {
		errEnd(ctx)
		return
	}
	requestData.Name = name.(string)
	phoneNumber, ok := busiInfo["phoneNumber"]
	if !ok {
		errEnd(ctx)
		return
	}
	requestData.Phone = phoneNumber.(string)
	starttime, ok := busiInfo["starttime"]
	if !ok {
		errEnd(ctx)
		return
	}
	requestData.StartTime = starttime.(string)

	requestDataByte, err := json.Marshal(requestData)
	if err != nil {
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

	encryptKey := "0102030405060708"

	dataReq.SetParam("encryptKey", encryptKey)

	ctx.AddQueue(dataReq)
}

func rsaEncryptParamFunc(ctx *Context) {
	fmt.Println("rsaEncryptParamFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("rsa encrypt failed")
		errEnd(ctx)
		return
	}

	header := &fasthttp.RequestHeader{}
	header.SetContentType("application/json;charset=UTF-8")
	header.SetMethod("POST")

	args := make(map[string]string, 0)
	args["cliKey"] = ""
	args["data"] = string(ctx.DataResponse.Body)

	ctx.AddQueue(&request.DataRequest{
		Rule:         "execquery",
		Method:       "POSTARGS",
		TransferType: request.FASTHTTP,
		Url:          "",
		Reloadable:   true,
		HeaderArgs:   header,
		PostArgs:     args,
		//Parameters:   ctx.DataResponse.Body,
	})
}

func querySmartResponseFunc(ctx *Context) {
	fmt.Println("queryResponseFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("exec edunwang query failed")
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

	dataRequest.SetParam("encryptKey", ctx.DataResponse.BodyStr)

	ctx.AddQueue(dataRequest)
}

func rsaDecryptFunc(ctx *Context) {
	fmt.Println("rsaDecryptFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("exec edunwang query failed")
		errEnd(ctx)
		return
	}

	responseData := &ResponseData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		errEnd(ctx)
		return
	}

	respData := &RespDetail{}
	//respData.Tag = "疑似仿冒包装"
	//respData.EvilScore = 77

	//if err := json.Unmarshal(ctx.DataResponse.Body, respData); err != nil {
	//	fmt.Println("convert respData to struct failed")
	//	errEnd(ctx)
	//	return
	//}

	pubRespMsgByte, err := json.Marshal(respData)
	if err != nil {
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().BodyChan <- pubRespMsgByte

	errEnd(ctx)
	//ctx.AddQueue(&request.DataRequest{
	//	Rule:         "buildresp",
	//	Method:       "Get",
	//	TransferType: request.NONETYPE,
	//	Reloadable:   true,
	//	Bobject:      pubRespMsg,
	//})

}
