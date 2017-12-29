package transfer

import "github.com/valyala/fasthttp"

/**
    Author: luzequan
    Created: 2017-12-29 15:03:01
*/
type SftpTransfer struct {

}

func NewSftpTransfer() *SftpTransfer {
	s := new(SftpTransfer)
	return s
}

// 封装sftp服务
func (ft *SftpTransfer) ExecuteMethod(req *fasthttp.Request) (resp *fasthttp.Response, err error) {


	return nil, nil
}
