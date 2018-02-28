package dataman

/**
    Author: luzequan
    Created: 2018-02-02 10:14:01
*/
import (
	"sync"
	"drcs/runtime/status"
	"drcs/core/interaction"
	"time"
)

type (
	// 全局调度资源池
	CarrierPool interface {
		Reset(carrierNum int) int              // 设置调度资源池数量，根据可用资源、带宽、sftp服务器负载等按需分配
		Use() interaction.Carrier              // 并发安全的使用调度资源
		Free(carrier interaction.Carrier)              // 释放调度资源
		Stop()                                 // 主动停止
		GetOneById(id int) interaction.Carrier // 根据id获取Carrier
	}
	carrierPool struct {
		capacity int                         // 调度资源规模
		count    int                         // 信使在岗数量
		usable   chan interaction.Carrier    // 信使可用（空闲）队列
		all      []interaction.Carrier       // 信使团队
		ranks    map[int]interaction.Carrier // 在运行的信使团队，根据编号
		status   int                         // 信使团队状态
		sync.RWMutex
	}
)

func NewCarrierPool() CarrierPool {
	return &carrierPool{
		status: status.RUN,
		//ranks:  make(map[int]DataMan),
		all:    make([]interaction.Carrier, 0, 5),
	}
}

func (cp *carrierPool) Reset(carrierNum int) int {
	cp.Lock()
	defer cp.Unlock()
	var wantNum int
	if carrierNum < 5 {
		wantNum = carrierNum
	} else {
		wantNum = 5
	}
	if wantNum <= 0 {
		wantNum = 1
	}
	cp.capacity = wantNum
	cp.count = 0
	cp.usable = make(chan interaction.Carrier, wantNum)
	for _, carrier := range cp.all {
		if cp.count < cp.capacity {
			cp.usable <- carrier
			cp.count++
		}
	}
	cp.status = status.RUN
	return wantNum
}

func (cp *carrierPool) Use() interaction.Carrier {
	var carrier interaction.Carrier
	for {
		cp.Lock()
		if cp.status == status.STOP {
			cp.Unlock()
			return nil
		}
		select {
		case carrier = <-cp.usable:
			//dmp.ranks[dataMan.GetId()] = dataMan
			cp.Unlock()
			return carrier
		default:
			if cp.count < cp.capacity {
				carrier = interaction.NewCross()
				cp.all = append(cp.all, carrier)
				cp.count++
				//dmp.ranks[dataMan.GetId()] = dataMan
				cp.Unlock()
				return carrier
			}
		}
		cp.Unlock()
		time.Sleep(time.Second)
	}
	return nil
}

func (cp *carrierPool) Free(carrier interaction.Carrier) {
	cp.RLock()
	defer cp.RUnlock()
	if cp.status == status.STOP {
		return
	}
	//delete(dmp.ranks, dataMan.GetId())
	cp.usable <- carrier
}

func (cp *carrierPool) Stop() {

}

func (cp *carrierPool) GetOneById(id int) interaction.Carrier {
	return nil
}