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
    Created: 2018-01-02 19:58:54
*/
type DemService struct {}

func NewDemService() *DemService {
	return &DemService{}
}

/**
  * 以一对一批量碰撞为例
  * 1. 需方exid文件，发往供方前置机
  */
func (d *DemService) SendDemReqToSup(ctx *fasthttp.RequestCtx) {

	filePath := string(ctx.FormValue("filePath"))
	fmt.Println("filePath:" + filePath)
	//filePath := "D:/dds_send/JON20171102000000276_ID010201_20171213175701_0002.TARGET"

	// 1.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName("demsend")
	if b == nil {
		fmt.Println("databox is nil!")
	}
	b.SetDataFilePath(filePath)

	addrs := []*request.NodeAddress{}
	addrs = append(addrs, &request.NodeAddress{MemberId: "000079", IP: "127.0.0.1", Host: "8081", URL: "/api/sup/rec", Priority: 0})
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000108", IP: "127.0.0.1", Host: "8082", URL: "/api/sup/rec", Priority: 1})
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000109", IP: "127.0.0.1", Host: "8083", URL: "/api/sup/rec", Priority: 2})
	//addrs = append(addrs, &request.NodeAddress{MemberId: "000115", IP: "127.0.0.1", Host: "8084", URL: "/api/sup/rec", Priority: 3})

	b.SetNodeAddress(addrs)

	// 1.2 setDataBoxQueue
	setDataBoxQueue(b)

	// 1.3 执行，单条http执行碰撞请求
}

func setDataBoxQueue(box *databox.DataBox) {
	dataBoxs := []*databox.DataBox{}
	dataBoxs = append(dataBoxs, box)
	assetnode.AssetNodeEntity.PushDataBox(dataBoxs)
}

/**
  * 以一对一批量碰撞为例
  * 4. 需方前置机接收到供方返回文件，推送给需方
  */
func (d *DemService) RecSupRespAndPushToDem(ctx *fasthttp.RequestCtx) {

	// 1.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName("demrec")

	if b == nil {
		fmt.Println("databox is nil!")
	}
	// 1.2 setDataBoxQueue
	setDataBoxQueue(b)
	// 1.3 执行DataBox，sftp推送文件，核验
}





func (d *DemService) SendDemToSup(ctx *fasthttp.RequestCtx) {

	fmt.Println("hello data")

	df := assetnode.AssetNodeEntity.GetDataBoxByName("demtest")

	if df == nil {
		fmt.Println("databox is nil!")
	}

	context := databox.GetContext(df, &request.DataRequest{})
	dresp := context.SyncParse("ruleTest3")

	ctx.Response.SetStatusCode(dresp.StatusCode)
}