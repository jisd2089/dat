package posixmq

/**
    Author: luzequan
    Created: 2017-12-28 16:20:25
*/
import (
	"bitbucket.org/avd/go-ipc/mq"
	"os"
	"time"
	"fmt"
)

type PosixMQ struct {
	mq *mq.FastMq
}

func New(path string) (*PosixMQ, error) {
	mq, err := mq.OpenFastMq(path, os.O_WRONLY)
	if err != nil {
		return nil, err
	}
	return &PosixMQ{mq}, nil
}

func (pm *PosixMQ) Record(record *Record) {
	fmt.Println("posixmq record start~")

	tm := time.Now()
	record.rdate = tm.Format("20060102")
	record.rtime = tm.Format("150405")
	record.stepCount = len(record.StepInfos)

	buffer := _bufferPool.Get()
	defer _bufferPool.Put(buffer)

	flatRecordIntoBuffer(record, buffer)

	fmt.Println("posixmq record before send~", string(buffer.Bytes()))
	err := pm.mq.Send(buffer.Bytes())
	if err != nil {
		fmt.Println("mq send err: ", err.Error())
	}
	fmt.Println("posixmq record end~")
}

