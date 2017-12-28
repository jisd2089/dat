package dataman

/**
    Author: luzequan
    Created: 2017-12-27 14:14:48
*/
import (
	"sync"
	"time"

	"dat/runtime/status"
)

type (
	DataManPool interface {
		Reset(spiderNum int) int
		Use() DataMan
	}
	dataManPool struct {
		capacity int
		count    int
		usable   chan DataMan
		all      []DataMan
		status   int
		sync.RWMutex
	}
)

func NewDataManPool() DataManPool {
	return &dataManPool{}
}

// 根据要执行的蜘蛛数量设置CrawlerPool
// 在二次使用Pool实例时，可根据容量高效转换
func (dmp *dataManPool) Reset(dataFlowNum int) int {
	dmp.Lock()
	defer dmp.Unlock()
	var wantNum int
	//if spiderNum < config.CRAWLS_CAP {
	//	wantNum = spiderNum
	//} else {
	//	wantNum = config.CRAWLS_CAP
	//}
	//if wantNum <= 0 {
	//	wantNum = 1
	//}
	//dmp.capacity = wantNum
	//dmp.count = 0
	//dmp.usable = make(chan Crawler, wantNum)
	//for _, crawler := range dmp.all {
	//	if dmp.count < dmp.capacity {
	//		dmp.usable <- crawler
	//		dmp.count++
	//	}
	//}
	//dmp.status = status.RUN
	return wantNum
}

// 并发安全地使用资源
func (dmp *dataManPool) Use() DataMan {
	var crawler DataMan
	for {
		dmp.Lock()
		if dmp.status == status.STOP {
			dmp.Unlock()
			return nil
		}
		select {
		case crawler = <-dmp.usable:
			dmp.Unlock()
			return crawler
		default:
			if dmp.count < dmp.capacity {
				crawler = New(dmp.count)
				dmp.all = append(dmp.all, crawler)
				dmp.count++
				dmp.Unlock()
				return crawler
			}
		}
		dmp.Unlock()
		time.Sleep(time.Second)
	}
	return nil
}