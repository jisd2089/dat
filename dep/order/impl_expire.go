package order

import (
	"math"
	"sync/atomic"
	"time"
	"drcs/dep/management"
)

// ExpireOrderManager 支持在指定过期时间间隔后自动更新订单信息的OrderManager装饰器类
type ExpireOrderManager struct {
	// 过期时间间隔，单位秒
	elapseS         int64
	delegated       UpdatableOrderManager
	nextExpireTimeS int64
	triableLock     *management.TriableLock
}

// GetOrderInfo 实现 OrderManager 接口
func (manager *ExpireOrderManager) GetOrderInfo() *OrderInfo {
	if manager.matchUpdateCondition() {
		manager.tryLockAndCallUpdate()
	}
	return manager.delegated.GetOrderInfo()
}

// 尝试获得更新锁，如果成功获得锁，调用delegated的更新方法
func (manager *ExpireOrderManager) tryLockAndCallUpdate() {
	if manager.triableLock.TryLock() {
		defer manager.triableLock.UnLock()
		if manager.matchUpdateCondition() {
			err := manager.delegated.Update()
			if err != nil {
				// TODO 日志
			}
			nextExpireTimeS := calcNextExpireTimeS(manager.elapseS)
			atomic.StoreInt64(&manager.nextExpireTimeS, nextExpireTimeS)
		}
	}
}

// 检查是否满足更新条件
func (manager *ExpireOrderManager) matchUpdateCondition() bool {
	nowTimeS := time.Now().Unix()
	return nowTimeS > manager.nextExpireTimeS
}

func calcNextExpireTimeS(elapseS int64) int64 {
	if elapseS <= 0 {
		return math.MaxInt64
	}

	nowTimeS := time.Now().Unix()
	return nowTimeS + elapseS
}

// NewExpireOrderManager 创建 ExpireOrderManager 实例
func NewExpireOrderManager(elapseS int64, delegated UpdatableOrderManager) *ExpireOrderManager {
	nextExpireTimeS := calcNextExpireTimeS(elapseS)
	triableLock := management.NewTriableLock()
	return &ExpireOrderManager{elapseS, delegated, nextExpireTimeS, triableLock}
}
