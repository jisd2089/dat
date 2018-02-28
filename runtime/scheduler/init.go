package scheduler

import (
	"sync"

	"github.com/robfig/cron"
)

var (
	scheduler *cron.Cron
	mutex     sync.Mutex
)

func Init() {
	if scheduler != nil {
		scheduler.Start()
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	if scheduler == nil {
		scheduler = cron.New()
		scheduler.Start()
		scheduler.Run()
	}

}
