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

type SshTransfer struct {
	sshClient *SSH.SSHClient
	lock      sync.RWMutex
}

func NewSshTransfer() *SshTransfer {
	return &SshTransfer{
		sshClient: SSH.New(),
	}
}

var (
	sshOnce      sync.Once
	sshRetryOnce sync.Once
)

// 封装sftp服务
func (st *SshTransfer) ExecuteMethod(req Request) Response {

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("ssh transfer recover error: ", err)
		}
	}()

	st.lock.Lock()
	defer st.lock.Unlock()

	var (
		cmdline string
		err        error
		retryTimes = 0
	)

	st.connect(req.GetFileCatalog())

RETRY:
	switch req.GetMethod() {
	case "STRING":
		cmdline = req.GetCommandName()
	case "SLICE":
		cmdline = req.GetCommandLine()
	}

	client := st.sshClient.Client
	session, err := client.NewSession()

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

	defer func() {
		sshRetryOnce = sync.Once{}
	}()

	output, err := session.Output(cmdline)
	if err != nil {
		return &response.DataResponse{
			Body:       output,
			StatusCode: 500,
			ReturnCode: "999999",
		}
	}

	fmt.Println(string(output))
	fmt.Println("sftp success ^^^^^^^^^^^^^^^^^^^^: ", req.GetFileCatalog().RemoteFileName)

	return &response.DataResponse{
		Body:       output,
		StatusCode: 200,
		ReturnCode: "000000",
	}
}

func (st *SshTransfer) connect(fileCataLog *sftp.FileCatalog) {
	sshOnce.Do(func() {
		st.connectSSH(fileCataLog)
	})
}

func (st *SshTransfer) refresh(fileCataLog *sftp.FileCatalog) {
	sshRetryOnce.Do(func() {
		st.connectSSH(fileCataLog)
	})
}

func (st *SshTransfer) Close() {
	st.sshClient.Client.Close()
}

func (ft *SshTransfer) connectSSH(fileCataLog *sftp.FileCatalog) error {
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
