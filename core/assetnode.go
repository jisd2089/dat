package assetnode

/**
    Author: luzequan
    Created: 2017-12-27 13:33:06
*/
import (
	"sync"
	"dat/core/distribute"
	"dat/core/dataman"
	"dat/runtime/status"
)

// 数据资产方
type (
	AssetNode interface {
		Init() AssetNode      // 初始化
		Empower() AssetNode   // 资产方赋权
		SetConfig() AssetNode // 设置全局参数
		Run()
		Status() int // 返回当前状态
	}
	assetNode struct {
		id       int          //资产方系统ID
		rights   []string     //资产方权利
		roleType string       //资产方角色类型
		*distribute.TaskBase  //服务器与客户端间传递任务的存储库
		dataman.DataFlowQueue //当前任务的数产品流队列
		dataman.DataManPool   //配送回收池
		status   int          // 运行状态
		sync.RWMutex
	}
)

// 全局唯一的核心接口实例
var LogicApp = New()

func New() AssetNode {
	return &assetNode{
	}
}

// 必要的初始化
func (a *assetNode) Init() AssetNode {
	return a
}

// 给资产方赋权
func (a *assetNode) Empower() AssetNode {
	return a
}

// 设置全局参数
func (a *assetNode) SetConfig() AssetNode {
	defer func() {
		if err := recover(); err != nil {
			//logs.Log.Error(fmt.Sprintf("%v", err))
		}
	}()

	return a
}

// 系统启动
func (a *assetNode) Run() {

}

// 返回当前运行状态
func (a *assetNode) Status() int {
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()
	return a.status
}

// 开始执行任务
func (a *assetNode) exec() {
	count := a.DataFlowQueue.Len()

	a.DataManPool.Reset(count)
}

// 任务执行
func (a *assetNode) goRun(count int) {
	m := a.DataManPool.Use()
	if m != nil {
		go func(i int, c dataman.DataMan) {
			// 执行并返回结果消息
			m.Init(a.DataFlowQueue.GetByIndex(i)).Run()
			// 任务结束后回收该蜘蛛
			a.RWMutex.RLock()
			if a.status != status.STOP {
				a.DataManPool.Free(c)
			}
			a.RWMutex.RUnlock()
		}(i, m)
	}
}
