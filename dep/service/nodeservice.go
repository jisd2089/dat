package service

/**
    Author: luzequan
    Created: 2018-02-27 10:48:39
*/

import (
	"drcs/dep/security"
	"drcs/dep/or"
	"drcs/dep/member"
	logger "dds/log"
	SSH "drcs/common/ssh"
	"drcs/common/sftp"

	"golang.org/x/crypto/ssh"
	"sync"
	"net"
	"time"
	"fmt"
)

var (
	once      sync.Once
)

type NodeService struct {
	sshClient  *SSH.SSHClient
	SftpClient *sftp.SFTPClient
	lock       sync.RWMutex
}

func NewNodeService() *NodeService {
	return &NodeService{
		sshClient:  SSH.New(),
		SftpClient: sftp.New(),
	}
}

func (s *NodeService) Init() {

	//
	if err := security.Initialize(); err != nil {
		return // TODO
	}

	// 初始化xml订单文件信息
	or.InitOrderRouteFile()

	memberManager, err := member.GetMemberManager()
	if err != nil {
		logger.Debug("common init get member manager failed ", err)
	}
	logger.Info("init with memberManager:%+v", memberManager)
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