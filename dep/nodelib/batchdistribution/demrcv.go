package batchdistribution

/**
    Author: luzequan
    Created: 2018-06-25 18:08:48
*/
import (
	"path"
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"drcs/common/sftp"
	"drcs/runtime/status"
	"fmt"
	"strings"
	"io"
	//"time"
	"time"
	"os"
	"bufio"
)

func init() {
	BATCH_DEM_RCV.Register()
}

var BATCH_DEM_RCV = &DataBox{
	Name:        "batch_dem_rcv",
	Description: "batch_dem_rcv",
	RuleTree: &RuleTree{
		Root: batchDemRcvRootFunc,

		Trunk: map[string]*Rule{
			"checkMD5": {
				ParseFunc: checkMD5DemFunc,
			},
			"pushToServer": {
				ParseFunc: pushToDemServerFunc,
			},
			"sendRecord": {
				ParseFunc: sendDemRcvRecordFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func batchDemRcvRootFunc(ctx *Context) {
	fmt.Println("batchDemRcvRootFunc ...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "checkMD5",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func checkMD5DemFunc(ctx *Context) {
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
		Rule:         "pushToServer",
		Method:       "PUT",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func pushToDemServerFunc(ctx *Context) {
	fmt.Println("pushToDemServerFunc ...")

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
		Rule:         "sendRecord",
		Method:       "PUT",
		TransferType: request.NONETYPE,
		//TransferType: request.SFTP,
		FileCatalog:  fileCatalog,
		Reloadable:   true,
	})
}

func sendDemRcvRecordFunc(ctx *Context) {
	fmt.Println("sendDemRcvRecordFunc ...")

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
			"supMemID":   "0000140",
			"taskID":     strings.Replace(ctx.GetDataBox().Param("TaskId"), "|@|", ".", -1),
			"seqNo":      ctx.GetDataBox().Param("seqNo"),
			"dmpSeqNo":   "",
			"recordType": "2",
			"succCount":  "0.0.0",
			"flowStatus": "01",
			"usedTime":   11,
			"errCode":    "031008",
			//"stepInfoM":  stepInfoM,
		})
	}

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()

}
