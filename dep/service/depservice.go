package service

/**
    Author: luzequan
    Created: 2018-05-10 14:59:03
*/
import (
	"sync"
	"gopkg.in/yaml.v2"
	logger "drcs/log"

	"drcs/dep/agollo"
	"drcs/core"
	"drcs/core/interaction/request"
	. "drcs/dep/nodelib/crp/common"
	st "drcs/settings"
	"path/filepath"
	"drcs/dep/nodelib"
	"path"
	"bytes"
	"github.com/valyala/fasthttp"
	"drcs/core/databox"
	"fmt"
	"encoding/json"
	"drcs/dep/order"
	"time"
	"drcs/dep/member"
)

type DepService struct {
	DataPath  string
	JobId     string
	PartnerId string
	lock      sync.RWMutex
}

type BatchParams struct {
	SeqNo     []byte
	TaskId    []byte
	UserId    []byte
	JobId     []byte
	IdType    []byte
	DataRange []byte
	MaxDelay  []byte
	MD5       []byte
}

func NewDepService() *DepService {
	return &DepService{}
}

func (s *DepService) Init() {
	path := filepath.Join(SettingPath, "trans.properties")
	go initTransConfig(filepath.Clean(path))
}

func (s *DepService) Process() {

	transInfo := GetTransInfo()

	common := st.GetCommonSettings()

	logger.Info("transInfo", transInfo)
	logger.Info("common setting", common)

	nodeMemberId := transInfo.Trans.MemberId

	fsAddress := &request.FileServerAddress{
		Host:      common.Sftp.Hosts,
		Port:      common.Sftp.Port,
		UserName:  common.Sftp.Username,
		Password:  common.Sftp.Password,
		TimeOut:   common.Sftp.DefualtTimeout,
		LocalDir:  common.Sftp.LocalDir,
		RemoteDir: common.Sftp.RemoteDir,
	}

	dataAddrs := []*Dest{}
	algAddrs := []*Dest{}

	for _, val := range transInfo.Trans.Dest {
		switch val.Type {
		case "data":
			dataAddrs = append(dataAddrs, val)
		case "algorithm":
			algAddrs = append(algAddrs, val)
		}
	}

	logger.Info("dataAddrs", dataAddrs)
	logger.Info("algAddrs", algAddrs)

	runDataBox(dataAddrs, "datasend", nodeMemberId, fsAddress)

	runDataBox(algAddrs, "algorithmsend", nodeMemberId, fsAddress)

}

func initTransConfig(configDir string) {
	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		fmt.Println("initTransConfig")

		changesCnt := changeEvent.Changes["content"]
		value := changesCnt.NewValue

		transInfo := &TransmissionInfo{}
		err := yaml.Unmarshal([]byte(value), transInfo)
		if err != nil {
		}

		SetTransInfo(transInfo)
	}
}

func runDataBox(addrs []*Dest, boxName string, nodeMemberId string, fsAddress *request.FileServerAddress) {
	for _, v := range addrs {

		b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
		if b == nil {
			logger.Error("databox is nil!")
			return
		}
		b.SetDataFilePath(v.DataPath)

		addrs := []*request.NodeAddress{}
		addrs = append(addrs, &request.NodeAddress{
			MemberId: nodeMemberId,
			Host:     v.DestHost,
			Port:     v.DestPort,
			URL:      v.Api,
			Priority: 0,})

		//b.SetNodeAddress(addrs)
		b.FileServerAddress = fsAddress
		// 算法依赖文件hdfs路径
		b.Params = v.RelyDatas

		setDataBoxQueue(b)
	}
}

