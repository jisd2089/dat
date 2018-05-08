package handler

import (
	"github.com/valyala/fasthttp"
	"github.com/golang/protobuf/proto"


	"crypto/sha256"
	"strings"
	"encoding/hex"

	. "drcs/settings"
	"drcs/dep/handler/msg"
	logger "dds/log"
	"drcs/dep/security"
)

/**
    Author: luzequan
    Created: 2018-05-08 15:19:02
*/
type NodeHandler struct{}

func NewNodeHandler() *NodeHandler {
	return &NodeHandler{}
}

func (n *NodeHandler) InitSecurityConfig(ctx *fasthttp.RequestCtx) {

	memberId := GetCommonSettings().Node.MemberId
	userkey := GetCommonSettings().Node.Userkey
	token := GetCommonSettings().Node.Token
	services_type := GetCommonSettings().Node.Role
	url := GetCommonSettings().Node.DlsUrl

	if isLocal := string(ctx.FormValue("isLocal")); strings.ToLower(isLocal) == "true" {
		hash := sha256.New()
		hash.Write([]byte(memberId + token))
		sum := hash.Sum(nil)
		seed := hex.EncodeToString(sum)
		security.SaveDataToxml(seed, memberId, userkey)
		ctx.Response.SetBody([]byte("success"))
		return
	}

	req_init_msg := &msg_dem.PBDDlsReqMsg{}
	res_init_msg := &msg_dem.PBDDlsResMsg{}
	req_init_msg.MemId = &memberId
	req_init_msg.UserPswd = &userkey
	req_init_msg.Token = &token
	req_init_msg.Role = &services_type
	body, _ := proto.Marshal(req_init_msg)

	request := &fasthttp.Request{}
	request.SetRequestURI(url)
	request.Header.SetMethod("POST")
	request.SetBody(body)
	response := &fasthttp.Response{}
	err0 := fasthttp.Do(request, response)
	if err0 != nil {
		logger.Error("post dls init node err ", err0)
		ctx.Response.SetBody([]byte("post dls init node failed"))
		return
	}
	data := response.Body()

	err := proto.Unmarshal(data, res_init_msg)
	if err != nil {
		logger.Error("failed to unmarshal data to res_init_msg", err)
		ctx.Response.SetBody([]byte("failed to unmarshal data to res_init_msg"))
		return
	}
	status := res_init_msg.Status
	if *status == "0" {
		ctx.Response.SetBody([]byte("success"))
		return
	}
	ctx.Response.SetBody([]byte(*res_init_msg.ErrMsg))
}

func (n *NodeHandler) GenKeys(ctx *fasthttp.RequestCtx) {
	data := ctx.Request.Body()
	req_keyGen_msg := &msg_dem.DReqKeyGenMsg{}
	proto.Unmarshal(data, req_keyGen_msg)
	seed := *req_keyGen_msg.KeySeed
	memId := *req_keyGen_msg.MemId
	userkey := *req_keyGen_msg.UserPswd
	pubkey, err := security.SaveDataToxml(seed, memId, userkey)

	res_keyGen_msg := &msg_dem.DResKeyGenMsg{}
	var status, errNo, errMsg, pubKey string
	if err != nil {
		status = "-1"
		errNo = "000001"
		errMsg = "failed to generate public key"
		pubKey = ""
		logger.Error("failed to generate public key")
	} else {
		status = "0"
		errNo = ""
		errMsg = ""
		pubKey = pubkey
	}

	go security.Initialize()

	res_keyGen_msg.Status = &status
	res_keyGen_msg.ErrNO = &errNo
	res_keyGen_msg.ErrMsg = &errMsg
	res_keyGen_msg.PubKey = &pubKey
	logger.Info(" generate public key resp: ", res_keyGen_msg)
	body, _ := proto.Marshal(res_keyGen_msg)
	ctx.Response.SetBody(body)
}