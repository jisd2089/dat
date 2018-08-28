package service

/**
    Author: luzequan
    Created: 2018-02-27 10:48:39
*/

import (
	SSH "drcs/common/ssh"
	"drcs/common/sftp"
	"drcs/dep/agollo"
	"drcs/settings"
	"drcs/dep/handler/msg"
	logger "drcs/log"

	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/ssh"
	"sync"
	"net"
	"time"
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/valyala/fasthttp"
	"path/filepath"
	"drcs/dep/security"
	"strconv"
	"os"
	"drcs/common/balance"
	redisLib "drcs/common/redis"
	"errors"
	"math/rand"
	"drcs/dep/member"
	"drcs/runtime/scheduler"
)

var (
	once        sync.Once
	SettingPath string
	wg          sync.WaitGroup
	nsOnce      sync.Once
	redisCli    redisLib.RedisClient
)

type NodeService struct {
	nodeCh     chan bool
	sshClient  *SSH.SSHClient
	SftpClient *sftp.SFTPClient
	lock       sync.RWMutex
	redisCli   redisLib.RedisClient
}

func NewNodeService() *NodeService {
	return &NodeService{
	}
}

func (s *NodeService) Init() {

	wg.Add(2)
	NewDepService().Init()
	NewMemberService().Init()
	NewOrderService().Init()
	//NewRouteService().Init()

	s.init()

	wg.Wait()

	s.connectRedis()
	go InitCrpScheduler(s)

	//fmt.Println("init end")
}

func (s *NodeService) init() {
	s.nodeCh = make(chan bool, 1)

	path := filepath.Join(SettingPath, "setting.properties")

	go s.initApollo(filepath.Clean(path))

	// 初始化日志 和 security
	select {
	case ret := <-s.nodeCh:
		fmt.Println("logger init", ret)

		logger.Initialize()

		if err := security.Initialize(); err != nil {
			logger.Error("security initialize failed %s", err.Error())
		}

		wg.Done()
		break
	}

	//defaultInitSecurityConfig()

	//// 初始化xml订单文件信息
	//or.InitOrderRouteFile()
	//
	//memberManager, err := member.GetMemberManager()
	//if err != nil {
	//	logger.Debug("common init get member manager failed ", err)
	//}
	//logger.Info("init with memberManager:%+v", memberManager)
}

func (s *NodeService) initApollo(configDir string) {

	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		fmt.Println("initApollo")

		changesCnt := changeEvent.Changes["content"]
		value := changesCnt.NewValue

		switch changesCnt.ChangeType {
		case 0:
			common := &settings.CommonSettings{}
			err := yaml.Unmarshal([]byte(value), common)
			if err != nil {
			}

			settings.SetCommonSettings(common)
		case 1:
			common := &settings.CommonSettings{}
			err := yaml.Unmarshal([]byte(value), common)
			if err != nil {
			}
			settings.SetCommonSettings(common)
		}

		s.nodeCh <- true
	}
}

// 默认初始化公私钥
func defaultInitSecurityConfig() {

	memberId := settings.GetCommonSettings().Node.MemberId
	userkey := settings.GetCommonSettings().Node.Userkey
	token := settings.GetCommonSettings().Node.Token
	services_type := settings.GetCommonSettings().Node.Role
	url := settings.GetCommonSettings().Node.DlsUrl

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
	err := fasthttp.Do(request, response)
	if err != nil {
		logger.Error("post dls init node err ", err)
		return
	}
	data := response.Body()

	err = proto.Unmarshal(data, res_init_msg)
	if err != nil {
		logger.Error("failed to unmarshal data to res_init_msg", err)
		return
	}
	//status := res_init_msg.Status
	//if *status == "0" {
	//	return
	//}
}

func (s *NodeService) connect(fileCataLog *sftp.FileCatalog) {
	once.Do(func() {
		s.connectSSH(fileCataLog)
		s.connectSFTP()
	})
}

