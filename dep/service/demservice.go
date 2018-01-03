package service

import "github.com/valyala/fasthttp"

/**
    Author: luzequan
    Created: 2018-01-02 19:58:54
*/
type DemService struct {}

func NewDemService() *DemService {
	return &DemService{}
}

func (d *DemService) SendDemToSup(ctx *fasthttp.RequestCtx) {

}