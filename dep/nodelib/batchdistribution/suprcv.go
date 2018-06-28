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
	"crypto/md5"
	"os"
	"io"
	"bufio"
	"bytes"
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
				ParseFunc: pingRedisFunc,
			},
			"pushToServer": {
				ParseFunc: pushToServerFunc,
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

	dataFilePath := ctx.GetDataBox().DataFilePath

	dataFile, err := os.Open(dataFilePath)
	defer dataFile.Close()
	if err != nil {
		return
	}

	buf := bufio.NewReader(dataFile)

	md5Hash := md5.New()
	lineCnt := 300
	cntBuf := &bytes.Buffer{}
	c := lineCnt
	for {
		c--
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF && cntBuf.Len() == 0 {
				break
			} else if err == io.EOF && cntBuf.Len() > 0 {
				md5Hash.Write(cntBuf.Bytes())
				break
			} else {
				errEnd(ctx)
				return
			}
		}

		cntBuf.Write(line)
		cntBuf.WriteByte('\n')

		if c == 0 {

			md5Hash.Write(cntBuf.Bytes())

			c = lineCnt
			cntBuf.Reset()
		}
	}

	md5Str := fmt.Sprintf("%x", md5Hash.Sum(nil))

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

func pushToServerFunc(ctx *Context) {
	fmt.Println("pushToServerFunc ...")

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
		TransferType: request.SFTP,
		FileCatalog:  fileCatalog,
		Reloadable:   true,
	})
}

func sendRcvRecordFunc(ctx *Context) {
	fmt.Println("sendRcvRecordFunc ...")

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
			"exID":       line,
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
	}

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()

}
