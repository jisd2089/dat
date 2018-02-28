package transfer

/**
    Author: luzequan
    Created: 2017-12-29 15:03:01
*/
import (
	"net"
	"fmt"
	SSH "drcs/common/ssh"
	"drcs/common/sftp"
	"golang.org/x/crypto/ssh"
	"time"
	"drcs/core/interaction/response"
	"sync"
)

type SftpTransfer struct {
	sshClient  *SSH.SSHClient
	sftpClient *sftp.SFTPClient
	lock       sync.RWMutex
}

func NewSftpTransfer() *SftpTransfer {
	return &SftpTransfer{
		sshClient:  SSH.New(),
		sftpClient: sftp.New(),
	}
}

var (
	sshCli    *SSH.SSHClient
	sftpCli   *sftp.SFTPClient
	once      sync.Once
	retryOnce sync.Once
	lock      sync.RWMutex
)

// 封装sftp服务
func (st *SftpTransfer) ExecuteMethod(req Request) Response {

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("sftp recover error: ", err)
		}
	}()

	st.lock.Lock()
	defer st.lock.Unlock()

	var (
		err        error
		retryTimes = 0
	)

	st.connect(req.GetFileCatalog())

RETRY:
	switch req.GetMethod() {
	case "GET":
		fmt.Println("sftp get ^^^^^^^^^^^^^^^^^^^^: ", req.GetFileCatalog().RemoteFileName)
		err = st.sftpClient.RemoteGet(req.GetFileCatalog())
	case "PUT":
		fmt.Println("sftp put ^^^^^^^^^^^^^^^^^^^^: ", req.GetFileCatalog().LocalFileName)
		err = st.sftpClient.RemotePut(req.GetFileCatalog())
	case "CLOSE":
		st.Close()
	}

	if err != nil {
		fmt.Println(err.Error())

		if retryTimes >= 1 {
			fmt.Println("sftp failed ^^^^^^^^^^^^^^^^^^^^: ", req.GetFileCatalog().RemoteFileName)
			return &response.DataResponse{
				StatusCode: 400,
				ReturnCode: "999999",
			}
		}

		retryTimes ++

		st.refresh(req.GetFileCatalog())
		goto RETRY
	}

	fmt.Println("sftp success ^^^^^^^^^^^^^^^^^^^^: ", req.GetFileCatalog().RemoteFileName)
	defer func() {
		retryOnce = sync.Once{}
	}()

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: "000000",
	}
}

func (st *SftpTransfer) connect(fileCataLog *sftp.FileCatalog) {
	once.Do(func() {
		st.connectSSH(fileCataLog)
		st.connectSFTP()
	})
}

func (st *SftpTransfer) refresh(fileCataLog *sftp.FileCatalog) {
	retryOnce.Do(func() {
		st.connectSSH(fileCataLog)
		st.connectSFTP()
	})
}

func (st *SftpTransfer) Close() {
	st.sshClient.Client.Close()
	st.sftpClient.CloseSession()
}

func (ft *SftpTransfer) connectSSH(fileCataLog *sftp.FileCatalog) error {
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

func (ft *SftpTransfer) connectSFTP() error {

	if err := ft.sftpClient.Init(ft.sshClient.Client); err != nil {
		return err
	}
	return nil
}
