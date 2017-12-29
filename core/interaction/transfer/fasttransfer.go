package transfer

import (
	"github.com/valyala/fasthttp"
)

/**
    Author: luzequan
    Created: 2017-12-28 15:41:49
*/

type FastTransfer struct {

}

func NewFastTransfer() *FastTransfer {
	s := new(FastTransfer)
	return s
}

// 封装fasthttp服务
func (ft *FastTransfer) ExecuteMethod(req *fasthttp.Request) (resp *fasthttp.Response, err error) {


	return nil, nil
}