// 处理批量配送——发送
func (s *DepService) ProcessBatchDis(ctx *fasthttp.RequestCtx) {

	reqFilePath := string(ctx.FormValue("reqFilePath"))
	if reqFilePath == "" {
		logger.Error("reqFilePath missing")
		return
	}
	boxName := string(ctx.FormValue("boxName"))
	if boxName == "" {
		logger.Error("box name missing")
		return
	}

	reqFileName := path.Base(reqFilePath)

	dataFileName := &nodelib.DataFileName{}
	if err := dataFileName.ParseAndValidFileName(reqFileName); err != nil {
		logger.Error("Parse and valid fileName: [%s] error: %s", reqFileName, err)
		return
	}

	//prefixName := dataFileName.GetPrefixName()
	jobId := dataFileName.JobId
	idType := dataFileName.IdType
	batchNo := dataFileName.BatchNo
	fileNo := dataFileName.FileNo

	common := st.GetCommonSettings()
	logger.Info("common setting", common)

	fsAddress := &request.FileServerAddress{
		Host:      common.Sftp.Hosts,
		Port:      common.Sftp.Port,
		UserName:  common.Sftp.Username,
		Password:  common.Sftp.Password,
		TimeOut:   common.Sftp.DefualtTimeout,
		LocalDir:  common.Sftp.LocalDir,
		RemoteDir: common.Sftp.RemoteDir,
	}

	//boxName = "batch_sup_send"
	//boxName := "batch_dem_send"
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
	if b == nil {
		logger.Error("databox is nil!")
		return
	}
	b.SetDataFilePath(reqFilePath)
	b.SetParam("jobId", jobId)
	b.SetParam("idType", idType)
	b.SetParam("batchNo", batchNo)
	b.SetParam("fileNo", fileNo)
	b.SetParam("NodeMemberId", common.Node.MemberId)

	b.Params = common.Redis.Addr

	b.FileServerAddress = fsAddress

	setDataBoxQueue(b)
}

// 处理批量配送——接收
func (s *DepService) ProcessBatchRcv(ctx *fasthttp.RequestCtx, targetFilePath string) {

	boxName := string(ctx.Request.Header.Peek("boxName"))
	if boxName == "" {
		logger.Error("box name missing")
		return
	}

	respFileName := path.Base(targetFilePath)

	dataFileName := &nodelib.DataFileName{}
	if err := dataFileName.ParseAndValidFileName(respFileName); err != nil {
		logger.Error("Parse and valid fileName: [%s] error: %s", respFileName, err)
		return
	}

	//boxName = "batch_dem_rcv"
	//boxName := "batch_sup_rcv"
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
	if b == nil {
		logger.Error("databox is nil!")
		return
	}

	if err := setRcvParams(ctx, b); err != nil {
		logger.Error("rcv params err [%s]", err.Error())
		return
	}

	common := st.GetCommonSettings()
	logger.Info("common setting", common)

	fsAddress := &request.FileServerAddress{
		Host:      common.Sftp.Hosts,
		Port:      common.Sftp.Port,
		UserName:  common.Sftp.Username,
		Password:  common.Sftp.Password,
		TimeOut:   common.Sftp.DefualtTimeout,
		LocalDir:  common.Sftp.LocalDir,
		RemoteDir: common.Sftp.RemoteDir,
	}

	b.SetDataFilePath(targetFilePath)
	b.FileServerAddress = fsAddress
	b.SetParam("jobId", dataFileName.JobId)
	b.SetParam("batchNo", dataFileName.BatchNo)
	b.SetParam("fileNo", dataFileName.FileNo)
	b.SetParam("NodeMemberId", common.Node.MemberId)

	b.Params = common.Redis.Addr

	setDataBoxQueue(b)

}

// 金融消费
func (s *DepService) ProcessCrpTrans(ctx *fasthttp.RequestCtx) {
	//logger.Info("DepService ProcessCrpTrans start")

	bodyChan := make(chan []byte)

	//timeOut := time.Duration(300000) * time.Millisecond

	//boxName, err := getCrpBoxName(ctx.Request.Body())
	//if err != nil {
	//	return
	//}
	common := st.GetCommonSettings()

	boxName := "dem_request"
	//boxName = "smart_request"
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
	if b == nil {
		logger.Error("databox [%s] is nil!", boxName)
		return
	}

	pubkey := member.GetMemberInfoList().MemberDetailList.MemberDetailInfo[0].PubKey
	b.SetParam("pubkey", pubkey)
	b.SetParam("balanceUrl", common.Other.Crp.BalanceUrl)
	b.SetParam("prdtIdCd", "1003004")

	// redis address
	b.Params = common.Redis.Addr
	b.BodyChan = bodyChan

	b.HttpRequestBody = ctx.Request.Body()

	setDataBoxQueue(b)

	select {
	case body := <-bodyChan:
		ctx.SetBody(body)
		close(bodyChan)
	//case <-time.After(timeOut):
	//	logger.Error("http response timeout")
	//	break
	}

	//fmt.Println("DepService middle time: ", time.Since(middle))
	//fmt.Println("DepService run time: ", time.Since(start))
}

