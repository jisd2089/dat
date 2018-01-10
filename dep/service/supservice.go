package service

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"dat/core"
	"dat/core/databox"
	"dat/core/interaction/request"
)

/**
    Author: luzequan
    Created: 2018-01-10 13:54:30
*/
type SupService struct {}

func NewSupService() *SupService {
	return &SupService{}
}

func (d *SupService) SendSupToDem(ctx *fasthttp.RequestCtx) {

	fmt.Println("hello data")

	df := assetnode.AssetNodeEntity.GetDataBoxByName("demtest")

	if df == nil {
		fmt.Println("databox is nil!")
	}

	context := databox.GetContext(df, &request.DataRequest{})
	dresp := context.SyncParse("ruleTest3")

	ctx.Response.SetStatusCode(dresp.StatusCode)
}

/**
* 以一对一批量碰撞为例
* 2. 供方前置机接收需方exid单条请求，单批次结束批量文件推送给供方
*/
func (d *SupService) RecDemReqAndPushToSup(ctx *fasthttp.RequestCtx) {

	// 2.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName("suprec")
	if b == nil {
		fmt.Println("databox is nil!")
	}
	// 2.2 执行碰撞rule，同步返回碰撞结果
	// 2.3 碰撞结束，执行推送rule，推送文件至供方
}

/**
* 以一对一批量碰撞为例
* 3. 供方准备好返回文件，发送至需方前置机
*/
func (d *SupService) SupRespSendToDem(ctx *fasthttp.RequestCtx) {

	// 3.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName("supsend")
	if b == nil {
		fmt.Println("databox is nil!")
	}
	// 3.2 setDataBoxQueue
	// 3.3 执行DataBox，通过Sftp传输
}