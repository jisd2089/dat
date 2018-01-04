package response

import (
	"github.com/valyala/fasthttp"
)

/**
    Author: luzequan
    Created: 2018-01-02 14:52:09
*/

type DataResponse struct {
	DataFlowName string                   //规则名
	Header       *fasthttp.ResponseHeader //response头
	Body         []byte                   //返回消息体
	StatusCode   int                      //返回码
}

func (resp *DataResponse) GetHeader() *fasthttp.ResponseHeader {
	return resp.Header
}

func (resp *DataResponse) SetHeader(header *fasthttp.ResponseHeader) *DataResponse {
	resp.Header = header
	return resp
}

func (resp *DataResponse) SetStatusCode(statusCode int) *DataResponse {
	resp.StatusCode = statusCode
	return resp
}

func (resp *DataResponse) GetStatusCode() int {
	return resp.StatusCode
}

func (resp *DataResponse) SetBody(body []byte) *DataResponse {
	resp.Body = body
	return resp
}

func (resp *DataResponse) GetBody() []byte {
	return resp.Body
}
