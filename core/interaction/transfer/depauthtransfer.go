package transfer

import (
	"drcs/core/interaction/response"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

/**
    Author: luzequan
    Created: 2018-08-01 14:15:01
*/
type DepAuthTransfer struct{}

func NewDepAuthTransfer() Transfer {
	return &DepAuthTransfer{}
}

// 封装dep auth服务
func (ft *DepAuthTransfer) ExecuteMethod(req Request) Response {

	// dep 认证成功返回“000000”， 认证失败返回“000004”
	memberId := req.Param("memberId")
	serialNo := req.Param("serialNo")
	reqSign := req.Param("reqSign")
	pubkey := req.Param("pubkey")
	jobId := req.Param("jobId")

	retCode := "000000"
	retMsg := "authentication success"

	switch req.GetMethod() {
	case "APPKEY":
		if !appKeyAuthentication(memberId, serialNo, reqSign, pubkey, jobId) {
			retCode = "000004"
			retMsg = "authentication failed"
		}
	case "Prikey":

	}

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: retCode,
		//Bobject:    req.GetBobject(),
		ReturnMsg:  retMsg,
	}
}

func appKeyAuthentication(memId, serialNo, reqSign, pubkey, jobId string) bool {

	authenticationInfo := memId + serialNo + jobId + pubkey
	hash := sha256.New()
	hash.Write([]byte(authenticationInfo))
	md := hash.Sum(nil)
	authenticationHash := hex.EncodeToString(md)

	return strings.EqualFold(string(authenticationHash), reqSign)
}

func (ft *DepAuthTransfer) Close() {}
