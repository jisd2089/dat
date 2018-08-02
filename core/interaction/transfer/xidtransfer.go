package transfer

type XidTransfer struct{}

func NewXidTransfer() Transfer {
	return &XidTransfer{}
}

// 封装xid服务
func (xt *XidTransfer) ExecuteMethod(req Request) Response {

	return nil
}

func (xt *XidTransfer) Close() {}