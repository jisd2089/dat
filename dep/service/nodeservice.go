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
)

var (
	once        sync.Once
	SettingPath string
)

type NodeService struct {
	nodeCh     chan bool
	sshClient  *SSH.SSHClient
	SftpClient *sftp.SFTPClient
	lock       sync.RWMutex
}

func NewNodeService() *NodeService {
	return &NodeService{}
}

func (s *NodeService) Init() {

	NewDepService().Init()
	NewMemberService().Init()
	NewOrderService().Init()
	NewRouteService().Init()

	s.init()
	fmt.Println("init end")
}

func (s *NodeService) init() {
	s.nodeCh = make(chan bool, 1)

	path := filepath.Join(SettingPath, "setting.properties")

	go s.initApollo(filepath.Clean(path))

	// 初始化日志
	select {
	case ret := <-s.nodeCh:
		fmt.Println("logger init", ret)
		logger.Initialize()
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

func (ft *NodeService) connectSSH(fileCataLog *sftp.FileCatalog) error {
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

	if err := ft.sshClient.Init(address, clientConfig); err != nil {
		return err
	}
	return nil
}

func (ft *NodeService) connectSFTP() error {

	if err := ft.SftpClient.Init(ft.sshClient.Client); err != nil {
		return err
	}
	return nil
}
