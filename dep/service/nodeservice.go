package service

/**
    Author: luzequan
    Created: 2018-02-27 10:48:39
*/

import (
	//"drcs/dep/security"
	//"drcs/dep/or"
	//"drcs/dep/member"
	//logger "dds/log"
	SSH "drcs/common/ssh"
	"drcs/common/sftp"
	"drcs/dep/agollo"
	"drcs/settings"
	"drcs/dep/handler/msg"
	logger "dds/log"

	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/ssh"
	"sync"
	"net"
	"time"
	"fmt"
	//"encoding/json"
	yaml "gopkg.in/yaml.v2"
	"github.com/valyala/fasthttp"
)

var (
	once      sync.Once
)

func init() {

	//NewNodeService().Init()

}

type NodeService struct {
	sshClient  *SSH.SSHClient
	SftpClient *sftp.SFTPClient
	lock       sync.RWMutex
}

func NewNodeService() *NodeService {
	return &NodeService{
		//sshClient:  SSH.New(),
		//SftpClient: sftp.New(),
	}
}

func (s *NodeService) Init() {

	initApollo("D:/GoglandProjects/src/drcs/settings/setting.properties")

	defaultInitSecurityConfig()

	//
	//// 初始化xml订单文件信息
	//or.InitOrderRouteFile()
	//
	//memberManager, err := member.GetMemberManager()
	//if err != nil {
	//	logger.Debug("common init get member manager failed ", err)
	//}
	//logger.Info("init with memberManager:%+v", memberManager)
}

func initApollo(configDir string) {
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
		//bytes, _ := json.Marshal(changeEvent)
		//fmt.Println("event:", string(bytes))
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