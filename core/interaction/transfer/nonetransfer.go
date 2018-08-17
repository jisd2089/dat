package transfer

/**
    Author: luzequan
    Created: 2018-01-10 19:34:32
*/
import (
	. "drcs/core/interaction/response"
	"sync"
)

type NoneTransfer struct {
	lock sync.RWMutex
}

func NewNoneTransfer() *NoneTransfer {
	return &NoneTransfer{}
}

// 封装NoneType服务
func (nt *NoneTransfer) ExecuteMethod(req Request) Response {

	preRule := req.GetPreRuleName()
	if preRule == "" {
		preRule = req.GetRuleName()
	}

	return &DataResponse{
		StatusCode: 200,
		ReturnCode: "000000",
		Bobject:    req.GetBobject(),
		Body:       req.GetParameters(),
		BodyStrs:   req.GetCommandParams(),
		PreRule:    preRule,
	}
}

func (ft *NoneTransfer) Close() {

}
