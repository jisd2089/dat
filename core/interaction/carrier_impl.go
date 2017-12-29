package interaction

/**
    Author: luzequan
    Created: 2017-12-28 14:20:45
*/
import (
	"errors"
	"github.com/valyala/fasthttp"

	"dat/core/dataflow"
	"dat/core/interaction/request"
	"dat/core/interaction/transfer"
)

type Cross struct {
	transfer transfer.Transfer
}

var CrossHandler = &Cross{
}

func (c *Cross) Handle(df *dataflow.DataFlow, cReq *request.DataRequest) *dataflow.Context {
	ctx := dataflow.GetContext(df, cReq)

	var resp *fasthttp.Response
	var err error

	switch cReq.GetDownloaderID() {
	case request.SURF_ID:
		resp, err = c.transfer.ExecuteMethod(cReq)
	}
	//case request.PHANTOM_ID:
	//	resp, err = self.phantom.Download(cReq)
	//}

	if resp.StatusCode() >= 400 {
		err = errors.New("响应状态 " + resp.Status)
	}

	ctx.SetResponse(resp).SetError(err)

	return ctx
}