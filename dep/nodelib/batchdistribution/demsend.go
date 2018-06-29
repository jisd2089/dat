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
	"drcs/dep/member"
	"drcs/runtime/status"
	"os"
	"io"
	"crypto/md5"
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

	fmt.Println("NodeAddress: %s", ctx.GetDataBox().GetNodeAddress())
	ctx.AddQueue(&request.DataRequest{
		Rule:         "getPolicy",
		Method:       "GET",
		TransferType: request.NONETYPE, // TEST
		//TransferType: request.SFTP,
		FileCatalog: fileCatalog,
		Reloadable:  true,
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

	getMD5beforeSend(ctx, "singleRouteSend", batchRequestInfo)
}

func staticRouteSendPreFunc(ctx *Context) {
	fmt.Println("staticRouteSendPreFunc ...")

	getMD5beforeSend(ctx, "staticRouteSend", batchRequestInfo)
}

func getMD5beforeSend(ctx *Context, nextRule string, batchRequest *BatchRequest) {

	md5Str, err := getMD5(ctx.GetDataBox().DataFilePath)
	if err != nil {
		errEnd(ctx)
		return
	}

	fmt.Println("rcv md5: ", md5Str)

	batchRequest.MD5 = md5Str

	ctx.AddQueue(&request.DataRequest{
		Rule:         nextRule,
		TransferType: request.NONETYPE,
		Priority:     1,
		Reloadable:   true,
	})
}

func getMD5beforeSendOld(ctx *Context, nextRule string) {
	dataFilePath := ctx.GetDataBox().DataFilePath

	dataFile, err := os.Open(dataFilePath)
	defer dataFile.Close()
	if err != nil {
		errEnd(ctx)
		return
	}

	md5Buf := make([]byte, 1024)
	//md5Buf := make([]byte, 104857600)

	md5Hash := md5.New()
	for {
		_, err := dataFile.Read(md5Buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				errEnd(ctx)
				return
			}
		}

		md5Hash.Write(md5Buf)
	}

	md5Str := fmt.Sprintf("%x", md5Hash.Sum(nil))

	fmt.Println("rcv md5: ", md5Str)

	ctx.GetDataBox().SetParam("md5", md5Str)

	ctx.AddQueue(&request.DataRequest{
		Rule:         nextRule,
		TransferType: request.NONETYPE,
		Priority:     1,
		Reloadable:   true,
	})
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
	dataRequest.SetParam("taskId", batchRequestInfo.TaskId)
	dataRequest.SetParam("orderId", batchRequestInfo.JobId)
	dataRequest.SetParam("userId", batchRequestInfo.UserId)
	dataRequest.SetParam("idType", batchRequestInfo.IdType)
	dataRequest.SetParam("dataRange", batchRequestInfo.DataRange)
	dataRequest.SetParam("maxDelay", string(batchRequestInfo.MaxDelay))
	dataRequest.SetParam("md5", batchRequestInfo.MD5)

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

	svcUrls, _ := getMemberUrls(orPolicyMap.MemTaskIdMap)

	for _, targetUrl := range svcUrls {
		fmt.Println(targetUrl)
		dataRequest := &request.DataRequest{
			Rule: "sendRecord",
			//TransferType: request.NONETYPE, // TEST
			TransferType: request.FASTHTTP,
			//Url:          targetUrl,
			Url:        "http://127.0.0.1:8096/api/rcv/batch",
			Method:     "FILESTREAM",
			Priority:   1,
			PostData:   ctx.GetDataBox().DataFilePath,
			Reloadable: true,
		}

		dataRequest.SetParam("seqNo", batchRequestInfo.SeqNo)
		dataRequest.SetParam("taskId", batchRequestInfo.TaskId)
		dataRequest.SetParam("orderId", batchRequestInfo.JobId)
		dataRequest.SetParam("userId", batchRequestInfo.UserId)
		dataRequest.SetParam("idType", batchRequestInfo.IdType)
		dataRequest.SetParam("dataRange", batchRequestInfo.DataRange)
		dataRequest.SetParam("maxDelay", string(batchRequestInfo.MaxDelay))
		dataRequest.SetParam("md5", batchRequestInfo.MD5)

		ctx.AddQueue(dataRequest)
	}
}

func sendRecordFunc(ctx *Context) {
	fmt.Println("sendRecordFunc ...")

	stepInfoM := []map[string]interface{}{}
	stepInfo1 := map[string]interface{}{"no": 1, "memID": "0000161", "stepStatus": "1", "signature": "407a6871ef5d1bd043322c2c5da35401bf9bf4a0afcaf7b899a57d262ca0f3d39097a7ec8e1da4548b124c7f374c6598da94533b9541549647417f1739aa0630"}
	stepInfo2 := map[string]interface{}{"no": 2, "memID": "0000162", "stepStatus": "1", "signature": "407a6871ef5d1bd043322c2c5da35401bf9bf4a0afcaf7b899a57d262ca0f3d39097a7ec8e1da4548b124c7f374c6598da94533b9541549647417f1739aa0630"}
	stepInfo3 := map[string]interface{}{"no": 3, "memID": "0000163", "stepStatus": "1", "signature": "407a6871ef5d1bd043322c2c5da35401bf9bf4a0afcaf7b899a57d262ca0f3d39097a7ec8e1da4548b124c7f374c6598da94533b9541549647417f1739aa0630"}
	stepInfoM = append(stepInfoM, stepInfo1)
	stepInfoM = append(stepInfoM, stepInfo2)
	stepInfoM = append(stepInfoM, stepInfo3)

	ctx.Output(map[string]interface{}{
		"exID":       "",
		"demMemID":   batchRequestInfo.UserId,
		"supMemID":   "0000140",
		"taskID":     strings.Replace(batchRequestInfo.TaskId, "|@|", ".", -1),
		"seqNo":      batchRequestInfo.SeqNo,
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
		SeqNo:     busiSerialNo,
		TaskId:    taskIdItemStr,
		UserId:    taskInfo.DemMemId,
		JobId:     jobId,
		IdType:    idType,
		DataRange: dataRange,
		MaxDelay:  orPolicyMap.MaxDelay,
	}

	return nil
}

func getMemberUrls(taskInfoMap map[string]string) ([]string, []string) {
	var svcUrls []string
	var supMemId []string

	for k, _ := range taskInfoMap {
		if p, err := member.GetPartnerInfoById(k); err == nil {
			url := p.SvrURL
			if len(url) > 0 {
				svcUrls = append(svcUrls, url)
				supMemId = append(supMemId, k)
			}
		}
	}
	return svcUrls, supMemId
}
