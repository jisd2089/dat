package interaction

import (
	"dat/core/interaction/response"
	"dat/core/interaction/transfer"
	"dat/core/interaction/request"
)

/**
    Author: luzequan
    Created: 2018-01-12 22:10:13
*/
type Reflect struct {
	fastHttpTsf transfer.Transfer
	sftpTsf     transfer.Transfer
	noneTsf     transfer.Transfer
}

var ReflectHandler = &Reflect{
	fastHttpTsf: transfer.NewFastTransfer(),
	sftpTsf:     transfer.NewSftpTransfer(),
	noneTsf:     transfer.NewNoneTransfer(),
}

func (c *Reflect) Handle(cReq *request.DataRequest) *response.DataResponse {

	var resp *response.DataResponse
	//var err error

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

	return resp
}