func (ns *NodeService) connectSSH(fileCataLog *sftp.FileCatalog) error {
	userName := fileCataLog.UserName
	password := fileCataLog.Password
	host := fileCataLog.Host
	port := fileCataLog.Port
	timeout := fileCataLog.TimeOut

	sshAuth := make([]ssh.AuthMethod, 0)
	sshAuth = append(sshAuth, ssh.Password(password))

	if timeout == 0 {
		timeout = 30 * time.Second
	}

	clientConfig := &ssh.ClientConfig{
		User:    userName,
		Auth:    sshAuth,
		Timeout: timeout,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	address := fmt.Sprintf("%s:%d", host, port)

	if err := ns.sshClient.Init(address, clientConfig); err != nil {
		return err
	}
	return nil
}

func (ns *NodeService) connectSFTP() error {

	if err := ns.SftpClient.Init(ns.sshClient.Client); err != nil {
		return err
	}
	return nil
}

type RecordProcess struct {
	ns *NodeService
}

func (p RecordProcess) Run() {
	p.ns.RecordProcess()
}

type CheckStatus struct {
	ns *NodeService
}

func (c CheckStatus)Run() {
	c.ns.CheckAndUpdateStatus()
}

func InitCrpScheduler(ns *NodeService) {
	scheduler.Init()
	scheduler.Schedule(time.Second * time.Duration(10), RecordProcess{ns})
	scheduler.Schedule(time.Second * time.Duration(60), CheckStatus{ns})
}

func (ns *NodeService) connectRedis() {
	nsOnce.Do(func() {
		common := settings.GetCommonSettings()
		o := &redisLib.ConnectOptions{
			AddressList: common.Redis.Addr,
			Password:    "",
			DBIndex:     common.Redis.DB,
			PoolSize:    common.Redis.PoolSize,
		}
		redisCli, _ = redisLib.GetRedisClient(o)
		//if err != nil {
		//	return
		//}
	})
	ns.redisCli = redisCli
}

func (ns *NodeService) CheckAndUpdateStatus() (bool, error) {

	pids, err := ns.redisCli.HGetAll("process")
	if err != nil {
		logger.Error(fmt.Sprintf("[NodeService] get pids from Redis err :%v", err))
		return false, errors.New(fmt.Sprintf("get pids from Redis err :%v", err))
	}

	pid_dead := checkStatus(pids)
	if pid_dead == "" {
		logger.Error("[NodeService] pids no dead", pids)
		return true, nil
	}

	duration := time.Duration(rand.Intn(30))
	time.Sleep(duration * time.Second)

	pid_deadx := checkStatus(pids)

	if pid_dead != pid_deadx {
		return true, nil
	}

	err = ns.reAllocateQuota(pid_dead)
	if err != nil {
		logger.Error(fmt.Sprintf("[NodeService] error while check process status"))
		return false, err
	}
	return true, nil
}

func checkStatus(pids map[string]string) string {
	now := time.Now().UnixNano()
	interval := int64(30e9) //30 seconds
	for pid, v := range pids {
		stmp, _ := strconv.ParseInt(v, 10, 64)
		if now-stmp > interval {
			return pid
		}
	}
	return ""
}

//lucky process stratery: the process who does reallocation job gets the corresponding quota
func (ns *NodeService) reAllocateQuota(pid_dead string) error {

	mem_deads, err := ns.redisCli.HGetAll(pid_dead)
	if err != nil {
		logger.Error(fmt.Sprintf("Get mem_deads from Redis err :%v", err))
		return errors.New(fmt.Sprintf("Get mem_deads from Redis err :%v", err))
	}

	pid := strconv.Itoa(os.Getpid())

	memId2del := []string{}
	for memId, v := range mem_deads {
		memId2del = append(memId2del, memId)

		//yagrusLog.Error("quota job###################: ", v)
		if quota, err := strconv.ParseInt(v, 10, 64); err == nil {
			ns.redisCli.HIncrBy(pid, memId, quota)
			balance.UpdateBalance(memId, float64(quota)/1000)
		}
	}

	ns.redisCli.HDelString(pid_dead, memId2del)
	ns.redisCli.HDelString("process", []string{pid_dead})

	return nil
}

func (ns *NodeService) RecordProcess() error {

	key := "process"
	pid := strconv.Itoa(os.Getpid())

	if err := ns.redisCli.HSetIn64(key, pid, time.Now().UnixNano()); err != nil {
		return errors.New(fmt.Sprintf("Get Value from Redis err :%v", err))
	}

	ns.updateQuota(pid) // 将内存中的额度同步到redis

	return nil
}

func (ns *NodeService) updateQuota(pid string) {
	memberId := member.GetMemberInfoList().MemberDetailList.MemberDetailInfo[0].MemberId
	logger.Info("[NodeService] updateQuota sync balance >>>>>>>>>>>>>>>", memberId, *balance.BanlanceMutex[memberId])
	balance := balance.BanlanceMutex[memberId]
	ns.redisCli.HSetIn64(pid, memberId, *balance)
}
