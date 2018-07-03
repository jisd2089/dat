package batchdistribution

/**
    Author: luzequan
    Created: 2018-06-25 18:08:34
*/
import (
	"path"
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"drcs/runtime/status"
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
	"time"
	"drcs/common/sftp"
)

func init() {
	BATCH_SUP_RCV.Register()
}

var BATCH_SUP_RCV = &DataBox{
	Name:        "batch_sup_rcv",
	Description: "batch_sup_rcv",
	RuleTree: &RuleTree{
		Root: batchSupRcvRootFunc,

		Trunk: map[string]*Rule{
			"checkMD5": {
				ParseFunc: checkMD5Func,
			},
			"pingRedis": {
				ParseFunc: rcvPingRedisFunc,
			},
			"pushToServer": {
				ParseFunc: pushToServerFunc,
			},
			"saveSeqNo": {
				ParseFunc: saveSeqNoFunc,
			},
			"sendRecord": {
				ParseFunc: sendRcvRecordFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func batchSupRcvRootFunc(ctx *Context) {
	fmt.Println("batchSupRcvRootFunc ...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "checkMD5",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func checkMD5Func(ctx *Context) {
	fmt.Println("batchSupRcvRootFunc ...")

	md5Str, _, err := getMD5(ctx.GetDataBox().DataFilePath)

	if err != nil {
		errEnd(ctx)
		return
	}

	if !strings.EqualFold(ctx.GetDataBox().Param("md5"), md5Str) {
		errEnd(ctx)
		return
	}

	ctx.AddQueue(&request.DataRequest{
		Rule:         "pingRedis",
		Method:       "PUT",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func rcvPingRedisFunc(ctx *Context) {
	fmt.Println("rcvPingRedisFunc ...")

	dr := &request.DataRequest{
		Rule:         "pushToServer",
		Method:       "PING",
		//TransferType: request.REDIS,
		TransferType: request.NONETYPE,
		Reloadable:   true,
		CommandParams: ctx.GetDataBox().Params,
	}

	//dr.SetParam("redisAddrs", ctx.GetDataBox().Param("redisAddr"))

	ctx.AddQueue(dr)
}

func pushToServerFunc(ctx *Context) {
	fmt.Println("pushToServerFunc ...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		errEnd(ctx)
		return
	}

	fsAddress := ctx.GetDataBox().FileServerAddress
	filePath := ctx.GetDataBox().GetDataFilePath()
	dataFile := path.Base(filePath)

	// 1. push local file to hadoop server
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

	ctx.AddQueue(&request.DataRequest{
		Rule:         "saveSeqNo",
		Method:       "PUT",
		TransferType: request.NONETYPE, // TEST
		//TransferType: request.SFTP,
		FileCatalog:  fileCatalog,
		Reloadable:   true,
	})
}

func saveSeqNoFunc(ctx *Context) {
	fmt.Println("sendRcvRecordFunc ...")

	r := &request.DataRequest{
		Rule:         "sendRecord",
		Method:       "HSET_STRING",
		//TransferType: request.REDIS,
		TransferType: request.NONETYPE,
		Reloadable:   true,
	}

	jobId := ctx.GetDataBox().Param("jobId")
	idType := ctx.GetDataBox().Param("idType")
	batchNo := ctx.GetDataBox().Param("batchNo")
	fileNo := ctx.GetDataBox().Param("fileNo")

	key := jobId + "_" + idType + "_" + batchNo + "_" + fileNo

	r.SetParam("key", key)
	r.SetParam("field", "seqNo")
	r.SetParam("value", ctx.GetDataBox().Param("seqNo"))

	ctx.AddQueue(r)
}

func sendRcvRecordFunc(ctx *Context) {
	fmt.Println("sendRcvRecordFunc ...")

	if ctx.DataResponse.StatusCode == 200 && !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		errEnd(ctx)
		return
	}

	taskIdStr := ctx.GetDataBox().Param("TaskId")
	succNumStr := "1"
	taskIdList := strings.Split(taskIdStr, "|@|")
	for i := 1; i < len(taskIdList); i ++ {
		succNumStr += "." + succNumStr
	}

	dataFilePath := ctx.GetDataBox().GetDataFilePath()

	dataFile, err := os.Open(dataFilePath)
	defer dataFile.Close()
	if err != nil {
		errEnd(ctx)
		return
	}

	buf := bufio.NewReader(dataFile)

	cnt := 0

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				errEnd(ctx)
				return
			}
		}

		cnt ++
		if cnt == 1 {
			continue
		}

		ctx.Output(map[string]interface{}{
			"exID":       string(line),
			"demMemID":   ctx.GetDataBox().Param("UserId"),
			"supMemID":   ctx.GetDataBox().Param("NodeMemberId"),
			"taskID":     strings.Replace(ctx.GetDataBox().Param("TaskId"), "|@|", ".", -1),
			"seqNo":      ctx.GetDataBox().Param("seqNo"),
			"dmpSeqNo":   ctx.GetDataBox().Param("fileNo"),
			"recordType": "2",
			"succCount":  succNumStr,
			"flowStatus": "11",
			"usedTime":   11,
			"errCode":    "031014",
			//"stepInfoM":  stepInfoM,
		})
	}

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()

}



