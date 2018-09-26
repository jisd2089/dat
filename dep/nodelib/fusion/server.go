package fusion

import (
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	. "drcs/dep/nodelib/fusion/taifinance"
	. "drcs/dep/nodelib/fusion/common"
	"fmt"
	"time"
	"github.com/valyala/fasthttp"
	"strings"

	logger "drcs/log"
	"os"
	"bufio"
	"io"
	"strconv"
	"path"
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
			"execuploadsuccess": {
				ParseFunc: uploadTFSuccessFunc,
			},
			"execpredict": {
				ParseFunc: execPredictFunc,
			},
			"predictresponse": {
				ParseFunc: execPredictResponseFunc,
			},
			"execGetProcessedDataSet": {
				ParseFunc: getProcessDataSetSuccessFunc,
			},
			"sendrecord": {
				ParseFunc: sendRecordFunc,
			},
			"parseparam": {
				ParseFunc: parseResponseParamFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func serverRootFunc(ctx *Context) {
	logger.Info("serverRootFunc start ", ctx.GetDataBox().GetId())

	switch ctx.GetDataBox().Param("processType") {
	case "upload":

		ctx.AddChanQueue(&request.DataRequest{
			Rule:         "execuploaddataset",
			Method:       "GET",
			TransferType: request.NONETYPE,
			Reloadable:   true,
			ConnTimeout:  time.Duration(time.Second * 3000),
		})

	case "api":
		ctx.AddChanQueue(&request.DataRequest{
			Rule:         "execpredict",
			Method:       "GET",
			TransferType: request.NONETYPE,
			Reloadable:   true,
			ConnTimeout:  time.Duration(time.Second * 3000),
		})
	default:
		logger.Error("[serverRootFunc] processType [%s] error", ctx.GetDataBox().Param("processType"))
		errEnd(ctx)
	}
}

/**
*******************************************************************************************
*   数据集上传
*******************************************************************************************
 */
func execUploadDataSetFunc(ctx *Context) {
	logger.Info("execUploadDataSetFunc start ", ctx.GetDataBox().GetId())

	// 调用 泰融 上传接口
	header := &fasthttp.RequestHeader{}
	header.Set("tfapi-key", TFAPI_KEY)

	dataRequest := &request.DataRequest{
		Rule:         "execuploadsuccess",
		Method:       "POSTFILE",
		Url:          TFAPI_UPLOAD_URL,
		TransferType: request.FASTHTTP,
		Reloadable:   true,
		HeaderArgs:   header,
		PostData:     ctx.GetDataBox().DataFilePath,
		ConnTimeout:  time.Duration(time.Second * 300),
	}

	ctx.AddChanQueue(dataRequest)
}

func uploadTFSuccessFunc(ctx *Context) {
	logger.Info("uploadTFSuccessFunc start ", ctx.GetDataBox().GetId())

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[uploadTFSuccessFunc] resultmsg encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	// 调用 泰融 获取数据集和特征值接口
	header := &fasthttp.RequestHeader{}
	header.Set("tfapi-key", TFAPI_KEY)
	dataRequest := &request.DataRequest{
		Rule:         "execGetProcessedDataSet",
		Method:       "POSTBODY",
		Url:          TFAPI_PROCESSED_DATASETS_URL,
		TransferType: request.FASTHTTP,
		Reloadable:   true,
		HeaderArgs:   header,
		ConnTimeout:  time.Duration(time.Second * 300),
	}

	dataRequest.SetParam("name", "testNow")
	dataRequest.SetParam("datasetAbsPath", string(ctx.DataResponse.Body))
	dataRequest.SetParam("expansionType", "UNRELATED_ITEM")
	dataRequest.SetParam("modelType", "CREDIT_SCORE")

	ctx.AddChanQueue(dataRequest)
}


//获取泰融数据集和特征值
func getProcessDataSetSuccessFunc(ctx *Context) {
	logger.Info("getProcessDataSetSuccessFunc start ", ctx.GetDataBox().GetId())

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[getProcessDataSetSuccessFunc] resultmsg encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}
	fmt.Println(string(ctx.DataResponse.Body))

//	ctx.DataResponse.Body = []byte(`[{
//	"id": 273,
//	"name": "示例数据集",
//	"description": "",
//	"rowNumber": 2810,
//	"obsStartDate": "2018-08-31T15:57:00",
//	"obsEndDate": "2018-08-31T15:57:59",
//	"behaviorStartDate": "2018-08-31T15:57:00",
//	"behaviorEndDate": "2018-08-31T15:57:59",
//	"datasetAbsPath": "/creditScoreAPI/datasets/88/7d18303b-a380-4eb5-b64e-f30b079c6c5c/demoDataset.csv",
//	"defaultRatio": 0.39466193,
//	"labelIntroduction": "",
//	"datasetIntroduction": "",
//	"creditScoreFeatureCollectionId": 177,
//	"creditScoreFeatureCollectionName": "示例数据集的特征集",
//	"valid": true,
//	"numOfMetadataFeatures": 22,
//	"expansionType": "UNRELATED_ITEM",
//	"modelType": "CREDIT_SCORE"
//}]`)

	ctx.GetDataBox().BodyChan <- ctx.DataResponse.Body

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "sendrecord",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func sendRecordFunc(ctx *Context) {
	logger.Info("sendRecordFunc start ", ctx.GetDataBox().GetId())

	dataFile, err := os.Open(ctx.GetDataBox().DataFilePath)
	defer dataFile.Close()
	if err != nil {
		logger.Error("[sendRecordFunc] open dataset file failed [%s]", err.Error())
		errEnd(ctx)
		return
	}

	buf := bufio.NewReader(dataFile)

	cnt := 0
	for {
		_, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				errEnd(ctx)
				return
			}
		}

		cnt ++
	}

	go os.RemoveAll(path.Dir(ctx.GetDataBox().DataFilePath))

	ctx.Output(map[string]interface{}{
		"exID":       "",
		"demMemID":   ctx.GetDataBox().Param("UserId"),
		"supMemID":   ctx.GetDataBox().Param("NodeMemberId"),
		"taskID":     strings.Replace(ctx.GetDataBox().Param("TaskId"), "|@|", ".", -1),
		"seqNo":      ctx.GetDataBox().Param("seqNo"),
		"dmpSeqNo":   ctx.GetDataBox().Param("fileNo"),
		"recordType": "2",
		"succCount":  strconv.Itoa(cnt),
		"flowStatus": "11",
		"usedTime":   11,
		"errCode":    "031014",
		//"stepInfoM":  stepInfoM,
	})
}

