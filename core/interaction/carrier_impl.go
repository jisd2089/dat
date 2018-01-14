package interaction

/**
    Author: luzequan
    Created: 2017-12-28 14:20:45
*/
import (
	"dat/core/databox"
	"dat/core/interaction/request"
	"dat/core/interaction/transfer"
	"dat/core/interaction/response"
)

type Cross struct {
	fastHttpTsf transfer.Transfer
	sftpTsf     transfer.Transfer
	noneTsf     transfer.Transfer
}

var CrossHandler = &Cross{
	fastHttpTsf: transfer.NewFastTransfer(),
	sftpTsf:     transfer.NewSftpTransfer(),
	noneTsf:     transfer.NewNoneTransfer(),
}

func (c *Cross) Handle(df *databox.DataBox, cReq *request.DataRequest) *databox.Context {
	ctx := databox.GetContext(df, cReq)

	var resp *response.DataResponse
	var err error

	switch cReq.GetTransferType() {
	case request.HTTP:
		resp = c.fastHttpTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.SFTP:
		resp = c.sftpTsf.ExecuteMethod(cReq).(*response.DataResponse)
	case request.NONETYPE:
		resp = c.noneTsf.ExecuteMethod(cReq).(*response.DataResponse)
	}

	if resp.GetStatusCode() >= 400 {
		//err = errors.New("响应状态 " + resp.Status)
	}

	ctx.SetResponse(resp).SetError(err)

	return ctx
}
