package transfer

/**
    Author: luzequan
    Created: 2018-01-10 19:34:32
*/
import (
	. "drcs/core/interaction/response"
	"sync"
	"time"
)

type NoneTransfer struct {
	lock sync.RWMutex
}

func NewNoneTransfer() *NoneTransfer {
	return &NoneTransfer{}
}

// 封装NoneType服务
func (nt *NoneTransfer) ExecuteMethod(req Request) Response {

	//retCode := "000000"
	//if strings.EqualFold(req.GetUrl(), "127.0.0.1/send01") {
	//	retCode = "021003"
	//	fmt.Println(req.GetBobject())
	//} else {
	//
	//}
	//fmt.Println("NoneTransfer: %s", req.GetBobject())
	return &DataResponse{StatusCode: 200, ReturnCode: "000000"}
}

func (ft *NoneTransfer) Close() {

}