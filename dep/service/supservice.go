package service

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"dat/core"
	"sync"
)

/**
    Author: luzequan
    Created: 2018-01-10 13:54:30
*/
type SupService struct {}

func NewSupService() *SupService {
	return &SupService{}
}

/**
* 以一对一批量碰撞为例
* 2. 供方前置机接收需方exid单条请求，单批次结束批量文件推送给供方
*/
func (d *SupService) RecDemReqAndPushToSup(ctx *fasthttp.RequestCtx) {



	reqType := string(ctx.FormValue("ReqType"))
	reqParam := string(ctx.FormValue("reqParam"))
	fmt.Println("reqType: ", reqType)

	switch reqType {
	case "start":
		b := assetnode.AssetNodeEntity.GetDataBoxByName("suprec")
		if b == nil {
			fmt.Println("databox is nil!")
		}
		var wg sync.WaitGroup
		wg.Add(1)
		b.WG = &wg

		assetnode.AssetNodeEntity.PushActiveDataBox(b)
		wg.Wait()
		fmt.Println("waitgroup end")
		//defer close(cb.BlockChan)

		ab := assetnode.AssetNodeEntity.GetActiveDataBoxByName("suprec")
		fmt.Println("active databox name", ab.Name)
		dataResp := assetnode.AssetNodeEntity.RunActiveBox(ab, "1111")
		fmt.Println("dataResp:", dataResp)
	case "normal":
		ab := assetnode.AssetNodeEntity.GetActiveDataBoxByName("suprec")
		fmt.Println("active databox name", ab.Name)
		dataResp := assetnode.AssetNodeEntity.RunActiveBox(ab, reqParam)
		fmt.Println("dataResp:", dataResp)

		//defer close(ab.BlockChan)
	case "end":
		b:= assetnode.AssetNodeEntity.GetActiveDataBoxByName("suprec")
		assetnode.AssetNodeEntity.StopActiveBox(b)
	}
	// 2.1 匹配相应的DataBox
	// 2.2 执行碰撞rule，同步返回碰撞结果
	// 2.3 碰撞结束，执行推送rule，推送文件至供方
	//dataResp := assetnode.AssetNodeEntity.SyncRun()
	//fmt.Println(dataResp)


	// 1.1) 接收到start请求后，实例化一个DataBox单例



	//for {
	//	succflag := <- cb.StartSuccChan
	//	fmt.Println(succflag)
	//	if succflag {
	//
	//		close(cb.StartSuccChan)
	//		//ab := assetnode.AssetNodeEntity.GetActiveDataBoxByName("suprec")
	//		//fmt.Println("active databox name", ab.Name)
	//		//assetnode.AssetNodeEntity.RunActiveBox(ab, "1111")
	//
	//		//fmt.Println("active databox name", cb.Name)
	//		//assetnode.AssetNodeEntity.RunActiveBox(cb, "1111")
	//		break
	//	}
	//}


	//ab := assetnode.AssetNodeEntity.GetActiveDataBoxByName("suprec")
	//fmt.Println("active databox name", ab.Name)


	// 1.2) 初始化， 启动DataBox，
	// 2) 接收normal请求，DataBox处理
	//ab := assetnode.AssetNodeEntity.GetActiveDataBoxByName("suprec")
	//assetnode.AssetNodeEntity.RunActiveBox(ab, "1111")

	// 3.1) 接收end请求，DataBox处理
	// 3.2) 关闭DataBox
}

//func setActiveDataBoxQueue(box *databox.DataBox) *databox.DataBox {
//	dataBoxs := []*databox.DataBox{}
//	dataBoxs = append(dataBoxs, box)
//	assetnode.AssetNodeEntity.PushActiveDataBox(dataBoxs)
//	return
//}

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