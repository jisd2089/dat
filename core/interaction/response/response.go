package response

import (
	"github.com/valyala/fasthttp"
)

/**
    Author: luzequan
    Created: 2018-01-02 14:52:09
*/

type DataResponse struct {
	DataBoxName string                   // 规则名
	Header      *fasthttp.ResponseHeader // response头
	Body        []byte                   // 返回消息体
	BodyStr     string                   // 消息体字符串
	BodyStrs    []string                 // 消息体字符串数组
	StatusCode  int                      // 返回码
	ReturnCode  string                   // 业务返回码
	ReturnMsg   string                   // 业务返回信息
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

func (dr *DataResponse) GetBodyStr() string {
	return dr.BodyStr
}

func (dr *DataResponse) GetBodyStrs() []string {
	return dr.BodyStrs
}
