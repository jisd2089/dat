package transfer

/**
    Author: luzequan
    Created: 2017-12-29 15:03:01
*/
type SftpTransfer struct {}

func NewSftpTransfer() *SftpTransfer {
	return new(SftpTransfer)
}

// 封装sftp服务
func (ft *SftpTransfer) ExecuteMethod(req Request) Response {



	return nil
}
