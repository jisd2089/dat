package management

import "sync/atomic"

const (
	unlockedValue = 0
	lockedValue   = 1
)

// TriableLock 一个粗糙的可以尝试获取的锁
type TriableLock struct {
	intPtr *int32
}

// NewTriableLock 创建TriableLock
func NewTriableLock() *TriableLock {
	var intPtr int32
	intPtr = unlockedValue
	return &TriableLock{&intPtr}
}

// TryLock 尝试获得锁，如果成功返回true，否则立即返回false
func (lock *TriableLock) TryLock() bool {
	return atomic.CompareAndSwapInt32(lock.intPtr, unlockedValue, lockedValue)
}

// UnLock 释放锁，未获得锁的协程一定不能调用此方法
func (lock *TriableLock) UnLock() {
	atomic.CompareAndSwapInt32(lock.intPtr, lockedValue, unlockedValue)
}
