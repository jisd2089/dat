package transfer

/**
    Author: luzequan
    Created: 2018-01-13 19:59:16
*/
import (
	. "dat/core/interaction/response"
)

type RedisTransfer struct {}

func NewRedisTransfer() Transfer {
	return &RedisTransfer{}
}

// 封装fasthttp服务
func (ft *RedisTransfer) ExecuteMethod(req Request) Response {
	dataResponse := &DataResponse{}


	return dataResponse
}