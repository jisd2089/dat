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
	tmpCustomerFilePath = "D:/dds_accept/tmp"
)

type DrcsServerHandler struct{}

func NewServerHandler() *DrcsServerHandler {
	return &DrcsServerHandler{}
}

func (n *DrcsServerHandler) AcceptCSVfile(ctx *fasthttp.RequestCtx) {
	logger.Info("DrcsServerHandler AcceptCSVfile start")

	bodyChan := make(chan []byte)

	acpFile, err := ctx.FormFile("file")
	if err != nil {
		logger.Error("[DrcsServerHandler] AcceptCSVfile get file from ctx err [%s]", err.Error())
		return
	}
	acpFileContent, err := acpFile.Open()
	defer acpFileContent.Close()
	if err != nil {
		logger.Error("[DrcsServerHandler] AcceptCSVfile open file[%s] err [%s]", acpFile.Filename, err.Error())
		return
	}
	targetFilePath := tmpCustomerFilePath+"/"+strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	os.MkdirAll(filepath.Clean(targetFilePath), 0644)

	targetFile, err := os.OpenFile(path.Join(targetFilePath, acpFile.Filename), os.O_WRONLY|os.O_CREATE, 0644)
	defer targetFile.Close()
	if err != nil {
		logger.Error("[DrcsServerHandler] AcceptCSVfile open file[%s] err [%s]", acpFile.Filename, err.Error())
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

func (n *DrcsServerHandler) PredictCreditScore(ctx *fasthttp.RequestCtx) {
	logger.Info("DrcsServerHandler PredictCreditScore start")
	bodyChan := make(chan []byte)
	boxName := "server_response"
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
	if b == nil {
		logger.Error("databox [%s] is nil!", boxName)
		return
	}
	b.SetParam("processType", "api")
	b.HttpRequestBody = ctx.Request.Body()
	b.BodyChan = bodyChan
	setDataBoxQueue(b)

	select {
	case body := <-bodyChan:
		ctx.SetBody(body)
		close(bodyChan)
	}

}


func (n *DrcsServerHandler) PredictCreditScoreCard(ctx *fasthttp.RequestCtx) {
	logger.Info("DrcsServerHandler PredictCreditScoreCard start")
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