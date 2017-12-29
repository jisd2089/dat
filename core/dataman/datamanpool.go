package dataman

/**
    Author: luzequan
    Created: 2017-12-27 14:14:48
*/
import (
	"sync"
	"time"

	"dat/runtime/status"
	"dat/config"
)

type (
	DataManPool interface {
		Reset(dataFlowNum int) int // 设置数据信使数量，根据数据流产品数量按需分配
		Use() DataMan              // 并发安全的使用数据信使
		Free(DataMan)              // 释放信使资源
		Stop()
	}
	dataManPool struct {
		capacity int          // 信使团队规模
		count    int          // 信使在岗数量
		usable   chan DataMan // 信使可用队列
		all      []DataMan    // 信使团队
		status   int          // 信使团队状态
		sync.RWMutex
	}
)

func NewDataManPool() DataManPool {
	return &dataManPool{
		status: status.RUN,
		all:    make([]DataMan, 0, config.DATAMANS_CAP),
	}
}

// 根据要执行的dataFlow数量设置DataManPool
// 在二次使用Pool实例时，可根据容量高效转换
func (dmp *dataManPool) Reset(dataFlowNum int) int {
	dmp.Lock()
	defer dmp.Unlock()
	var wantNum int
	if dataFlowNum < config.DATAMANS_CAP {
		wantNum = dataFlowNum
	} else {
		wantNum = config.DATAMANS_CAP
	}
	if wantNum <= 0 {
		wantNum = 1
	}
	dmp.capacity = wantNum
	dmp.count = 0
	dmp.usable = make(chan DataMan, wantNum)
	for _, crawler := range dmp.all {
		if dmp.count < dmp.capacity {
			dmp.usable <- crawler
			dmp.count++
		}
	}
	dmp.status = status.RUN
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

func (dmp *dataManPool) Free(dataMan DataMan) {
	dmp.RLock()
	defer dmp.RUnlock()
	if dmp.status == status.STOP || !dataMan.CanStop() {
		return
	}
	dmp.usable <- dataMan
}

// 主动终止所有信使任务
func (dmp *dataManPool) Stop() {
	dmp.Lock()
	// println("CrawlerPool^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	if dmp.status == status.STOP {
		dmp.Unlock()
		return
	}
	dmp.status = status.STOP
	close(dmp.usable)
	dmp.usable = nil
	dmp.Unlock()

	for _, dataMan := range dmp.all {
		dataMan.Stop()
	}
	// println("CrawlerPool$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
}
