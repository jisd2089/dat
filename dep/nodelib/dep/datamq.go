package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"path"
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"drcs/common/sftp"
	"drcs/runtime/status"
	"fmt"
	"strings"
	"time"
)

func init() {
	DATAMQ.Register()
}

var DATAMQ = &DataBox{
	Name:        "datamq",
	Description: "datamq",
	RuleTree: &RuleTree{
		Root: datamqRootFunc,

		Trunk: map[string]*Rule{
			"start": {
				ParseFunc: procMqEndFunc,
			},
			//"start": {
			//	ParseFunc: startmqFunc,
			//},
			"pull": {
				ParseFunc: processmqPullFunc,
			},
			"end": {
				ParseFunc: procMqEndFunc,
			},
		},
	},
}

func datamqRootFunc(ctx *Context) {
	fmt.Println("demsend Root ...")

	fmt.Println("NodeAddress: %s", ctx.GetDataBox().GetNodeAddress())
	ctx.AddQueue(&request.DataRequest{
		Rule:         "start",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func startmqFunc(ctx *Context) {
	fmt.Println("datasend start rule...")

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
		Rule:   "pull",
		Method: "GET",
		//TransferType: request.NONETYPE,
		TransferType: request.SFTP,
		FileCatalog:  fileCatalog,
		Reloadable:   true,
	})
}

func processmqPullFunc(ctx *Context) {
	fmt.Println("datasend process ...")
	if ctx.DataResponse.StatusCode == 200 && strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		addressList := ctx.GetDataBox().GetNodeAddress()

		addr := addressList[0]

		ctx.AddQueue(&request.DataRequest{
			Rule:   "end",
			Url:    addr.GetUrl(),
			Method: "POSTFILE",
			//TransferType: request.NONETYPE,
			TransferType: request.FASTHTTP,
			Priority:     1,
			PostData:     ctx.GetDataBox().DataFilePath,
			Reloadable:   true,
		})
	} else {
		ctx.AddQueue(&request.DataRequest{
			Rule:         "end",
			TransferType: request.NONETYPE,
			Priority:     1,
			Reloadable:   true,
		})
	}
}

func procMqEndFunc(ctx *Context) {
	fmt.Println("mq end start ...")

	for i := 0; i < 151; i++ {
		stepInfoM := []map[string]interface{}{}
		stepInfo1 := map[string]interface{}{"no": 1, "memID": "0000161", "stepStatus": "1", "signature": "407a6871ef5d1bd043322c2c5da35401bf9bf4a0afcaf7b899a57d262ca0f3d39097a7ec8e1da4548b124c7f374c6598da94533b9541549647417f1739aa0630"}
		stepInfo2 := map[string]interface{}{"no": 2, "memID": "0000162", "stepStatus": "1", "signature": "407a6871ef5d1bd043322c2c5da35401bf9bf4a0afcaf7b899a57d262ca0f3d39097a7ec8e1da4548b124c7f374c6598da94533b9541549647417f1739aa0630"}
		stepInfo3 := map[string]interface{}{"no": 3, "memID": "0000163", "stepStatus": "1", "signature": "407a6871ef5d1bd043322c2c5da35401bf9bf4a0afcaf7b899a57d262ca0f3d39097a7ec8e1da4548b124c7f374c6598da94533b9541549647417f1739aa0630"}
		stepInfoM = append(stepInfoM, stepInfo1)
		stepInfoM = append(stepInfoM, stepInfo2)
		stepInfoM = append(stepInfoM, stepInfo3)

		ctx.Output(map[string]interface{}{
			"exID":       "00001092018050216520530529412619986",
			"demMemID":   "0000109",
			"supMemID":   "0000140",
			"taskID":     "CTN20171220000014000001620000876.CTN20171220000014000001620000877.CTN20171220000014000001620000878",
			"seqNo":      "00001092018050216520530529412619986",
			"dmpSeqNo":   "00001352016111607462087321234567",
			"recordType": "1",
			"succCount":  "0.0.0",
			"flowStatus": "22",
			"usedTime":   11,
			"errCode":    "030002",
			"stepInfoM":  stepInfoM,
		})
	}

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
}
