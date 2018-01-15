package transfer

/**
    Author: luzequan
    Created: 2018-01-13 19:59:16
*/
import (
	. "dat/core/interaction/response"
	"fmt"
)

type RedisTransfer struct {}

func NewRedisTransfer() Transfer {
	return &RedisTransfer{}
}

// 封装fasthttp服务
func (ft *RedisTransfer) ExecuteMethod(req Request) Response {
	fmt.Println("RedisTransfer")

	return  &DataResponse{StatusCode: 200, ReturnCode: "000000"}
}