/**
*******************************************************************************************
*   模型api调用
*******************************************************************************************
 */
func execPredictFunc(ctx *Context) {
	logger.Info("execPredictFunc start ", ctx.GetDataBox().GetId())

	header := &fasthttp.RequestHeader{}
	header.Set("tfapi-key", TFAPI_KEY)
	header.SetMethod("GET")

	postArgs := &PredictCreditScoreReq{}
	if err := json.Unmarshal(ctx.GetDataBox().HttpRequestBody, postArgs); err != nil {
		logger.Error("[execPredictFunc] json unmarshal failed [%s]", err.Error())
		errEnd(ctx)
		return
	}

	postArgsMap := make(map[string]string)
	postArgsMap["modelUID"] = postArgs.ModelUID
	postArgsMap["instancesAmount"] = postArgs.InstancesAmount
	postArgsMap["instancesArray"] = postArgs.InstancesArray

	// TODO
	var tfApiUrl string
	switch ctx.GetDataBox().Param("jobId") {
	case "JON20180912000000781":
		tfApiUrl = TFAPI_PREDICT_CREDIT_SCORE_URL
	case "JON20180913000000782":
		tfApiUrl = TFAPI_PREDICT_CREDIT_SCORE_CARD_URL
	}

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "predictresponse",
		Method:       "POSTARGS",
		Url:          tfApiUrl,
		TransferType: request.FASTHTTP,
		Reloadable:   true,
		HeaderArgs:   header,
		PostArgs:     postArgsMap,
		ConnTimeout:  time.Duration(time.Second * 300),
	})
}

func execPredictResponseFunc(ctx *Context) {
	logger.Info("execPredictResponseFunc start ", ctx.GetDataBox().GetId())

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		logger.Error("[execPredictResponseFunc] resultmsg encode failed [%s]", ctx.DataResponse.ReturnMsg)
		errEnd(ctx)
		return
	}

	// mock
//	ctx.DataResponse.Body = []byte(`{
//	"resultCode": 0,
//	"resultMessage": "API执行成功",
//	"resultData": {
//		"defaultProbability": [
//			0.6799222597181479,
//			0.43125540974628596
//		],
//		"creditScore": [338.46999117398, 420.21021112098003]
//	}
//}`)
	//mock-end

	responseData := &ResponseData{}
	if err := json.Unmarshal(ctx.DataResponse.Body, responseData); err != nil {
		logger.Error("[execPredictResponseFunc] json unmarshal response data err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	// 请求真实供方 成功返回
	pubRespMsg := &PubResProductMsg{}

	pubAnsInfo := &PubAnsInfo{}
	pubAnsInfo.ResCode = GetCenterCodeFromTAIFIN(responseData.RespCode)
	pubAnsInfo.ResMsg = "成功"
	pubAnsInfo.SerialNo = ctx.GetDataBox().Param("serialNo")
	pubAnsInfo.BusiSerialNo = ctx.GetDataBox().Param("busiSerialNo")
	pubAnsInfo.TimeStamp = strconv.Itoa(int(time.Now().UnixNano() / 1e6))

	pubRespMsg.PubAnsInfo = pubAnsInfo
	pubRespMsg.DetailInfo = responseData.RespDetail

	pubRespMsgByte, err := json.Marshal(pubRespMsg)
	if err != nil {
		logger.Error("[execPredictResponseFunc] json marshal PubResProductMsg err [%s]", err.Error())
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().BodyChan <- pubRespMsgByte

	procEndFunc(ctx)
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

	requestData := &PredictCreditScoreReq{}
	modelUID, ok := busiInfo["modelUID"]
	if !ok {
		logger.Error("[parseResponseParamFunc] request data param [%s] is nil", "modelUID")
		errEnd(ctx)
		return
	}
	requestData.ModelUID = modelUID.(string)
	instancesAmount, ok := busiInfo["instancesAmount"]
	if !ok {
		logger.Error("[parseResponseParamFunc] request data param [%s] is nil", "instancesAmount")
		errEnd(ctx)
		return
	}
	requestData.InstancesAmount = instancesAmount.(string)
	instancesArray, ok := busiInfo["instancesArray"]
	if !ok {
		logger.Error("[parseResponseParamFunc] request data param [%s] is nil", "instancesArray")
		errEnd(ctx)
		return
	}
	requestData.InstancesArray = instancesArray.(string)

	requestDataByte, err := json.Marshal(requestData)
	if err != nil {
		logger.Error("[parseResponseParamFunc] json marshal request data err [%v]", err.Error())
		errEnd(ctx)
		return
	}

	dataReq := &request.DataRequest{
		Rule:         "execPredictCreditScore",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
		Parameters:   requestDataByte,
	}

	//dataReq.SetParam("encryptKey", SMARTSAIL_PUBLIC_KEY)

	ctx.AddChanQueue(dataReq)
}
