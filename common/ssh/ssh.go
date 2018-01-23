package ssh

import (
	"sync"
	"golang.org/x/crypto/ssh"
)

/**
    Author: luzequan
    Created: 2017-12-29 15:35:15
*/

var (
	err        error
	once       sync.Once
)

type SSHClient struct {
	Client *ssh.Client
}

func New() *SSHClient {
	return &SSHClient{}
}

func (s *SSHClient) Init(address string, config *ssh.ClientConfig) error {
	sshClient, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return err
	}
	s.Client = sshClient

	return nil
}

//刷新producer
func Refresh() {
	once.Do(func() {
		//conf := sarama.NewConfig()
		//conf.Producer.RequiredAcks = sarama.WaitForAll //等待所有备份返回ack
		//conf.Producer.Retry.Max = 10                   // 重试次数
		//brokerList := config.KAFKA_BORKERS
		//producer, err = sarama.NewSyncProducer(strings.Split(brokerList, ","), conf)
		//if err != nil {
		//	logs.Log.Error("Kafka:%v\n", err)
		//}
	})
}
