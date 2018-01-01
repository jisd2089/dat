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
	"dat/core/dataflow"
	"dat/core/pipeline/collector"
	"time"
	"dat/runtime/cache"
	"dat/core/scheduler"
	"fmt"
	"strings"
	"reflect"
	"github.com/henrylee2cn/pholcus/logs"
	"github.com/henrylee2cn/teleport"
)

// 数据资产方
type (
	AssetNode interface {
		Init() AssetNode                                         // 初始化
		Empower() AssetNode                                      // 资产方赋权
		GetAppConf(k ...string) interface{}                      // 获取全局参数
		SetConfig(k string, v interface{}) AssetNode             // 设置全局参数
		DataFlowPrepare(original []*dataflow.DataFlow) AssetNode // 须在设置全局运行参数后Run()前调用（client模式下不调用该方法）
		Run()                                                    // 阻塞式运行直至任务完成（须在所有应当配置项配置完成后调用）
		Stop()                                                   // Offline 模式下中途终止任务（对外为阻塞式运行直至当前任务终止）
		IsRunning() bool                                         // 检查任务是否正在运行
		IsPause() bool                                           // 检查任务是否处于暂停状态
		IsStopped() bool                                         // 检查任务是否已经终止
		PauseRecover()                                           // Offline 模式下暂停\恢复任务
		Status() int                                             // 返回当前状态
		GetDataFlowLib() []*dataflow.DataFlow                    // 获取全部Dataflow种类
		GetDataFlowByName(string) *dataflow.DataFlow             // 通过名字获取某DataFlow
		GetDataFlowQueue() dataman.DataFlowQueue                 // 获取DataFlow队列接口实例
		GetOutputLib() []string                                  // 获取全部输出方式
		GetTaskBase() *distribute.TaskBase                       // 返回任务库
		distribute.Distributer                                   // 实现分布式接
	}
	NodeEntity struct {
		id           int           //资产方系统ID
		rights       []string      //资产方权利
		roleType     string        //资产方角色类型
		*cache.AppConf             // 全局配置
		*dataflow.DataFlowSpecies  //数据产品流种类
		*distribute.TaskBase       //服务器与客户端间传递任务的存储库
		dataman.DataFlowQueue      //当前任务的数产品流队列
		dataman.DataManPool        //配送回收池
		teleport.Teleport          // socket长连接双工通信接口，json数据传输
		status       int           // 运行状态
		sum          [2]uint64     // 执行计数
		takeTime     time.Duration // 执行计时
		finish       chan bool
		finishOnce   sync.Once
		canSocketLog bool
		sync.RWMutex
	}
)

// 全局唯一的核心接口实例
var AssetNodeEntity = New()

func New() AssetNode {
	return newNodeEntity()
}

func newNodeEntity() *NodeEntity {
	return &NodeEntity{}
}

// 必要的初始化
func (a *NodeEntity) Init() AssetNode {
	//a.Teleport = teleport.New()
	a.TaskBase = distribute.NewTaskBase()
	a.DataFlowQueue = dataman.NewDataFlowQueue()
	a.DataManPool = dataman.NewDataManPool()

	//switch a.AppConf.Mode {
	//case status.SERVER:
	//	if a.checkPort() {
	//		logs.Log.Informational("                                                                                               ！！当前运行模式为：[ 服务器 ] 模式！！")
	//		a.Teleport.SetAPI(distribute.MasterApi(a)).Server(":" + strconv.Itoa(a.AppConf.Port))
	//	}
	//
	//case status.CLIENT:
	//	if a.checkAll() {
	//		logs.Log.Informational("                                                                                               ！！当前运行模式为：[ 客户端 ] 模式！！")
	//		a.Teleport.SetAPI(distribute.SlaveApi(a)).Client(a.AppConf.Master, ":"+strconv.Itoa(a.AppConf.Port))
	//		// 开启节点间log打印
	//		a.canSocketLog = true
	//		go a.socketLog()
	//	}
	//case status.OFFLINE:
	//	logs.Log.Informational("                                                                                               ！！当前运行模式为：[ 单机 ] 模式！！")
	//	return a
	//default:
	//	logs.Log.Warning(" *    ——请指定正确的运行模式！——")
	//	return a
	//}
	return a
}

// 切换运行模式时使用
//func (self *NodeEntity) ReInit(mode int, port int, master string, w ...io.Writer) AssetNode {
//	if !self.IsStopped() {
//		self.Stop()
//	}
//	self.LogRest()
//	if self.Teleport != nil {
//		self.Teleport.Close()
//	}
//	// 等待结束
//	if mode == status.UNSET {
//		self = newLogic()
//		self.AppConf.Mode = status.UNSET
//		return self
//	}
//	// 重新开启
//	self = newNodeEntity().Init(mode, port, master, w...).(*NodeEntity)
//	return self
//}

// 给资产方赋权
func (a *NodeEntity) Empower() AssetNode {
	return a
}

