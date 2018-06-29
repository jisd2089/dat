package batchdistribution

/**
    Author: luzequan
    Created: 2018-06-25 18:08:42
*/
import (
	"path"
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"drcs/common/sftp"
	"drcs/runtime/status"
	"fmt"
	"time"
	"strings"
	"drcs/dep/or"
)

func init() {
	BATCH_SUP_SEND.Register()
}

var batchResponseInfo *BatchRequest

var BATCH_SUP_SEND = &DataBox{
	Name:        "batch_sup_send",
	Description: "batch_sup_send",
	RuleTree: &RuleTree{
		Root: batchSupSendRootFunc,

		Trunk: map[string]*Rule{
			"pingRedis": {
				ParseFunc: pingRedisFunc,
			},
			"qrySeqNo": {
				ParseFunc: qrySeqNoFunc,
			},
			"pullRespFile": {
				ParseFunc: pullResponseFileFunc,
			},
			"setBatchResp": {
				ParseFunc: setBatchResponseFunc,
			},
			"postBatchRespPre": {
				ParseFunc: getRespDataMD5Func,
			},
			"postBatchResp": {
				ParseFunc: postBatchRespDataFunc,
			},
			"sendRespRecord": {
				ParseFunc: sendRespRecordFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func batchSupSendRootFunc(ctx *Context) {
	fmt.Println("batchSupSendRootFunc ...")
	ctx.AddQueue(&request.DataRequest{
		Rule:         "pingRedis",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func pingRedisFunc(ctx *Context) {
	fmt.Println("pingRedisFunc ...")

	dr := &request.DataRequest{
		Rule:         "qrySeqNo",
		Method:       "PING",
		TransferType: request.REDIS,
		Reloadable:   true,
	}

	dr.SetParam("redisAddrs", "10.101.12.45:6379")

	ctx.AddQueue(dr)
}

func qrySeqNoFunc(ctx *Context) {
	fmt.Println("qrySeqNoFunc ...")
	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		errEnd(ctx)
		return
	}

	r := &request.DataRequest{
		Rule:         "pullRespFile",
		Method:       "HGET_STRING",
		TransferType: request.REDIS,
		Reloadable:   true,
	}

	jobId := ctx.GetDataBox().Param("jobId")
	idType := ctx.GetDataBox().Param("idType")
	batchNo := ctx.GetDataBox().Param("batchNo")
	fileNo := ctx.GetDataBox().Param("fileNo")

	key := jobId + "_" + idType + "_" + batchNo + "_" + fileNo

	r.SetParam("key", key)
	r.SetParam("field", "seqNo")

	ctx.AddQueue(r)
}

func pullResponseFileFunc(ctx *Context) {
	fmt.Println("pullResponseFileFunc ...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		errEnd(ctx)
		return
	}

	ctx.GetDataBox().SetParam("seqNo", ctx.DataResponse.BodyStr)

	filePath := ctx.GetDataBox().GetDataFilePath()
	dataFile := path.Base(filePath)
	dataFilePath := path.Dir(filePath)

	fmt.Println(dataFilePath + "@" + ctx.GetDataBox().GetDataFilePath())

	fsAddress := ctx.GetDataBox().FileServerAddress

	// 1. 从sftp服务器（需方dmp服务器）拉取文件到节点服务器本地
	fileCatalog := &sftp.FileCatalog{
		UserName:       fsAddress.UserName,
		Password:       fsAddress.Password,
		Host:           fsAddress.Host,
		Port:           fsAddress.Port,
		TimeOut:        time.Duration(fsAddress.TimeOut) * time.Second,
		LocalDir:       fsAddress.LocalDir,
		LocalFileName:  dataFile,
		RemoteDir:      fsAddress.RemoteDir,
		RemoteFileName: dataFile,
	}

	ctx.GetDataBox().SetDataFilePath(path.Join(fsAddress.LocalDir, dataFile))

	fmt.Println("NodeAddress: %s", ctx.GetDataBox().GetNodeAddress())
	ctx.AddQueue(&request.DataRequest{
		Rule:         "setBatchResp",
		Method:       "GET",
		TransferType: request.NONETYPE, // TEST
		//TransferType: request.SFTP,
		FileCatalog: fileCatalog,
		Reloadable:  true,
	})
}

func setBatchResponseFunc(ctx *Context) {
	fmt.Println("setBatchResponseFunc ...")

	dmpSeqNo := ctx.GetDataBox().Param("fileNo")

	jobId := ctx.GetDataBox().Param("jobId")

	batchResponseInfo = &BatchRequest{
		SeqNo:     ctx.GetDataBox().Param("seqNo"),
		DmpSeqNo:  dmpSeqNo,
		TaskIdStr:    ctx.GetDataBox().Param("batchNo"),
		JobId:     jobId,
		IdType:    "sup",
		UserId:    ctx.GetDataBox().Param("NodeMemberId"),
		DataRange: "sup",
		MaxDelay:  0,
	}

	ctx.AddQueue(&request.DataRequest{
		Rule:         "postBatchRespPre",
		TransferType: request.NONETYPE,
		Priority:     1,
		Reloadable:   true,
	})
}

func getRespDataMD5Func(ctx *Context) {
	fmt.Println("getRespDataMD5Func ...")

	genMD5beforeSend(ctx, "postBatchResp", batchResponseInfo)
}

func postBatchRespDataFunc(ctx *Context) {
	fmt.Println("postBatchRespDataFunc ...")

	jobId := ctx.GetDataBox().Param("jobId")
	// 根据jobid获取orderroute map
	orPolicyMap, ok := or.OrderRoutePolicyMap[jobId]
	if !ok {
		errEnd(ctx)
		return
	}

	svcUrls, _ := getMemberUrls(orPolicyMap.MemTaskIdMap)

	for _, targetUrl := range svcUrls {
		fmt.Println(targetUrl)
		dataRequest := &request.DataRequest{
			Rule: "sendRespRecord",
			//TransferType: request.NONETYPE, // TEST
			TransferType: request.FASTHTTP,
			//Url:          targetUrl,
			Url:        "http://127.0.0.1:8095/api/rcv/batch",
			Method:     "FILESTREAM",
			Priority:   1,
			PostData:   ctx.GetDataBox().DataFilePath,
			Reloadable: true,
		}

		dataRequest.SetParam("seqNo", batchResponseInfo.SeqNo)
		dataRequest.SetParam("taskId", batchResponseInfo.TaskIdStr)
		dataRequest.SetParam("orderId", batchResponseInfo.JobId)
		dataRequest.SetParam("userId", batchResponseInfo.UserId)
		dataRequest.SetParam("idType", batchResponseInfo.IdType)
		dataRequest.SetParam("dataRange", batchResponseInfo.DataRange)
		dataRequest.SetParam("maxDelay", string(batchResponseInfo.MaxDelay))
		dataRequest.SetParam("md5", batchResponseInfo.MD5)

		ctx.AddQueue(dataRequest)
	}
}

func sendRespRecordFunc(ctx *Context) {
	fmt.Println("sendRespRecordFunc ...")

	stepInfoM := []map[string]interface{}{}
	stepInfo1 := map[string]interface{}{"no": 1, "memID": "0000161", "stepStatus": "1", "signature": "407a6871ef5d1bd043322c2c5da35401bf9bf4a0afcaf7b899a57d262ca0f3d39097a7ec8e1da4548b124c7f374c6598da94533b9541549647417f1739aa0630"}
	stepInfo2 := map[string]interface{}{"no": 2, "memID": "0000162", "stepStatus": "1", "signature": "407a6871ef5d1bd043322c2c5da35401bf9bf4a0afcaf7b899a57d262ca0f3d39097a7ec8e1da4548b124c7f374c6598da94533b9541549647417f1739aa0630"}
	stepInfo3 := map[string]interface{}{"no": 3, "memID": "0000163", "stepStatus": "1", "signature": "407a6871ef5d1bd043322c2c5da35401bf9bf4a0afcaf7b899a57d262ca0f3d39097a7ec8e1da4548b124c7f374c6598da94533b9541549647417f1739aa0630"}
	stepInfoM = append(stepInfoM, stepInfo1)
	stepInfoM = append(stepInfoM, stepInfo2)
	stepInfoM = append(stepInfoM, stepInfo3)

	ctx.Output(map[string]interface{}{
		"exID":       "",
		"demMemID":   batchResponseInfo.UserId,
		"supMemID":   "0000140",
		"taskID":     strings.Replace(batchResponseInfo.TaskIdStr, "|@|", ".", -1),
		"seqNo":      batchResponseInfo.SeqNo,
		"dmpSeqNo":   "",
		"recordType": "2",
		"succCount":  "0.0.0",
		"flowStatus": "01",
		"usedTime":   11,
		"errCode":    "031008",
		//"stepInfoM":  stepInfoM,
	})

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
}
