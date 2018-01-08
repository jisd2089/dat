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