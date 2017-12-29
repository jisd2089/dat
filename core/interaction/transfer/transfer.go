package transfer

/**
    Author: luzequan
    Created: 2017-12-28 14:39:09
*/
import (
	"github.com/valyala/fasthttp"
)

type Transfer interface {
	ExecuteMethod(Request) (resp *fasthttp.Response, err error)

}