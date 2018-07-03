package handler

import (
	"github.com/valyala/fasthttp"
	"github.com/golang/protobuf/proto"


	"crypto/sha256"
	"strings"
	"encoding/hex"
	"io"

	. "drcs/settings"
	"drcs/dep/handler/msg"
	logger "drcs/log"
	"drcs/dep/security"
	"drcs/core"
	"os"
	"path"
	"drcs/core/interaction/request"
	"drcs/dep/service"
	"bytes"
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

func (n *NodeHandler) RcvData(ctx *fasthttp.RequestCtx) {

	dataFile, err := ctx.FormFile("file")
	if err != nil {
		logger.Error("filePath err:", err)
		return
	}
	logger.Info("filePath***********: ", dataFile.Filename)

	//bn := ctx.FormValue("boxname")
	//boxName := string(bn)
	boxName := "datareceive" //TODO

	common := GetCommonSettings()
	hdfsInputDir := common.Hdfs.InputDir
	hdfsOutputDir := common.Hdfs.OutputDir
	targetFileDir := common.Sftp.LocalDir

	targetFilePath := path.Join(targetFileDir, dataFile.Filename)

	targetFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer targetFile.Close()
	if err != nil {
		logger.Error("open target file err:", err)
		return
	}

	dataFileContent, err := dataFile.Open()
	defer dataFileContent.Close()
	if err != nil {
		logger.Error("open form file err:", err)
		return
	}

	io.Copy(targetFile, dataFileContent)

	fsAddress := &request.FileServerAddress{
		Host:      common.Sftp.Hosts,
		Port:      common.Sftp.Port,
		UserName:  common.Sftp.Username,
		Password:  common.Sftp.Password,
		TimeOut:   common.Sftp.DefualtTimeout,
		LocalDir:  common.Sftp.LocalDir,
		RemoteDir: common.Sftp.RemoteDir,
	}

	// 1.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)

	if b == nil {
		logger.Error("databox is nil!")
		return
	}

	b.DataFilePath = targetFilePath
	b.FileServerAddress = fsAddress
	b.SetParam("dataPath", "/logistic/input/logistic3.txt") // TODO
	b.SetParam("hdfsInputDir", hdfsInputDir) // TODO
	b.SetParam("hdfsOutputDir", hdfsOutputDir) // TODO

	// 1.2 setDataBoxQueue
	setDataBoxQueue(b)
}

func (n *NodeHandler) RcvAlg(ctx *fasthttp.RequestCtx) {

	// 算法文件
	dataFile, err := ctx.FormFile("file")
	if err != nil {
		logger.Error("filePath err:", err)
		return
	}
	logger.Info("filePath***********: ", dataFile.Filename)

	// 算法依赖文件hdfs路径
	var dataFiles []string

	mf, err := ctx.MultipartForm()
	if err == nil && mf.Value != nil {
		vv := mf.Value["dataFiles"]
		if len(vv) > 0 {
			dataFiles = vv
		}
	}

	// boxname
	//bn := ctx.FormValue("boxname")
	//boxName := string(bn)
	boxName := "algorithmreceive" //TODO

	common := GetCommonSettings()
	hdfsInputDir := common.Hdfs.InputDir
	hdfsOutputDir := common.Hdfs.OutputDir // 算法输出结果路径
	targetFileDir := common.Sftp.LocalDir

	targetFilePath := path.Join(targetFileDir, dataFile.Filename)

	targetFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer targetFile.Close()
	if err != nil {
		logger.Error("open target file err:", err)
		return
	}

	dataFileContent, err := dataFile.Open()
	defer dataFileContent.Close()
	if err != nil {
		logger.Error("open form file err:", err)
		return
	}

	io.Copy(targetFile, dataFileContent)

	fsAddress := &request.FileServerAddress{
		Host:      common.Sftp.Hosts,
		Port:      common.Sftp.Port,
		UserName:  common.Sftp.Username,
		Password:  common.Sftp.Password,
		TimeOut:   common.Sftp.DefualtTimeout,
		LocalDir:  common.Sftp.LocalDir,
		RemoteDir: common.Sftp.RemoteDir,
	}

	// 1.1 匹配相应的DataBox
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)

	if b == nil {
		logger.Error("databox is nil!")
		return
	}

	b.DataFilePath = targetFilePath
	b.FileServerAddress = fsAddress
	b.Params = dataFiles
	b.SetParam("dataPath", "") // TODO
	b.SetParam("hdfsInputDir", hdfsInputDir) // TODO
	b.SetParam("hdfsOutputDir", hdfsOutputDir) // TODO

	// 1.2 setDataBoxQueue
	setDataBoxQueue(b)
}

func (n *NodeHandler) RunProcess(ctx *fasthttp.RequestCtx) {
	service.NewDepService().Process()
}

func (n *NodeHandler) RunBatchProcess(ctx *fasthttp.RequestCtx) {
	service.NewDepService().ProcessBatchDis(ctx)
}

func (n *NodeHandler) RunBatchRcv(ctx *fasthttp.RequestCtx) {

	if len(ctx.Request.Header.Peek("dataFile")) == 0 {
		logger.Error("data file name is null")
		return
	}
	dataFile := string(ctx.Request.Header.Peek("dataFile"))
	boxName := string(ctx.Request.Header.Peek("boxName"))
	logger.Info("filePath***********: ", dataFile)
	seqNo := string(ctx.Request.Header.Peek("seqNo"))
	logger.Info("filePath***********: ", seqNo)
	logger.Info("boxName***********: ", boxName)

	common := GetCommonSettings()
	//hdfsInputDir := common.Hdfs.InputDir
	//hdfsOutputDir := common.Hdfs.OutputDir
	targetFileDir := common.Sftp.LocalDir
	//targetFileDir = "D:/dds_receive/tmp"

	targetFilePath := path.Join(targetFileDir, dataFile)

	targetFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer targetFile.Close()
	if err != nil {
		logger.Error("open target file err:", err)
		return
	}

	io.Copy(targetFile, bytes.NewReader(ctx.Request.Body()))

	service.NewDepService().ProcessBatchRcv(ctx, targetFilePath)
}