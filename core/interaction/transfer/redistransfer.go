package transfer

/**
    Author: luzequan
    Created: 2018-01-13 19:59:16
*/
import (
	. "dat/core/interaction/response"
	redisLib "dat/common/redis"
	//"gopkg.in/redis.v5"
	"fmt"
	"sync"
)

type RedisTransfer struct {
	redisCli redisLib.RedisClient
}

func NewRedisTransfer() Transfer {
	return &RedisTransfer{}
}

var (
	redOnce sync.Once
)

// 封装fasthttp服务
func (rt *RedisTransfer) ExecuteMethod(req Request) Response {
	fmt.Println("RedisTransfer")

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("redisTransfer recover error: ", err)
		}
	}()

	//rt.connect()

	var (
		value     string
		byteValue []byte
		values    []string
		//isExist   bool
		//err       error
		retCode   = "000000"
	)

	//switch req.GetMethod() {
	//case "GET_STRING":
	//	value, err = rt.redisCli.GetString(req.GetPostData())
	//case "GET_BYTE":
	//	byteValue, err = rt.redisCli.Get(req.GetPostData())
	//case "GET_STRINGS":
	//	values, err = rt.redisCli.Keys(req.GetPostData())
	//case "EXIST":
	//	isExist, err = rt.redisCli.HExistString(req.GetPostData())
	//}
	//
	//if err != nil {
	//	fmt.Println("redis error: ", req.GetPostData())
	//	return &DataResponse{StatusCode: 400, ReturnCode: "999999"}
	//}
	//
	//if !isExist {
	//	retCode = "000001"
	//}

	return &DataResponse{
		Body:       byteValue,
		BodyStr:    value,
		BodyStrs:   values,
		StatusCode: 200,
		ReturnCode: retCode,
	}

}

func (rt *RedisTransfer) connect() {
	redOnce.Do(func() {
		//options := &redis.Options{
		//	Addr:     "",
		//	DB:       6,
		//	PoolSize: 10,
		//	// ReadOnly: readOnly,
		//}

		redisLib.GetRedisClient()
	})
}

func (rt *RedisTransfer) Close() {

}