// 获取全局参数
func (self *NodeEntity) GetAppConf(k ...string) interface{} {
	defer func() {
		if err := recover(); err != nil {
			logs.Log.Error(fmt.Sprintf("%v", err))
		}
	}()
	if len(k) == 0 {
		return self.AppConf
	}
	key := strings.Title(k[0])
	acv := reflect.ValueOf(self.AppConf).Elem()
	return acv.FieldByName(key).Interface()
}

// 设置全局参数
func (self *NodeEntity) SetAppConf(k string, v interface{}) AssetNode {
	defer func() {
		if err := recover(); err != nil {
			logs.Log.Error(fmt.Sprintf("%v", err))
		}
	}()
	if k == "Limit" && v.(int64) <= 0 {
		v = int64(dataflow.LIMIT)
	} else if k == "DockerCap" && v.(int) < 1 {
		v = int(1)
	}
	acv := reflect.ValueOf(self.AppConf).Elem()
	key := strings.Title(k)
	if acv.FieldByName(key).CanSet() {
		acv.FieldByName(key).Set(reflect.ValueOf(v))
	}

	return self
}

// DataFlowPrepare()必须在设置全局运行参数之后，Run()的前一刻执行
// original为dataflow包中未有过赋值操作的原始dataflow种类
// 已被显式赋值过的dataflow将不再重新分配Keyin
// client模式下不调用该方法
func (self *NodeEntity) DataFlowPrepare(original []*dataflow.DataFlow) AssetNode {
	self.DataFlowQueue.Reset()
	// 遍历任务
	for _, df := range original {
		spcopy := df.Copy()
		spcopy.SetPausetime(self.AppConf.Pausetime)
		if spcopy.GetLimit() == dataflow.LIMIT {
			spcopy.SetLimit(self.AppConf.Limit)
		} else {
			spcopy.SetLimit(-1 * self.AppConf.Limit)
		}
		self.DataFlowQueue.Add(spcopy)
	}
	// 遍历自定义配置
	self.DataFlowQueue.AddKeyins(self.AppConf.Keyins)
	return self
}

// 获取全部输出方式
func (self *NodeEntity) GetOutputLib() []string {
	return collector.DataOutputLib
}

// 获取全部蜘蛛种类
func (self *NodeEntity) GetDataFlowLib() []*dataflow.DataFlow {
	return self.DataFlowSpecies.Get()
}

// 通过名字获取某蜘蛛
func (self *NodeEntity) GetDataFlowByName(name string) *dataflow.DataFlow {
	return self.DataFlowSpecies.GetByName(name)
}

// 返回当前运行模式
func (self *NodeEntity) GetMode() int {
	return self.AppConf.Mode
}

// 返回任务库
func (self *NodeEntity) GetTaskBase() *distribute.TaskBase {
	return self.TaskBase
}

// 服务器客户端模式下返回节点数
func (self *NodeEntity) CountNodes() int {
	return self.Teleport.CountNodes()
}

// 获取蜘蛛队列接口实例
func (self *NodeEntity) GetDataFlowQueue() dataman.DataFlowQueue {
	return self.DataFlowQueue
}

// 系统启动
func (ne *NodeEntity) Run() {

	ne.finish = make(chan bool)
	ne.finishOnce = sync.Once{}
	// 重置计数
	ne.sum[0], ne.sum[1] = 0, 0
	// 重置计时
	ne.takeTime = 0
	// 设置状态
	ne.setStatus(status.RUN)
	defer ne.setStatus(status.STOPPED)
	// 任务执行
	switch ne.AppConf.Mode {
	case status.OFFLINE:
		ne.offline()
	case status.SERVER:
		ne.server()
	case status.CLIENT:
		ne.client()
	default:
		return
	}
	<-ne.finish
}

// 返回当前运行状态
func (a *NodeEntity) Status() int {
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()
	return a.status
}

// 开始执行任务
func (a *NodeEntity) exec() {
	count := a.DataFlowQueue.Len()

	a.DataManPool.Reset(count)
}

