package transfer

/**
    Author: luzequan
    Created: 2018-01-10 19:34:32
*/
import (
	. "dat/core/interaction/response"
	"strings"
	"fmt"
)

type NoneTransfer struct {}

func NewNoneTransfer() *NoneTransfer {
	return &NoneTransfer{}
}

// 封装NoneType服务
func (ft *NoneTransfer) ExecuteMethod(req Request) Response {

	retCode := "000000"
	if strings.EqualFold(req.GetUrl(), "127.0.0.1/send01") {
		retCode = "021003"
		fmt.Println(req.GetBobject())
	} else {

	}
	//fmt.Println("NoneTransfer: %s", req.GetBobject())
	return &DataResponse{StatusCode: 200, ReturnCode: retCode}
}

func (ft *NoneTransfer) Close() {

}