package service

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"dat/core"
	"sync"
	//"dat/core/databox"
	"dat/dep/management/entity"
	"dat/dep/management/constant"
	"strconv"
	"encoding/json"
)

/**
    Author: luzequan
    Created: 2018-01-10 13:54:30
*/

type SupService struct {
	lock sync.RWMutex
}

func NewSupService() *SupService {
	return &SupService{}
}

/**********************************************************************************************
* 以一对一批量碰撞为例
* 2. 供方前置机接收需方exid单条请求，单批次结束批量文件推送给供方
*
* 2.1 匹配相应的DataBox
* 2.2 执行碰撞rule，同步返回碰撞结果
* 2.3 碰撞结束，执行推送rule，推送文件至供方
*
* 1.1) 接收到start请求后，实例化一个DataBox单例
* 1.2) 初始化， 启动DataBox，
* 2) 接收normal请求，DataBox处理
* 3.1) 接收end请求，DataBox处理
* 3.2) 关闭DataBox
*/
func (s *SupService) RecDemReqAndPushToSup(ctx *fasthttp.RequestCtx) {

	requestData := ctx.Request.Body()
	batchReqestVo := &entity.BatchReqestVo{}
	if err := json.Unmarshal(requestData, batchReqestVo); err != nil {
		return
	}
	pairDataBoxId := batchReqestVo.DataBoxId
	activeDataBoxName := "suprec" + "_" + strconv.Itoa(pairDataBoxId)

	//fmt.Println("reqType: ", batchReqestVo.ReqType)

	switch batchReqestVo.ReqType {
	case constant.ReqType_Start:
		s.lock.RLock()
		defer s.lock.RUnlock()
		fmt.Println("start activeDataBoxName: ***************", activeDataBoxName)
		b := assetnode.AssetNodeEntity.GetDataBoxByName("suprec")
		if b == nil {
			fmt.Println("databox is nil!")
		}
		var wg sync.WaitGroup
		wg.Add(1)
		b.StartWG = &wg
		b.PairDataBoxId = pairDataBoxId

		assetnode.AssetNodeEntity.PushActiveDataBox(b)
		wg.Wait()
		//fmt.Println("waitgroup end")

		ab := assetnode.AssetNodeEntity.GetActiveDataBoxByName(activeDataBoxName)
		//fmt.Println("active databox name", ab.Name)
		dataResp := assetnode.AssetNodeEntity.RunActiveBox(ab, batchReqestVo)
		//fmt.Println("dataResp:", dataResp)
		ctx.SetStatusCode(dataResp.StatusCode)
	case constant.ReqType_Normal:
		fmt.Println("rec ReqType_Normal req")
		//ab := assetnode.AssetNodeEntity.GetActiveDataBoxByName(activeDataBoxName)
		//fmt.Println("activeDataBoxName: ***************", activeDataBoxName)
		//fmt.Println("active databox name", ab.Name)
		//dataResp := assetnode.AssetNodeEntity.RunActiveBox(ab, batchReqestVo)
		//fmt.Println("dataResp:", dataResp)

		//ctx.SetStatusCode(dataResp.StatusCode)
		ctx.SetStatusCode(200)

	case constant.ReqType_End:
		fmt.Println("end activeDataBoxName: ***************", activeDataBoxName)
		ab := assetnode.AssetNodeEntity.GetActiveDataBoxByName(activeDataBoxName)
		dataResp := assetnode.AssetNodeEntity.RunActiveBox(ab, batchReqestVo)
		//fmt.Println("dataResp:", dataResp)
		ctx.SetStatusCode(dataResp.StatusCode)
		//assetnode.AssetNodeEntity.StopActiveBox(ab)
	}

}

/**********************************************************************************************
* 以一对一批量碰撞为例
* 3. 供方准备好返回文件，发送至需方前置机
* (扫描到文件后，调用此服务，将反馈文件路径作参数传入)
*/
func (d *SupService) SupRespSendToDem(ctx *fasthttp.RequestCtx) {

	filePath := string(ctx.FormValue("filePath"))

	// 3.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName("supsend")
	if b == nil {
		fmt.Println("databox is nil!")
	}
	b.DataFilePath = filePath

	// 3.2 setDataBoxQueue
	setDataBoxQueue(b)
	// 3.3 执行DataBox，通过Sftp传输
}


/**********************************************************************************************
 *test
 */
func (d *SupService) SupRespWholeSendToDem(ctx *fasthttp.RequestCtx) {

	filePath := string(ctx.FormValue("filePath"))

	// 3.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName("supsendnotsplit")
	if b == nil {
		fmt.Println("databox is nil!")
	}
	b.DataFilePath = filePath

	// 3.2 setDataBoxQueue
	setDataBoxQueue(b)
	// 3.3 执行DataBox，通过Sftp传输
}