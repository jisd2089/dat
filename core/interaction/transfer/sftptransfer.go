package transfer

/**
    Author: luzequan
    Created: 2017-12-29 15:03:01
*/
import (
	"net"
	"fmt"
	SSH "dat/common/ssh"
	"dat/common/sftp"
	"golang.org/x/crypto/ssh"
	"time"
	"dat/core/interaction/response"
	"sync"
)

type SftpTransfer struct {
	sshClient  *SSH.SSHClient
	sftpClient *sftp.SFTPClient
}

func NewSftpTransfer() *SftpTransfer {
	return &SftpTransfer{
		sshClient:  SSH.New(),
		sftpClient: sftp.New(),
	}
}

var (
	sshCli  *SSH.SSHClient
	sftpCli *sftp.SFTPClient
	once    sync.Once
)

// 封装sftp服务
func (st *SftpTransfer) ExecuteMethod(req Request) Response {

	var err error
	st.connect(req.GetFileCatalog())

	switch req.GetMethod() {
	case "GET":
		fmt.Println("sftp get ^^^^^^^^^^^^^^^^^^^^: ", req.GetFileCatalog().RemoteFileName)
		err = st.sftpClient.RemoteGet(req.GetFileCatalog())
		if err != nil {
			fmt.Println("sftp get error: ", err)
		}
	case "PUT":
		fmt.Println("sftp put ^^^^^^^^^^^^^^^^^^^^: ", req.GetFileCatalog().LocalFileName)
		err = st.sftpClient.RemotePut(req.GetFileCatalog())

	}

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("sftp failed ^^^^^^^^^^^^^^^^^^^^: ", req.GetFileCatalog().RemoteFileName)
		return &response.DataResponse{
			StatusCode: 400,
			ReturnCode: "999999",
		}
	}

	fmt.Println("sftp success ^^^^^^^^^^^^^^^^^^^^: ", req.GetFileCatalog().RemoteFileName)

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: "000000",
	}

	return nil
}

func (st *SftpTransfer) connect(fileCataLog *sftp.FileCatalog) {
	once.Do(func() {

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
