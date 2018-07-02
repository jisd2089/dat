package transfer

/**
    Author: luzequan
    Created: 2018-01-13 19:59:16
*/
import (
	. "drcs/core/interaction/response"
	redisLib "drcs/common/redis"
	//"gopkg.in/redis.v5"
	"fmt"
	"sync"
	"strings"
	"strconv"
)

type RedisTransfer struct {
	redisCli redisLib.RedisClient
}

func NewRedisTransfer() Transfer {
	return &RedisTransfer{}
}

var (
	redOnce sync.Once
	redRtOnce sync.Once
)

// 封装redis服务
func (rt *RedisTransfer) ExecuteMethod(req Request) Response {
	//fmt.Println("RedisTransfer")

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("redisTransfer recover error: ", err)
		}
	}()

	rt.connect(req)

	var (
		value     string
		byteValue []byte
		values    []string
		isExist   bool
		retCode = "000000"
		err        error
		retryTimes = 0
	)

RETRY:
	switch req.GetMethod() {
	case "PING":
		if err = rt.redisCli.Ping(); err != nil {
			if retryTimes >= 1 {
				fmt.Println("redis failed ^^^^^^^^^^^^^^^^^^^^: ")
				return &DataResponse{
					StatusCode: 200,
					ReturnCode: "999999",
				}
			}
			retryTimes ++

			rt.refresh(req)
			goto RETRY
		}

		defer func() {
			redRtOnce = sync.Once{}
		}()

	case "GET_STRING":
		value, err = rt.redisCli.GetString(req.GetPostData())
	case "HGET_STRING":
		value, err = rt.redisCli.HGetString(req.Param("key"), req.Param("field"))
	case "HSET_STRING":
		err = rt.redisCli.HSetString(req.Param("key"), req.Param("field"), req.Param("value"))
	case "GET_BYTE":
		byteValue, err = rt.redisCli.Get(req.GetPostData())
	case "GET_STRINGS":
		values, err = rt.redisCli.Keys(req.GetPostData())
	case "EXIST":
		isExist, err = rt.redisCli.HExistString(req.GetPostData())
		if !isExist {
			retCode = "000001"
		}
	}

	if err != nil {
		fmt.Println("redis error: ", err.Error())

		return &DataResponse{
			StatusCode: 200,
			ReturnCode: "000002",
			ReturnMsg: err.Error(),
		}
	}

	return &DataResponse{
		Body:       byteValue,
		BodyStr:    value,
		BodyStrs:   values,
		StatusCode: 200,
		ReturnCode: retCode,
	}
}

func (rt *RedisTransfer) connect(req Request) {
	redOnce.Do(func() {
		dbIndex, _ := strconv.Atoi(req.Param("redisDB"))
		//if err != nil {
		//	return
		//}
		redisPoolSize, _ := strconv.Atoi(req.Param("redisPoolSize"))
		//if err != nil {
		//	return
		//}
		//addrStr := req.Param("redisAddrs")
		//addrList := strings.Split(addrStr, ",")
		o := &redisLib.ConnectOptions{
			AddressList: req.GetCommandParams(),
			Password:    req.Param("redisPwd"),
			DBIndex:     dbIndex,
			PoolSize:    redisPoolSize,
		}
		rt.redisCli, _ = redisLib.GetRedisClient(o)
		//if err != nil {
		//	return
		//}
	})
}

func (rt *RedisTransfer) refresh(req Request) {
	redRtOnce.Do(func() {
		dbIndex, _ := strconv.Atoi(req.Param("redisDB"))
		//if err != nil {
		//	return
		//}
		redisPoolSize, _ := strconv.Atoi(req.Param("redisPoolSize"))
		//if err != nil {
		//	return
		//}
		//addrStr := req.Param("redisAddrs")
		//addrList := strings.Split(addrStr, ",")
		o := &redisLib.ConnectOptions{
			AddressList: req.GetCommandParams(),
			Password:    req.Param("redisPwd"),
			DBIndex:     dbIndex,
			PoolSize:    redisPoolSize,
		}
		rt.redisCli, _ = redisLib.GetRedisClient(o)
		//if err != nil {
		//	return
		//}
	})
}


func (rt *RedisTransfer) Close() {

}
