package handler

import (
	"github.com/valyala/fasthttp"
	logger "drcs/log"
	"os"
	"io"
	"path/filepath"
	"path"
	"drcs/core"
	"time"
	"strconv"
)

/**
    Author: luzequan
    Created: 2018-05-08 15:19:02
*/
const (
	tmpServerFilePath = "tmp"
)

type CustomerHandler struct{}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{}
}

func (n *CustomerHandler) UploadCSVfile(ctx *fasthttp.RequestCtx) {
	logger.Info("CustomerHandler UploadCSVfile start")

	bodyChan := make(chan []byte)

	ctxFile, err := ctx.FormFile("file")
	if err != nil {
		logger.Error("[CustomerHandler] UploadCSVfile get file from ctx err [%s]", err.Error())
		return
	}

	targetFilePath :=path.Join(tmpServerFilePath, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
	os.MkdirAll(filepath.Clean(targetFilePath), 0644)

	sourcefile, err := ctxFile.Open()
	if err != nil {
		logger.Error("[CustomerHandler] UploadCSVfile open file[%s] err [%s]", ctxFile.Filename, err.Error())
		return
	}
	targetFile, err := os.OpenFile(path.Join(targetFilePath, ctxFile.Filename), os.O_WRONLY|os.O_CREATE, 0644)
	defer targetFile.Close()
	if err != nil {
		logger.Error("[CustomerHandler] UploadCSVfile open file[%s] err [%s]", ctxFile.Filename, err.Error())
		return
	}
	io.Copy(targetFile, sourcefile)

	jobId := string(ctx.FormValue("jobId"))
	if jobId == "" {
		logger.Error("[CustomerHandler] param jobId is nil")
		return
	}

	boxName := "customer_request"
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
	if b == nil {
		logger.Error("[CustomerHandler] databox [%s] is nil!", boxName)
		return
	}

	b.DataFilePath = path.Join(targetFilePath, ctxFile.Filename)

	b.SetParam("processType", "upload")
	b.SetParam("jobId", jobId)

	b.BodyChan = bodyChan
	setDataBoxQueue(b)

	select {
	case body := <-bodyChan:
		ctx.SetBody(body)
		close(bodyChan)
		go func() {
			time.Sleep(time.Microsecond*1000)
			if err := os.RemoveAll(targetFilePath); err != nil {
				logger.Error("[CustomerHandler] remove file path [%s] error [%s]", targetFilePath, err.Error())
			}
		}()
	}
}

func (n *CustomerHandler) Predict(ctx *fasthttp.RequestCtx) {
	logger.Info("CustomerHandler Predict start")
	bodyChan := make(chan []byte)

	boxName := "customer_request"
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

func (n *CustomerHandler) PredictCreditScoreCard(ctx *fasthttp.RequestCtx) {
	logger.Info("CustomerHandler PredictCreditScoreCard start")

	bodyChan := make(chan []byte)

	boxName := "customer_request"
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