// 任务执行
func (a *NodeEntity) goRun(count int) {
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

// Offline 模式下中途终止任务
func (self *NodeEntity) Stop() {
	if self.status == status.STOPPED {
		return
	}
	if self.status != status.STOP {
		// 不可颠倒停止的顺序
		self.setStatus(status.STOP)
		// println("scheduler.Stop()")
		scheduler.Stop()
		// println("self.DataManPool.Stop()")
		self.DataManPool.Stop()
	}
	// println("wait self.IsStopped()")
	for !self.IsStopped() {
		time.Sleep(time.Second)
	}
}

// 检查任务是否正在运行
func (self *NodeEntity) IsRunning() bool {
	return self.status == status.RUN
}

// 检查任务是否处于暂停状态
func (self *NodeEntity) IsPause() bool {
	return self.status == status.PAUSE
}

// 检查任务是否已经终止
func (self *NodeEntity) IsStopped() bool {
	return self.status == status.STOPPED
}

// 返回当前运行状态
func (self *NodeEntity) setStatus(status int) {
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()
	self.status = status
}

// ******************************************** 私有方法 ************************************************* \\
// 离线模式运行
func (self *NodeEntity) offline() {
	self.exec()
}

// 服务器模式运行，必须在DataFlowPrepare()执行之后调用才可以成功添加任务
// 生成的任务与自身当前全局配置相同
func (self *NodeEntity) server() {
	// 标记结束
	defer func() {
		self.finishOnce.Do(func() { close(self.finish) })
	}()

	// 便利添加任务到库
	tasksNum, dataFlowsNum := self.addNewTask()

	if tasksNum == 0 {
		return
	}

	// 打印报告
	logs.Log.Informational(" * ")
	logs.Log.Informational(` *********************************************************************************************************************************** `)
	logs.Log.Informational(" * ")
	logs.Log.Informational(" *                               —— 本次成功添加 %v 条任务，共包含 %v 条采集规则 ——", tasksNum, dataFlowsNum)
	logs.Log.Informational(" * ")
	logs.Log.Informational(` *********************************************************************************************************************************** `)
}

// 服务器模式下，生成task并添加至库
func (self *NodeEntity) addNewTask() (tasksNum, dataFlowsNum int) {
	length := self.DataFlowQueue.Len()
	t := distribute.Task{}
	// 从配置读取字段
	self.setTask(&t)

	for i, sp := range self.DataFlowQueue.GetAll() {

		t.DataFlows = append(t.DataFlows, map[string]string{"name": sp.GetName(), "keyin": sp.GetKeyin()})
		dataFlowsNum++

		// 每十个蜘蛛存为一个任务
		if i > 0 && i%10 == 0 && length > 10 {
			// 存入
			one := t
			self.TaskBase.Push(&one)
			// logs.Log.App(" *     [新增任务]   详情： %#v", *t)

			tasksNum++

			// 清空dataflow
			t.DataFlows = []map[string]string{}
		}
	}

	if len(t.DataFlows) != 0 {
		// 存入
		one := t
		self.TaskBase.Push(&one)
		tasksNum++
	}
	return
}

// 客户端模式运行
func (self *NodeEntity) client() {
	// 标记结束
	defer func() {
		self.finishOnce.Do(func() { close(self.finish) })
	}()

	for {
		// 从任务库获取一个任务
		t := self.downTask()

		if self.Status() == status.STOP || self.Status() == status.STOPPED {
			return
		}

		// 准备运行
		self.taskToRun(t)

		// 重置计数
		self.sum[0], self.sum[1] = 0, 0
		// 重置计时
		self.takeTime = 0

		// 执行任务
		self.exec()
	}
}

// 客户端模式下获取任务
func (self *NodeEntity) downTask() *distribute.Task {
ReStartLabel:
	if self.Status() == status.STOP || self.Status() == status.STOPPED {
		return nil
	}
	if self.CountNodes() == 0 && self.TaskBase.Len() == 0 {
		time.Sleep(time.Second)
		goto ReStartLabel
	}

	if self.TaskBase.Len() == 0 {
		self.Request(nil, "task", "")
		for self.TaskBase.Len() == 0 {
			if self.CountNodes() == 0 {
				goto ReStartLabel
			}
			time.Sleep(time.Second)
		}
	}
	return self.TaskBase.Pull()
}

// client模式下从task准备运行条件
func (self *NodeEntity) taskToRun(t *distribute.Task) {
	// 清空历史任务
	self.DataFlowQueue.Reset()

	// 更改全局配置
	self.setAppConf(t)

	// 初始化蜘蛛队列
	for _, n := range t.DataFlows {
		sp := self.GetDataFlowByName(n["name"])
		if sp == nil {
			continue
		}
		spcopy := sp.Copy()
		spcopy.SetPausetime(t.Pausetime)
		if spcopy.GetLimit() > 0 {
			spcopy.SetLimit(t.Limit)
		} else {
			spcopy.SetLimit(-1 * t.Limit)
		}
		if v, ok := n["keyin"]; ok {
			spcopy.SetKeyin(v)
		}
		self.DataFlowQueue.Add(spcopy)
	}
}

// 设置任务运行时公共配置
func (self *NodeEntity) setAppConf(task *distribute.Task) {
	self.AppConf.ThreadNum = task.ThreadNum
	self.AppConf.Pausetime = task.Pausetime
	self.AppConf.OutType = task.OutType
	self.AppConf.DockerCap = task.DockerCap
	self.AppConf.SuccessInherit = task.SuccessInherit
	self.AppConf.FailureInherit = task.FailureInherit
	self.AppConf.Limit = task.Limit
	self.AppConf.ProxyMinute = task.ProxyMinute
	self.AppConf.Keyins = task.Keyins
}
func (self *NodeEntity) setTask(task *distribute.Task) {
	task.ThreadNum = self.AppConf.ThreadNum
	task.Pausetime = self.AppConf.Pausetime
	task.OutType = self.AppConf.OutType
	task.DockerCap = self.AppConf.DockerCap
	task.SuccessInherit = self.AppConf.SuccessInherit
	task.FailureInherit = self.AppConf.FailureInherit
	task.Limit = self.AppConf.Limit
	task.ProxyMinute = self.AppConf.ProxyMinute
	task.Keyins = self.AppConf.Keyins
}
