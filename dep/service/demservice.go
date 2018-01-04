package service

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"dat/core"
	"dat/core/dataflow"
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

	df := assetnode.AssetNodeEntity.GetDataFlowByName("demtest")

	if df == nil {
		fmt.Println("dataflow is nil!")
	}

	context := dataflow.GetContext(df, &request.DataRequest{})
	dresp := context.SyncParse("ruleTest3")

	ctx.Response.SetStatusCode(dresp.StatusCode)
}