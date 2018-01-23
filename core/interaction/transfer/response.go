package transfer

/**
    Author: luzequan
    Created: 2018-01-02 15:04:17
*/
import (
	"github.com/valyala/fasthttp"
)

type (
	Response interface {
		GetHeader() *fasthttp.ResponseHeader
		//SetHeader(header *fasthttp.ResponseHeader) Response
		GetStatusCode() int
		//SetStatusCode(statusCode int) Response
		GetBody() []byte
		//SetBody(body []byte) Response
		GetBodyStr() string
		GetBodyStrs() []string
	}
)
