package handler

import (
	"github.com/valyala/fasthttp"
	logger "drcs/log"
	"os"
	"io"
	"path/filepath"
	"path"
	"drcs/core"
	"fmt"
	"time"
	"strconv"
)

/**
    Author: luzequan
    Created: 2018-05-08 15:19:02
*/
const (
	tmpCustomerFilePath = "tmp"
)

type ServerHandler struct{}

func NewServerHandler() *ServerHandler {
	return &ServerHandler{}
}

func (n *ServerHandler) AcceptCSVfile(ctx *fasthttp.RequestCtx) {
	logger.Info("ServerHandler AcceptCSVfile start")

	bodyChan := make(chan []byte)

	acpFile, err := ctx.FormFile("file")
	if err != nil {
		logger.Error("[ServerHandler] AcceptCSVfile get file from ctx err [%s]", err.Error())
		return
	}
	acpFileContent, err := acpFile.Open()
	defer acpFileContent.Close()
	if err != nil {
		logger.Error("[ServerHandler] AcceptCSVfile open file[%s] err [%s]", acpFile.Filename, err.Error())
		return
	}

	targetFilePath := path.Join(tmpCustomerFilePath, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
	os.MkdirAll(filepath.Clean(targetFilePath), 0644)

	targetFile, err := os.OpenFile(path.Join(targetFilePath, acpFile.Filename), os.O_WRONLY|os.O_CREATE, 0644)
	defer targetFile.Close()
	if err != nil {
		logger.Error("[ServerHandler] AcceptCSVfile open file[%s] err [%s]", acpFile.Filename, err.Error())
		return
	}
	io.Copy(targetFile, acpFileContent)
	ctx.Request.Body()
	boxName := "server_response"
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
	if b == nil {
		logger.Error("databox [%s] is nil!", boxName)
		return
	}

	b.SetParam("processType", "upload")
	b.DataFilePath = path.Join(targetFilePath, acpFile.Filename)

	b.BodyChan = bodyChan

	setDataBoxQueue(b)

	select {
	case body := <-bodyChan:
		ctx.SetBody(body)
		close(bodyChan)
		go func() {
			time.Sleep(time.Microsecond*1000)
			if err := os.RemoveAll(targetFilePath); err != nil {
				fmt.Println(err.Error())
			}
		}()
	}
}

func (n *ServerHandler) ExecPredict(ctx *fasthttp.RequestCtx) {
	logger.Info("ServerHandler ExecPredict start")

	prdtIdCd := string(ctx.Request.Header.Peek("prdtIdCd"))
	if prdtIdCd == "" {
		logger.Error("[ServerHandler] prdtIdC d is nil!")
		return
	}

	serialNo := string(ctx.Request.Header.Peek("serialNo"))
	if serialNo == "" {
		logger.Error("[ServerHandler] serialNo is nil!")
		return
	}

	busiSerialNo := string(ctx.Request.Header.Peek("busiSerialNo"))
	if busiSerialNo == "" {
		logger.Error("[ServerHandler] busiSerialNo is nil!")
		return
	}

	jobId := string(ctx.Request.Header.Peek("jobId"))
	if jobId == "" {
		logger.Error("[ServerHandler] jobId is nil!")
		return
	}

	boxName := "server_response"
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
	if b == nil {
		logger.Error("[ServerHandler] databox [%s] is nil!", boxName)
		return
	}
	b.SetParam("processType", "api")
	b.SetParam("serialNo", serialNo)
	b.SetParam("prdtIdCd", prdtIdCd)
	b.SetParam("busiSerialNo", busiSerialNo)
	b.SetParam("jobId", jobId)

	b.HttpRequestBody = ctx.Request.Body()

	bodyChan := make(chan []byte)
	b.BodyChan = bodyChan

	setDataBoxQueue(b)

	select {
	case body := <-bodyChan:
		ctx.SetBody(body)
		close(bodyChan)
	}
}

func (n *ServerHandler) PredictCreditScoreCard(ctx *fasthttp.RequestCtx) {
	logger.Info("ServerHandler PredictCreditScoreCard start")
	bodyChan := make(chan []byte)
	boxName := "server_response"
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
	if b == nil {
		logger.Error("databox [%s] is nil!", boxName)
		return
	}
	b.SetParam("processType", "apiCard")
	b.HttpRequestBody = ctx.Request.Body()
	b.BodyChan = bodyChan
	setDataBoxQueue(b)

	select {
	case body := <-bodyChan:
		ctx.SetBody(body)
		close(bodyChan)
	}
}