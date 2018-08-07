package crp

/**
    Author: luzequan
    Created: 2018-08-01 17:30:03
*/
import (
	. "drcs/core/databox"
	"fmt"
	"drcs/core/interaction/request"
	. "drcs/dep/nodelib/crp/edunwang"
	"encoding/json"
	"strings"
	"strconv"
	"github.com/valyala/fasthttp"
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
			//"getorderinfo": {
			//	ParseFunc: depAuthFunc,
			//},
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
	//fmt.Println("supResponseRootFunc root...")

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "parseparam",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func parseRespParamFunc(ctx *Context) {
	//fmt.Println("parseRespParamFunc rule...")

	reqBody := ctx.GetDataBox().HttpRequestBody

	busiInfo := map[string]interface{}{}
	err := json.Unmarshal(reqBody, &busiInfo)
	if err != nil {
		fmt.Println(err.Error())
		errEnd(ctx)
		return
	}

	//fmt.Println(busiInfo)

	requestData := &RequestData{}
	idNum, ok := busiInfo["identityNumber"]
	if !ok {
		errEnd(ctx)
		return
	}
	requestData.IdNum = idNum.(string)
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
	requestData.PhoneNum = phoneNumber.(string)
	timestampstr, ok := busiInfo["timestamp"]
	if !ok {
		errEnd(ctx)
		return
	}
	timestamp, err := strconv.Atoi(timestampstr.(string))
	if err != nil {
		errEnd(ctx)
		return
	}
	requestData.TimeStamp = timestamp

	requestDataByte, err := json.Marshal(requestData)
	if err != nil {
		errEnd(ctx)
		return
	}

	//fmt.Println("requestDataByte: ",  requestDataByte)
	//fmt.Println("requestDataByte: ",  string(requestDataByte))

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
	//fmt.Println("aesEncryptParamFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("aes encrypt failed")
		errEnd(ctx)
		return
	}

	//fmt.Println("aes encrypt response: ", ctx.DataResponse.Body)

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "base64encode",
		Method:       "BASE64ENCODE",
		TransferType: request.ENCODE,
		Reloadable:   true,
		Parameters:   ctx.DataResponse.Body,
	})
}

func base64EncodeFunc(ctx *Context) {
	//fmt.Println("base64EncodeFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("base encode failed")
		errEnd(ctx)
		return
	}

	//fmt.Println("base encode response: ", ctx.DataResponse.BodyStr)

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
	//fmt.Println("urlEncodeFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("url encode failed")
		errEnd(ctx)
		return
	}

	//fmt.Println("url encode response: ", ctx.DataResponse.BodyStr)

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
		Url:          "http://api.edunwang.com/test/black_check?appid=xxxx&secret_id=xxxx&seq_no=xxx&product_id=xxx&req_data=xxxx",
		TransferType: request.NONETYPE,
		Reloadable:   true,
		HeaderArgs:   header,
		PostArgs:     args,
	}

	//dataRequest.SetParam("appid", ctx.DataResponse.BodyStr)
	//dataRequest.SetParam("secret_id", ctx.DataResponse.BodyStr)
	//dataRequest.SetParam("seq_no", ctx.DataResponse.BodyStr)
	//dataRequest.SetParam("product_id", ctx.DataResponse.BodyStr)
	//dataRequest.SetParam("req_data", ctx.DataResponse.BodyStr)

	ctx.AddChanQueue(dataRequest)

}

func queryResponseFunc(ctx *Context) {
	//fmt.Println("queryResponseFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("exec edunwang query failed")
		errEnd(ctx)
		return
	}

	responseData := &ResponseData{}
	//if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
	//	fmt.Println("edunwang parse response failed")
	//	errEnd(ctx)
	//	return
	//}

	responseData.StatusCode = "100"
	responseData.Message = "null"
	responseData.RspData = ""

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
		Rule:         "aesdecrypt",
		Method:       "AESDECRYPT",
		TransferType: request.NONETYPE,
		Reloadable:   true,
		Parameters:   []byte(responseData.RspData),
	}

	//dataRequest.SetParam("urlstr", ctx.DataResponse.BodyStr)

	ctx.AddChanQueue(dataRequest)
}

func aesDecryptFunc(ctx *Context) {
	//fmt.Println("aesDecryptFunc rule...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		fmt.Println("exec edunwang query failed")
		errEnd(ctx)
		return
	}

	respData := &RspData{}
	respData.Tag = "疑似仿冒包装"
	respData.EvilScore = 77

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

	//ctx.GetDataBox().Callback(pubRespMsgByte)

	ctx.GetDataBox().BodyChan <- pubRespMsgByte

	//defer close(ctx.GetDataBox().BodyChan)

	procEndFunc(ctx)

}