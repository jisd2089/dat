package batchdistribution

/**
    Author: luzequan
    Created: 2018-06-25 18:08:23
*/
import (
	"path"
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"drcs/common/sftp"
	"fmt"
	"strings"
	"time"
	"drcs/dep/or"
	"drcs/dep/order"
	"drcs/dep/util"
	"strconv"
)

func init() {
	BATCH_DEM_SEND.Register()
}

var batchRequestInfo *BatchRequest

var BATCH_DEM_SEND = &DataBox{
	Name:        "batch_dem_send",
	Description: "batch_dem_send",
	RuleTree: &RuleTree{
		Root: batchDemSendRootFunc,

		Trunk: map[string]*Rule{
			"pullReqFile": {
				ParseFunc: pullRequestFileFunc,
			},
			"getPolicy": {
				ParseFunc: getRoutePolicyByJobId,
			},
			"execPattern": {
				ParseFunc: execTwoSidePatternFunc,
			},
			"execMultiPattern": {
				ParseFunc: execMultiPatternFunc,
			},
			"singleRoutePre": {
				ParseFunc: singleRouteSendPreFunc,
			},
			"singleRouteSend": {
				ParseFunc: singleRouteSendFunc,
			},
			"staticRoutePreSend": {
				ParseFunc: staticRouteSendPreFunc,
			},
			"staticRouteSend": {
				ParseFunc: staticRouteSendFunc,
			},
			"sendRecord": {
				ParseFunc: sendRecordFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func batchDemSendRootFunc(ctx *Context) {
	fmt.Println("batchDemSendRootFunc ...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "pullReqFile",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func pullRequestFileFunc(ctx *Context) {
	fmt.Println("pullRequestFileFunc ...")

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

	//fmt.Println("NodeAddress: %s", ctx.GetDataBox().GetNodeAddress())
	ctx.AddQueue(&request.DataRequest{
		Method:       "GET",
		//TransferType: request.NONETYPE, // TEST
		TransferType: request.SFTP,
		FileCatalog: fileCatalog,
		Reloadable:  true,
		Rule:        "getPolicy",
	})
}

func getRoutePolicyByJobId(ctx *Context) {
	fmt.Println("getRoutePolicyByJobId ...")

	jobId := ctx.GetDataBox().Param("jobId")

	// 根据jobid获取orderroute map
	orPolicyMap, ok := or.OrderRoutePolicyMap[jobId]
	if !ok {
		errEnd(ctx)
		return
	}

	var nextRule string

	switch orPolicyMap.Batch {
	case 1: // 需方供方点对点模式
		nextRule = "execPattern"
	case 2: // 单需方多供方模式
		nextRule = "execMultiPattern"
	}

	ctx.AddQueue(&request.DataRequest{
		Rule:         nextRule,
		TransferType: request.NONETYPE,
		Priority:     1,
		Reloadable:   true,
	})
}

// 需方供方点对点模式
func execTwoSidePatternFunc(ctx *Context) {
	fmt.Println("execTwoSidePatternFunc ...")

	jobId := ctx.GetDataBox().Param("jobId")
	// 根据jobid获取orderroute map
	orPolicyMap, ok := or.OrderRoutePolicyMap[jobId]
	if !ok {
		errEnd(ctx)
		return
	}

	if err := getBatchRequest(ctx); err != nil {
		errEnd(ctx)
		return
	}

	var nextRule string

	switch orPolicyMap.RouteMethod {
	case or.ROUTE_SINGLE:
		nextRule = "singleRoutePre"
	case or.ROUTE_STATIC:
		nextRule = "staticRoutePreSend"
	case or.ROUTE_DYNAMIC:
		//dynamic route
		fmt.Println("dynamic route policy not implement")
		errEnd(ctx)
		return
	case or.ROUTE_BROADCAST:
		//broadcast rtoute
		fmt.Println("broadcast route policy not implement")
		errEnd(ctx)
		return
	default:
		fmt.Println("route methold is illegal")
		errEnd(ctx)
		return
	}

	ctx.AddQueue(&request.DataRequest{
		Rule:         nextRule,
		TransferType: request.NONETYPE,
		Priority:     1,
		Reloadable:   true,
	})
}

// 单需方多供方模式
func execMultiPatternFunc(ctx *Context) {
	fmt.Println("execMultiPatternFunc ...")

}

func singleRouteSendPreFunc(ctx *Context) {
	fmt.Println("singleRouteSendPreFunc ...")

	genMD5beforeSend(ctx, "singleRouteSend", batchRequestInfo)
}

func staticRouteSendPreFunc(ctx *Context) {
	fmt.Println("staticRouteSendPreFunc ...")

	genMD5beforeSend(ctx, "staticRouteSend", batchRequestInfo)
}

func singleRouteSendFunc(ctx *Context) {
	fmt.Println("singleRouteSendFunc ...")

	jobId := ctx.GetDataBox().Param("jobId")
	// 根据jobid获取orderroute map
	orPolicyMap, ok := or.OrderRoutePolicyMap[jobId]
	if !ok {
		errEnd(ctx)
		return
	}

	svcUrls, _ := getMemberUrls(orPolicyMap.MemTaskIdMap)

	targetUrl := svcUrls[0]
	//ctx.GetDataBox().DetailCount = 1

	dataRequest := &request.DataRequest{
		Rule: "sendRecord",
		//TransferType: request.NONETYPE, // TEST
		TransferType: request.FASTHTTP,
		Url:          targetUrl,
		Method:       "FILESTREAM",
		Priority:     1,
		PostData:     ctx.GetDataBox().DataFilePath,
		Reloadable:   true,
		//CommandParams: relyDatas,
	}

	dataRequest.SetParam("seqNo", batchRequestInfo.SeqNo)
	dataRequest.SetParam("taskId", batchRequestInfo.TaskIdStr)
	dataRequest.SetParam("orderId", batchRequestInfo.JobId)
	dataRequest.SetParam("userId", batchRequestInfo.UserId)
	dataRequest.SetParam("idType", batchRequestInfo.IdType)
	dataRequest.SetParam("dataRange", batchRequestInfo.DataRange)
	dataRequest.SetParam("maxDelay", string(batchRequestInfo.MaxDelay))
	dataRequest.SetParam("md5", batchRequestInfo.MD5)
	dataRequest.SetParam("boxName", "batch_sup_rcv")

	ctx.AddQueue(dataRequest)
}

func staticRouteSendFunc(ctx *Context) {
	fmt.Println("staticRouteSendFunc ...")

	jobId := ctx.GetDataBox().Param("jobId")
	// 根据jobid获取orderroute map
	orPolicyMap, ok := or.OrderRoutePolicyMap[jobId]
	if !ok {
		errEnd(ctx)
		return
	}

	svcUrls, supMemId := getMemberUrls(orPolicyMap.MemTaskIdMap)

	//ctx.GetDataBox().DetailCount = len(svcUrls)

	for i, targetUrl := range svcUrls {
		fmt.Println(targetUrl)
		dataRequest := &request.DataRequest{
			Rule: "sendRecord",
			//TransferType: request.NONETYPE, // TEST
			TransferType: request.FASTHTTP,
			Url:          targetUrl,
			//Url:        "http://127.0.0.1:8096/api/rcv/batch", // TEST
			Method:     "FILESTREAM",
			Priority:   1,
			PostData:   ctx.GetDataBox().DataFilePath,
			Reloadable: true,
		}

		dataRequest.SetParam("seqNo", batchRequestInfo.SeqNo)
		dataRequest.SetParam("taskId", batchRequestInfo.TaskIdStr)
		dataRequest.SetParam("orderId", batchRequestInfo.JobId)
		dataRequest.SetParam("userId", batchRequestInfo.UserId)
		dataRequest.SetParam("idType", batchRequestInfo.IdType)
		dataRequest.SetParam("dataRange", batchRequestInfo.DataRange)
		dataRequest.SetParam("maxDelay", string(batchRequestInfo.MaxDelay))
		dataRequest.SetParam("md5", batchRequestInfo.MD5)
		dataRequest.SetParam("boxName", "batch_sup_rcv")

		dataRequest.SetParam("targetMemberId", supMemId[i])

		ctx.AddQueue(dataRequest)
	}
}

func sendRecordFunc(ctx *Context) {
	fmt.Println("sendRecordFunc ...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		errEnd(ctx)
		return
	}

	fmt.Println(ctx.DataRequest.Param("targetMemberId"))

	ctx.GetDataBox().TsfSuccCount ++

	errCode := "031008"
	if batchRequestInfo.LineCount == 0 {
		errCode = "031009"
	}

	succNumStr := strconv.Itoa(batchRequestInfo.LineCount)
	if batchRequestInfo.LineCount > 0 {
		for i := 1; i < len(batchRequestInfo.TaskIdList); i++ {
			succNumStr += "." + succNumStr
		}
	}

	ctx.Output(map[string]interface{}{
		"exID":       "",
		"demMemID":   batchRequestInfo.UserId,
		"supMemID":   ctx.DataRequest.Param("targetMemberId"),
		"taskID":     strings.Join(batchRequestInfo.TaskIdList, "."),
		"seqNo":      batchRequestInfo.SeqNo,
		"dmpSeqNo":   "",
		"recordType": "2",
		"succCount":  succNumStr,
		"flowStatus": "01",
		"usedTime":   0,
		"errCode":    errCode,
		//"stepInfoM":  stepInfoM,
	})

	//if ctx.GetDataBox().TsfSuccCount == ctx.GetDataBox().DetailCount {
	//	defer ctx.GetDataBox().SetStatus(status.STOP)
	//	defer ctx.GetDataBox().CloseRequestChan()
	//}
}

func getBatchRequest(ctx *Context) error {

	jobId := ctx.GetDataBox().Param("jobId")
	idType := ctx.GetDataBox().Param("idType")

	// 根据jobid获取orderroute map
	orPolicyMap, ok := or.OrderRoutePolicyMap[jobId]
	if !ok {
		return fmt.Errorf("")
	}

	memTaskIdMap := orPolicyMap.MemTaskIdMap

	// TODO
	taskIdItemStr := ""
	for _, taskIdItemStr = range memTaskIdMap {
		break
	}

	taskIdList := strings.Split(taskIdItemStr, "|@|")
	if len(taskIdList) < 1 {
		return fmt.Errorf("")
	}

	// 根据jobid获取orderinfo
	orderInfoMap := order.GetOrderInfoMap()
	orderInfo, ok := orderInfoMap[jobId]
	if !ok {
		return fmt.Errorf("")
	}
	// 根据第一个taskid获取taskinfo

	taskInfo, ok := orderInfo.TaskInfoMapById[taskIdList[0]]
	if !ok {
		return fmt.Errorf("")
	}

	busiSerialNo := util.SeqUtil{}.GenBusiSerialNo(taskInfo.DemMemId)

	dataRange := ""
	for _, taskId := range taskIdList {
		taskInfo, ok := orderInfo.TaskInfoMapById[taskId]
		if !ok {
			continue
		}
		dataRange += "|@|" + taskInfo.ConnObjId
	}

	batchRequestInfo = &BatchRequest{
		SeqNo:      busiSerialNo,
		TaskIdStr:  taskIdItemStr,
		TaskIdList: taskIdList,
		UserId:     taskInfo.DemMemId,
		JobId:      jobId,
		IdType:     idType,
		DataRange:  dataRange,
		MaxDelay:   orPolicyMap.MaxDelay,
	}

	return nil
}
