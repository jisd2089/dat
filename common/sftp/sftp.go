package sftp

import (
	"path"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"fmt"
	"os"
	"sync"
	"path/filepath"
)

/**
    Author: luzequan
    Created: 2017-12-29 15:30:55
*/

type SFTPClient struct {
	client *sftp.Client
	lock   sync.RWMutex
}

func New() *SFTPClient {
	return &SFTPClient{}
}

func (sc *SFTPClient) Init(sshClient *ssh.Client) error {
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	sc.client = sftpClient
	return nil
}

func (sc *SFTPClient) CloseSession() {
	sc.client.Close()
}

/**
在远程目录下创建文件夹
 */
func (sc *SFTPClient) RemoteMkdir(remoteDir, newRemoteFolder string) error {
	sc.lock.RLock()
	defer sc.lock.RUnlock()
	remotePath := path.Join(remoteDir, newRemoteFolder)
	//defer sc.CloseSession()
	_, err := sc.client.Stat(remotePath)
	if err == nil {
		return nil
	}
	err = sc.client.Mkdir(remotePath)
	if err != nil {
		return fmt.Errorf("%s", "sftp create folder err: "+err.Error())
	}
	return nil
}

/**
删除远程目录文件 (慎用)
*/
func (sc *SFTPClient) RemoteRM(remoteDir, remoteFileName string) error {
	sc.lock.RLock()
	defer sc.lock.RUnlock()

	//defer sc.CloseSession()

	err := sc.client.Remove(path.Join(remoteDir, remoteFileName))
	if err != nil {
		return fmt.Errorf("%s", "sftp remove file err: "+err.Error())
	}
	return nil
}

/**
从远程服务器指定目录拉取文件
*/
func (sc *SFTPClient) RemoteGet(fc *FileCatalog) error {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("sftp RemoteGet recover error: ", err)
		}
	}()

	sc.lock.RLock()
	defer sc.lock.RUnlock()

	//open 远程文件
	//defer sc.CloseSession()
	remoteFile, err := sc.client.Open(path.Join(fc.RemoteDir, fc.RemoteFileName))
	defer remoteFile.Close()
	if err != nil {
		fmt.Println("RemoteGet fc.RemoteFileName", fc.RemoteFileName)
		fmt.Println("RemoteGet error >>>>>>>>>>>>>>>", fmt.Errorf("RemoteGet Open Remote File [%s] Error: %s", fc.RemoteFileName, err.Error()).Error())
		return fmt.Errorf("RemoteGet Open Remote File [%s] Error: %s", fc.RemoteFileName, err.Error())
	}

	//创建本地文件
	if !isDirExists(fc.LocalDir) {
		os.MkdirAll(filepath.Clean(fc.LocalDir), 0777)
	}
	localFile, err := os.Create(path.Join(fc.LocalDir, fc.LocalFileName))
	defer localFile.Close()
	if err != nil {
		return fmt.Errorf("RemoteGet Create Local File [%s] Error: %s", fc.LocalFileName, err.Error())
	}

	//远程文件写入
	if _, err = remoteFile.WriteTo(localFile); err != nil {
		return fmt.Errorf("RemoteGet Write To Local File From Remote Error: %s", err.Error())
	}

	return nil
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

/**
SFTP推送文件到远程服务器
*/
func (sc *SFTPClient) RemotePut(fc *FileCatalog) error {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("sftp RemotePut recover error: ", err)
		}
	}()

	sc.lock.RLock()
	defer sc.lock.RUnlock()

	//defer sc.CloseSession()

	//参数验证
	if err := sc.Verify(fc); err != nil {
		return err
	}

	//创建远程文件
	dstFile, err := sc.client.Create(path.Join(fc.RemoteDir, fc.RemoteFileName))
	defer dstFile.Close()
	if err != nil {
		return fmt.Errorf("RemotePut Create Remote File [%s] Error: %s", fc.RemoteFileName, err.Error())
	}

	//打开本地文件
	srcFile, err := os.Open(path.Join(fc.LocalDir, fc.LocalFileName))
	defer srcFile.Close()
	if err != nil {
		return fmt.Errorf("RemotePut Open Local File [%s] Error: %s", fc.LocalFileName, err.Error())
	}
	//上传写入文件
	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return fmt.Errorf("RemotePut DstFile ReadFrom SrcFile Error: %s", err.Error())
	}

	return nil
}

/**
列出远程文件目录
*/
func (sc *SFTPClient) RemoteLS(remoteDir string) ([]os.FileInfo, error) {
	sc.lock.RLock()
	defer sc.lock.RUnlock()

	//defer sc.CloseSession()
	fileInfos, err := sc.client.ReadDir(remoteDir)
	if err != nil {
		return nil, fmt.Errorf("RemoteLS ReadDir [%s] Error: %s", remoteDir, err.Error())
	}
	return fileInfos, nil
}

/**
参数验证
*/
func (sc *SFTPClient) Verify(fc *FileCatalog) error {
	if fc.LocalDir == "" {
		return fmt.Errorf("%s", "LocalDir is nil")
	}
	if fc.LocalFileName == "" {
		return fmt.Errorf("%s", "LocalFileName is nil")
	}
	if fc.RemoteDir == "" {
		return fmt.Errorf("%s", "RemoteDir is nil")
	}
	if fc.RemoteFileName == "" {
		return fmt.Errorf("%s", "RemoteFileName is nil")
	}
	return nil
}
