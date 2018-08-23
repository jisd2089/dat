package dataman

/**
    Author: luzequan
    Created: 2017-12-27 14:14:48
*/
import (
	"sync"
	"time"

	"drcs/runtime/status"
	"drcs/config"
	"drcs/common/util"
)

type (
	DataManPool interface {
		Reset(dataBoxNum int) int  // 设置数据信使数量，根据数据流产品数量按需分配
		Use() DataMan              // 并发安全的使用数据信使
		UseOne() DataMan           // 并发安全的使用数据信使
		Free(DataMan)              // 释放信使资源
		Stop()                     // 主动停止
		GetOneById(id int) DataMan // 根据id获取dataman
	}
	dataManPool struct {
		capacity int             // 信使团队规模
		count    int             // 信使在岗数量
		usable   chan DataMan    // 信使可用（空闲）队列
		all      []DataMan       // 信使团队
		ranks    map[int]DataMan // 在运行的信使团队，根据编号
		status   int             // 信使团队状态
		sync.RWMutex
	}
)

func NewDataManPool() DataManPool {
	return &dataManPool{
		status: status.RUN,
		ranks:  make(map[int]DataMan),
		all:    make([]DataMan, 0, config.DATAMANS_CAP),
	}
}

var (
	idInc  = util.NewAutoInc(0, 1)
	dmPool = &sync.Pool{
		New: func() interface{} {
			return New(idInc.Id())
		},
	}
)

// 根据要执行的dataBox数量设置DataManPool
// 在二次使用Pool实例时，可根据容量高效转换
func (dmp *dataManPool) Reset(dataBoxNum int) int {
	dmp.Lock()
	defer dmp.Unlock()
	var wantNum int
	if dataBoxNum < config.DATAMANS_CAP {
		wantNum = dataBoxNum
	} else {
		wantNum = config.DATAMANS_CAP
	}
	if wantNum <= 0 {
		wantNum = 1
	}
	dmp.capacity = wantNum
	dmp.count = 0
	dmp.usable = make(chan DataMan, wantNum)
	for _, dataMan := range dmp.all {
		if dmp.count < dmp.capacity {
			dmp.usable <- dataMan
			dmp.count++
		}
	}
	dmp.status = status.RUN
	return wantNum
}

// 并发安全地使用资源
func (dmp *dataManPool) Use() DataMan {
	var dataMan DataMan
	//fmt.Println("dataManPool count:", dmp.count)
	for {
		dmp.Lock()
		if dmp.status == status.STOP {
			dmp.Unlock()
			return nil
		}
		select {
		case dataMan = <-dmp.usable:
			dmp.ranks[dataMan.GetId()] = dataMan
			dmp.Unlock()
			return dataMan
		default:
			if dmp.count < dmp.capacity {
				dataMan = New(dmp.count)
				dmp.all = append(dmp.all, dataMan)
				dmp.count++
				dmp.ranks[dataMan.GetId()] = dataMan
				dmp.Unlock()
				return dataMan
			}
		}
		dmp.Unlock()
		time.Sleep(time.Second)
	}
	return nil
}

func (dmp *dataManPool) UseOne() DataMan {
	return dmPool.Get().(DataMan)
}

func (dmp *dataManPool) Free(dataMan DataMan) {
	dmp.RLock()
	defer dmp.RUnlock()
	if dmp.status == status.STOP || !dataMan.CanStop() {
		return
	}
	//delete(dmp.ranks, dataMan.GetId())
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

// 根据id获取dataman
func (dmp *dataManPool) GetOneById(id int) DataMan {
	dmp.Lock()
	defer dmp.Unlock()
	return dmp.ranks[id]
}
