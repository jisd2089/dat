package member

import (
	"drcs/dep/management"
	"math"
	"sync/atomic"
	"time"
)

// ExpireMemberManager 支持在指定过期时间间隔后自动更新成员信息的 MemberManager装饰器类
type ExpireMemberManager struct {
	// 过期时间间隔，单位秒
	elapseS         int64
	delegated       UpdatableMemberManager
	nextExpireTimeS int64
	triableLock     *management.TriableLock
}

// GetMemberInfo 实现 MemberManager 接口
func (manager *ExpireMemberManager) GetMemberInfo(memID string) *MemberInfo {
	if manager.matchUpdateCondition() {
		manager.tryLockAndCallUpdate()
	}
	return manager.delegated.GetMemberInfo(memID)
}

// 尝试获得更新锁，如果成功获得锁，调用delegated的更新方法
func (manager *ExpireMemberManager) tryLockAndCallUpdate() {
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
func (manager *ExpireMemberManager) matchUpdateCondition() bool {
	nowTimeS := time.Now().Unix()
	return nowTimeS > manager.nextExpireTimeS
}

// NewExpireMemberManager 创建 ExpireMemberManage 实例
func NewExpireMemberManager(elapseS int64, delegated UpdatableMemberManager) *ExpireMemberManager {
	nextExpireTimeS := calcNextExpireTimeS(elapseS)
	triableLock := management.NewTriableLock()
	return &ExpireMemberManager{elapseS, delegated, nextExpireTimeS, triableLock}
}

func calcNextExpireTimeS(elapseS int64) int64 {
	if elapseS <= 0 {
		return math.MaxInt64
	}

	nowTimeS := time.Now().Unix()
	return nowTimeS + elapseS
}
