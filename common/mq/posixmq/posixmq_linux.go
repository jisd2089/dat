package posixmq

/**
    Author: luzequan
    Created: 2017-12-28 16:19:31
*/
import (
	"bitbucket.org/avd/go-ipc/mq"
	"os"
	"github.com/getlantern/lantern/src/github.com/oxtoacart/bpool"
	"github.com/ouqiang/gocron/modules/logger"
)

var _bufferPool = bpool.NewSizedBufferPool(1000, 1024)

type PosixMQ struct {
	mq *mq.LinuxMessageQueue
}

func New(path string) (*PosixMQ, error) {
	mq, err := mq.OpenLinuxMessageQueue(path, os.O_WRONLY)
	if err != nil {
		return nil, err
	}
	return &PosixMQ{mq}, nil
}

func (pm *PosixMQ) Record() {

	buffer := _bufferPool.Get()
	defer _bufferPool.Put(buffer)
	// buffer := new(bytes.Buffer)
	//flatRecordIntoBuffer(record, buffer)
	// TODO 检查错误
	pm.mq.Send(buffer.Bytes())
}