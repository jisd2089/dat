package transfer

import (
	"github.com/valyala/fasthttp"
	"fmt"
	. "dat/core/interaction/response"
)

/**
    Author: luzequan
    Created: 2017-12-28 15:41:49
*/

type FastTransfer struct {}

func NewFastTransfer() Transfer {
	return new(FastTransfer)
}

// 封装fasthttp服务
func (ft *FastTransfer) ExecuteMethod(req Request) Response {
	fmt.Println("execute fasthttp")
	freq := fasthttp.AcquireRequest()
	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(freq)
	defer fasthttp.ReleaseResponse(fresp)

	//resp = &fasthttp.Response{}

	//resp.SetStatusCode(200)

	return &DataResponse{
		StatusCode: 200,
	}
}
