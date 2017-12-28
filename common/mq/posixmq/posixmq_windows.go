package posixmq

/**
    Author: luzequan
    Created: 2017-12-28 16:20:25
*/
import (
	"bitbucket.org/avd/go-ipc/mq"
)

type PosixMQ struct {
	mq *mq.FastMq
}