func (s *DepService) ProcessCrpResponse(ctx *fasthttp.RequestCtx) {
	logger.Info("DepService ProcessCrpResponse start")

	timeOut := time.Duration(300000) * time.Millisecond

	prdtIdCd := string(ctx.Request.Header.Peek("prdtIdCd"))
	if prdtIdCd == "" {
		logger.Error("prdtIdC d is nil!")
		return
	}

	serialNo := string(ctx.Request.Header.Peek("serialNo"))
	if serialNo == "" {
		logger.Error("serialNo is nil!")
		return
	}

	busiSerialNo := string(ctx.Request.Header.Peek("busiSerialNo"))
	if busiSerialNo == "" {
		logger.Error("busiSerialNo is nil!")
		return
	}

	boxName := "sup_response"
	//boxName = "smart_response"
	b := assetnode.AssetNodeEntity.GetDataBoxByName(boxName)
	if b == nil {
		logger.Error("databox is nil!")
		return
	}

	b.SetParam("serialNo", serialNo)
	b.SetParam("busiSerialNo", busiSerialNo)

	b.HttpRequestBody = ctx.Request.Body()

	bodyChan := make(chan []byte)
	b.BodyChan = bodyChan

	setDataBoxQueue(b)

	select {
	case body := <-bodyChan:
		ctx.SetBody(body)
		close(bodyChan)
	case <-time.After(timeOut):
		logger.Error("http response timeout")
		break
	}
}

func setRcvParams(ctx *fasthttp.RequestCtx, b *databox.DataBox) error {

	batchParams := &BatchParams{
		SeqNo:     ctx.Request.Header.Peek("seqNo"),
		JobId:     ctx.Request.Header.Peek("orderId"),
		TaskId:    ctx.Request.Header.Peek("taskId"),
		UserId:    ctx.Request.Header.Peek("userId"),
		IdType:    ctx.Request.Header.Peek("idType"),
		DataRange: ctx.Request.Header.Peek("dataRange"),
		MaxDelay:  ctx.Request.Header.Peek("maxDelay"),
		MD5:       ctx.Request.Header.Peek("md5"),
	}

	if err := checkRcvParams(batchParams); err != nil {
		return err
	}

	b.SetParam("seqNo", string(batchParams.SeqNo))
	b.SetParam("taskId", string(batchParams.TaskId))
	b.SetParam("orderId", string(batchParams.JobId))
	b.SetParam("userId", string(batchParams.UserId))
	b.SetParam("idType", string(batchParams.IdType))
	b.SetParam("dataRange", string(batchParams.DataRange))
	b.SetParam("maxDelay", string(batchParams.MaxDelay))
	b.SetParam("md5", string(batchParams.MD5))

	return nil
}

func checkRcvParams(p *BatchParams) error {

	var errMsg bytes.Buffer
	if len(p.SeqNo) == 0 {
		errMsg.WriteString("[seqNo]")
	}
	if len(p.TaskId) == 0 {
		errMsg.WriteString("[taskId]")
	}
	if len(p.JobId) == 0 {
		errMsg.WriteString("[jobId]")
	}
	if len(p.IdType) == 0 {
		errMsg.WriteString("[idType]")
	}
	if len(p.UserId) == 0 {
		errMsg.WriteString("[userId]")
	}
	if len(p.DataRange) == 0 {
		errMsg.WriteString("[dataRange]")
	}
	if len(p.MaxDelay) == 0 {
		errMsg.WriteString("[maxDelay]")
	}
	if len(p.MD5) == 0 {
		errMsg.WriteString("[md5]")
	}
	if len(errMsg.String()) != 0 {
		return fmt.Errorf(" params %s missing", errMsg.String())
	}

	return nil
}

func getCrpBoxName(requestBody []byte) (string, error) {

	commonRequestData := &CommonRequestData{}
	if err := json.Unmarshal(requestBody, &commonRequestData); err != nil {
		logger.Error("marshal request body err ", )
		return "", fmt.Errorf("marshal request body err ")
	}

	jobId := commonRequestData.PubReqInfo.JobId

	common := st.GetCommonSettings()
	memberId := common.Node.MemberId

	var prdtIdCd string

	orderInfoMap := order.GetOrderInfoMap()

	orderData, ok := orderInfoMap[jobId]
	if !ok {
		logger.Error("marshal request body err ", )
		return "", fmt.Errorf("marshal request body err ")
	}
	for _, v := range orderData.TaskInfoMapById {
		prdtIdCd = v.PrdtIdCd
		break
	}

	return fmt.Sprint("%s_%s", memberId, prdtIdCd), nil
}
