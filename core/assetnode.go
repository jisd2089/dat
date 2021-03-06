package assetnode

/**
    Author: luzequan
    Created: 2017-12-27 13:33:06
*/
import (
	"sync"
	"drcs/core/distribute"
	"drcs/core/dataman"
	"drcs/runtime/status"
	"drcs/core/databox"
	"drcs/core/pipeline/collector"
	"time"
	"drcs/runtime/cache"
	"drcs/core/scheduler"
	"drcs/core/pipeline"
	"fmt"
	"strings"
	"reflect"

	"github.com/henrylee2cn/teleport"
	"drcs/core/interaction/response"
)

// 数据资产方
type (
	AssetNode interface {
		Init() AssetNode                                                         // 初始化
		Empower() AssetNode                                                      // 资产方赋权
		GetConfig(k ...string) interface{}                                       // 获取全局参数
		SetConfig(k string, v interface{}) AssetNode                             // 设置全局参数
		DataBoxPrepare(original []*databox.DataBox) AssetNode                    // 须在设置全局运行参数后Run()前调用（client模式下不调用该方法）
		PushDataBox(original []*databox.DataBox) AssetNode                       // 将DataBox放入channel
		PushActiveDataBox(original *databox.DataBox) *databox.DataBox            // 将DataBox放入active channel
		Run()                                                                    // 阻塞式运行直至任务完成（须在所有应当配置项配置完成后调用）
		SyncRun()                                                                // 同步运行ActiveBox
		RunActiveBox(b *databox.DataBox, obj interface{}) response.DataResponse // 执行ActiveBox请求，同步返回
		StopActiveBox(b *databox.DataBox)                                        // 停止Active DataBox
		Stop()                                                                   // Offline 模式下中途终止任务（对外为阻塞式运行直至当前任务终止）
		IsRunning() bool                                                         // 检查任务是否正在运行
		IsPause() bool                                                           // 检查任务是否处于暂停状态
		IsStopped() bool                                                         // 检查任务是否已经终止
		PauseRecover()                                                           // Offline 模式下暂停\恢复任务
		Status() int                                                             // 返回当前状态
		GetDataBoxLib() []*databox.DataBox                                       // 获取全部databox种类
		GetDataBoxByName(string) *databox.DataBox                                // 通过名字获取某DataBox
		GetActiveDataBoxByName(string) *databox.DataBox                          // 通过名字获取某活跃DataBox
		GetDataBoxQueue() dataman.DataBoxQueue                                   // 获取DataBox队列接口实例
		GetDataManPool() dataman.DataManPool                                     // 获取DataManPool
		GetOutputLib() []string                                                  // 获取全部输出方式
		GetTaskBase() *distribute.TaskBase                                       // 返回任务库
		distribute.Distributer                                                   // 实现分布式接口
	}
	NodeEntity struct {
		id           int           // 资产方系统ID
		rights       []string      // 资产方权利
		roleType     string        // 资产方角色类型
		*cache.AppConf             // 全局配置
		*databox.DataBoxSpecies    // 数据产品流种类
		*databox.DataBoxActivites  // DataBox活跃列表
		*distribute.TaskBase       // 服务器与客户端间传递任务的存储库
		dataman.DataBoxQueue       // 当前任务的数据产品流队列
		dataman.DataManPool        // 配送回收池
		dataman.CarrierPool        // 传输资源池
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
var (
	AssetNodeEntity = New()
	newNodeEntity   *NodeEntity
	once            sync.Once
)

func New() AssetNode {
	return getNewNodeEntity()
}

func getNewNodeEntity() *NodeEntity {
	once.Do(func() {
		newNodeEntity = &NodeEntity{
			AppConf:          cache.Task,
			DataBoxSpecies:   databox.Species,
			DataBoxActivites: databox.Activites,
			status:           status.STOPPED,
			Teleport:         teleport.New(),
			TaskBase:         distribute.NewTaskBase(),
			DataBoxQueue:     dataman.NewDataBoxQueue(),
			DataManPool:      dataman.NewDataManPool(),
			//CarrierPool:      dataman.NewCarrierPool(),
		}
	})
	return newNodeEntity
}

// 必要的初始化
func (a *NodeEntity) Init() AssetNode {
	a.TaskBase = distribute.NewTaskBase()
	a.DataBoxQueue = dataman.NewDataBoxQueue()
	a.DataManPool = dataman.NewDataManPool()

	return a
}

//func (n *NodeEntity) ReInit(mode int, port int, master string, w ...io.Writer) AssetNode {
//	if !n.IsStopped() {
//		n.Stop()
//	}
//	//n.LogRest()
//	if n.Teleport != nil {
//		n.Teleport.Close()
//	}
//	// 等待结束
//	if mode == status.UNSET {
//		n = newLogic()
//		n.AppConf.Mode = status.UNSET
//		return n
//	}
//	// 重新开启
//	n = newNodeEntity().Init().(*NodeEntity)
//	return n
//}

// 给资产方赋权
func (n *NodeEntity) Empower() AssetNode {
	return n
}

// 获取全局参数
func (self *NodeEntity) GetConfig(k ...string) interface{} {
	defer func() {
		if err := recover(); err != nil {
			//logs.Log.Error(fmt.Sprintf("%v", err))
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
func (self *NodeEntity) SetConfig(k string, v interface{}) AssetNode {
	defer func() {
		if err := recover(); err != nil {
			//logs.Log.Error(fmt.Sprintf("%v", err))
		}
	}()
	if k == "Limit" && v.(int64) <= 0 {
		v = int64(databox.LIMIT)
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

// DataBoxPrepare()必须在设置全局运行参数之后，Run()的前一刻执行
// original为databox包中未有过赋值操作的原始databox种类
// 已被显式赋值过的databox将不再重新分配Keyin
// client模式下不调用该方法
func (self *NodeEntity) DataBoxPrepare(original []*databox.DataBox) AssetNode {
	self.DataBoxQueue.Reset()
	// 遍历任务
	for _, df := range original {
		dfcopy := df.Copy()
		dfcopy.SetPausetime(self.AppConf.Pausetime)
		if dfcopy.GetLimit() == databox.LIMIT {
			dfcopy.SetLimit(self.AppConf.Limit)
		} else {
			dfcopy.SetLimit(-1 * self.AppConf.Limit)
		}
		self.DataBoxQueue.Add(dfcopy)
	}
	// 遍历自定义配置
	self.DataBoxQueue.AddKeyins(self.AppConf.Keyins)
	return self
}

func (self *NodeEntity) PushDataBox(original []*databox.DataBox) AssetNode {
	// 遍历任务
	for _, b := range original {
		dfcopy := b.Copy()
		dfcopy.SetPausetime(self.AppConf.Pausetime)
		if dfcopy.GetLimit() == databox.LIMIT {
			dfcopy.SetLimit(self.AppConf.Limit)
		} else {
			dfcopy.SetLimit(-1 * self.AppConf.Limit)
		}
		self.DataBoxQueue.AddChan(dfcopy)

		b.Refresh()
	}
	// 遍历自定义配置
	//self.DataBoxQueue.AddKeyins(self.AppConf.Keyins)
	return self
}

func (self *NodeEntity) PushActiveDataBox(original *databox.DataBox) *databox.DataBox {
	// 拷贝任务
	dfcopy := original.Copy()
	dfcopy.SetPausetime(self.AppConf.Pausetime)
	dfcopy.ActiveWG = &sync.WaitGroup{}
	if dfcopy.GetLimit() == databox.LIMIT {
		dfcopy.SetLimit(self.AppConf.Limit)
	} else {
		dfcopy.SetLimit(-1 * self.AppConf.Limit)
	}
	self.DataBoxQueue.AddActiveChan(dfcopy)
	// 遍历自定义配置
	//self.DataBoxQueue.AddKeyins(self.AppConf.Keyins)
	return dfcopy
}

// 获取全部输出方式
func (self *NodeEntity) GetOutputLib() []string {
	return collector.DataOutputLib
}

// 获取全部databox种类
func (self *NodeEntity) GetDataBoxLib() []*databox.DataBox {
	return self.DataBoxSpecies.Get()
}

// 通过名字获取某databox
func (self *NodeEntity) GetDataBoxByName(name string) *databox.DataBox {
	return self.DataBoxSpecies.GetByName(name)
}

// 通过名字获取某活跃databox
func (self *NodeEntity) GetActiveDataBoxByName(name string) *databox.DataBox {
	return self.DataBoxActivites.GetByName(name)
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

// 获取databox队列接口实例
func (self *NodeEntity) GetDataBoxQueue() dataman.DataBoxQueue {
	return self.DataBoxQueue
}

func (ne *NodeEntity) GetDataManPool() dataman.DataManPool {
	return ne.DataManPool
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
	ne.exec()

	//<-ne.finish //TODO 节点持续运行不退出
}

// 系统启动, 同步返回
func (ne *NodeEntity) SyncRun() {
	// 任务执行
	ne.syncExec()

}

// 返回当前运行状态
func (a *NodeEntity) Status() int {
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()
	return a.status
}

// 开始执行任务
func (ne *NodeEntity) exec() {
	//count := ne.DataBoxQueue.Len()

	cache.ResetPageCount()
	// 刷新输出方式的状态
	pipeline.RefreshOutput()
	// 初始化资源队列
	scheduler.Init()

	// 设置数据信使队列
	//dataManCap := ne.DataManPool.Reset(count)
	dataManCap := ne.DataManPool.Reset(10)
	//ne.CarrierPool.Reset(5)

	fmt.Println(" *     DataManPool池容量为 %v\n", dataManCap)

	// 开始计时
	cache.StartTime = time.Now()

	//TODO 根据节点支持业务类型启动 两类DataBox
	// goroutine 1
	go ne.runDataBox()
	// goroutine 2
	go ne.goSyncRun()
}

// 开始执行任务，同步开始
func (ne *NodeEntity) syncExec() {

	// 根据模式选择合理的并发
	go ne.goSyncRun()
}

// 任务执行
func (ne *NodeEntity) goRun(count int) {
	// 执行任务
	var i int
	for i = 0; i < count && ne.Status() != status.STOP; i++ {
	pause:
		if ne.IsPause() {
			time.Sleep(time.Second)
			goto pause
		}
		// 从数据信使队列取出空闲信使，并发执行
		m := ne.DataManPool.Use()
		if m != nil {
			go func(i int, m dataman.DataMan) {
				// 执行并返回结果消息
				m.Init(ne.DataBoxQueue.GetByIndex(i)).Run()
				// 任务结束后回收该信使
				ne.RWMutex.RLock()
				if ne.status != status.STOP {
					ne.DataManPool.Free(m)
				}
				ne.RWMutex.RUnlock()
			}(i, m)
		}
	}
	// 监控结束任务
	for ii := 0; ii < i; ii++ {
		s := <-cache.ReportChan
		if (s.DataNum == 0) && (s.FileNum == 0) {
			//logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   无采集结果，用时 %v！\n", s.DataBoxName, s.Keyin, s.Time)
			continue
		}
		//logs.Log.Informational(" * ")
		switch {
		case s.DataNum > 0 && s.FileNum == 0:
			//logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   共采集数据 %v 条，用时 %v！\n",
			//	s.DataBoxName, s.Keyin, s.DataNum, s.Time)
		case s.DataNum == 0 && s.FileNum > 0:
			//logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   共下载文件 %v 个，用时 %v！\n",
			//	s.DataBoxName, s.Keyin, s.FileNum, s.Time)
		default:
			//logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   共采集数据 %v 条 + 下载文件 %v 个，用时 %v！\n",
			//	s.DataBoxName, s.Keyin, s.DataNum, s.FileNum, s.Time)
		}

		ne.sum[0] += s.DataNum
		ne.sum[1] += s.FileNum
	}

	// 总耗时
	ne.takeTime = time.Since(cache.StartTime)
	var prefix = func() string {
		if ne.Status() == status.STOP {
			return "任务中途取消："
		}
		return "本次"
	}()
	// 打印总结报告
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)
	fmt.Println(` *********************************************************************************************************************************** `)
	//logs.Log.Informational(" * ")
	switch {
	case ne.sum[0] > 0 && ne.sum[1] == 0:
		fmt.Println(" *                            —— %s合计采集【数据 %v 条】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, ne.sum[0], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
		//logs.Log.App(" *                            —— %s合计采集【数据 %v 条】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
		//	prefix, ne.sum[0], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
	case ne.sum[0] == 0 && ne.sum[1] > 0:
		fmt.Println(" *                            —— %s合计采集【文件 %v 个】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, ne.sum[1], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
		//logs.Log.App(" *                            —— %s合计采集【文件 %v 个】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
		//	prefix, ne.sum[1], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
	case ne.sum[0] == 0 && ne.sum[1] == 0:
		fmt.Println(" *                            —— %s无采集结果，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
		//logs.Log.App(" *                            —— %s无采集结果，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
		//	prefix, cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
	default:
		fmt.Println(" *                            —— %s合计采集【数据 %v 条 + 文件 %v 个】，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, ne.sum[0], ne.sum[1], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
		//logs.Log.App(" *                            —— %s合计采集【数据 %v 条 + 文件 %v 个】，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
		//	prefix, ne.sum[0], ne.sum[1], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
	}
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)
	fmt.Println(` *********************************************************************************************************************************** `)

	// 单机模式并发运行，需要标记任务结束
	if ne.AppConf.Mode == status.OFFLINE {
		//ne.LogRest()
		ne.finishOnce.Do(func() { close(ne.finish) })
	}
}

// 任务执行
func (ne *NodeEntity) runDataBox() {
	// 执行任务
	//pause:
	//	if ne.IsPause() {
	//		time.Sleep(time.Second)
	//		goto pause
	//	}
	for {
		b := ne.DataBoxQueue.GetOne()
		go ne.goRunDataBox(b)
	}
}

func (ne *NodeEntity) goRunDataBox(b *databox.DataBox) {
	// 从数据信使队列取出空闲信使，并发执行
	m := ne.DataManPool.Use()
	if m != nil {
		// 执行并返回结果消息
		m.Init(b).Run()
		// 任务结束后回收该信使
		ne.RWMutex.RLock()
		if ne.status != status.STOP {
			//m.Stop()
			ne.DataManPool.Free(m)
		}
		ne.RWMutex.RUnlock()
	}
}

func (ne *NodeEntity) runMonitorTask(i int) {
	// 监控结束任务
	for ii := 0; ii < i; ii++ {
		s := <-cache.ReportChan
		if (s.DataNum == 0) && (s.FileNum == 0) {
			//logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   无采集结果，用时 %v！\n", s.DataBoxName, s.Keyin, s.Time)
			continue
		}
		//logs.Log.Informational(" * ")
		switch {
		case s.DataNum > 0 && s.FileNum == 0:
			//logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   共采集数据 %v 条，用时 %v！\n",
			//	s.DataBoxName, s.Keyin, s.DataNum, s.Time)
		case s.DataNum == 0 && s.FileNum > 0:
			//logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   共下载文件 %v 个，用时 %v！\n",
			//	s.DataBoxName, s.Keyin, s.FileNum, s.Time)
		default:
			//logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   共采集数据 %v 条 + 下载文件 %v 个，用时 %v！\n",
			//	s.DataBoxName, s.Keyin, s.DataNum, s.FileNum, s.Time)
		}

		ne.sum[0] += s.DataNum
		ne.sum[1] += s.FileNum
	}

	// 总耗时
	ne.takeTime = time.Since(cache.StartTime)
	var prefix = func() string {
		if ne.Status() == status.STOP {
			return "任务中途取消："
		}
		return "本次"
	}()

	// 打印总结报告
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)
	fmt.Println(` *********************************************************************************************************************************** `)
	//logs.Log.Informational(" * ")
	switch {
	case ne.sum[0] > 0 && ne.sum[1] == 0:
		fmt.Println(" *                            —— %s合计采集【数据 %v 条】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, ne.sum[0], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
		//logs.Log.App(" *                            —— %s合计采集【数据 %v 条】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
		//	prefix, ne.sum[0], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
	case ne.sum[0] == 0 && ne.sum[1] > 0:
		fmt.Println(" *                            —— %s合计采集【文件 %v 个】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, ne.sum[1], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
		//logs.Log.App(" *                            —— %s合计采集【文件 %v 个】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
		//	prefix, ne.sum[1], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
	case ne.sum[0] == 0 && ne.sum[1] == 0:
		fmt.Println(" *                            —— %s无采集结果，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
		//logs.Log.App(" *                            —— %s无采集结果，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
		//	prefix, cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
	default:
		fmt.Println(" *                            —— %s合计采集【数据 %v 条 + 文件 %v 个】，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, ne.sum[0], ne.sum[1], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
		//logs.Log.App(" *                            —— %s合计采集【数据 %v 条 + 文件 %v 个】，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
		//	prefix, ne.sum[0], ne.sum[1], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), ne.takeTime)
	}
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)
	fmt.Println(` *********************************************************************************************************************************** `)

	// 单机模式并发运行，需要标记任务结束
	if ne.AppConf.Mode == status.OFFLINE {
		//ne.LogRest()
		ne.finishOnce.Do(func() { close(ne.finish) })
	}
}

// 任务同步执行
func (ne *NodeEntity) goSyncRun() {
	// 执行任务
	//pause:
	//	if ne.IsPause() {
	//		time.Sleep(time.Second)
	//		goto pause
	//	}
	defer func() {
		fmt.Println("goSyncRun end ^^^^^^^^^^^^^")
	}()
	for {
		b := ne.DataBoxQueue.GetOneActive()
		fmt.Println("GetOneActive%%%%%%%%%%%%%%%%%%%%%")
		// 从数据信使队列取出空闲信使，并发执行
		go ne.goManRunBox(b)
	}
}

func (ne *NodeEntity) goManRunBox(b *databox.DataBox) {
	m := ne.DataManPool.Use()
	if m != nil {
		// 执行并返回结果消息
		m.Init(b).SyncRun()

		// 任务结束后回收该信使
		ne.RWMutex.RLock()
		if ne.status != status.STOP {
			ne.DataManPool.Free(m)
		}
		ne.RWMutex.RUnlock()
	}
}

// 执行ActiveBox请求，同步返回
func (ne *NodeEntity) RunActiveBox(b *databox.DataBox, obj interface{}) response.DataResponse {
	var context *databox.Context
	var dataResp response.DataResponse

	m := ne.DataManPool.Use()
	if m != nil {
		// 执行并返回结果消息
		ne.RWMutex.RLock()
		context = m.MiniInit(b).RunRequest(obj)

		dataResp = *context.DataResponse
		defer databox.PutContext(context)

		// 该条请求文件结果存入pipeline
		om := ne.DataManPool.GetOneById(b.OrigDataManId)
		//fmt.Println("dataman write to file, manId: ", om.GetId())
		//fmt.Println("dataman write content to file: ", context)
		for _, f := range context.PullFiles() {
			if om.GetPipeline().CollectFile(f) != nil {
				break
			}
		}
		// 该条请求文本结果存入pipeline
		for _, item := range context.PullItems() {
			if om.GetPipeline().CollectData(item) != nil {
				break
			}
		}

		// 任务结束后回收该信使
		if ne.status != status.STOP {
			m.Stop() // 停止信使
			ne.DataManPool.Free(m)
		}
		ne.RWMutex.RUnlock()
	}
	return dataResp
}

func (ne *NodeEntity) StopActiveBox(b *databox.DataBox) {
	go func() {
		b.ActiveWG.Wait()

		ne.RWMutex.RLock()
		defer ne.RWMutex.RUnlock()

		close(b.BlockChan)
		b.RemoveActiveDataBox()

		//count := 0
		//for {
		//	if b.IsRequestEmpty() {
		//		if count == 0 {
		//			time.Sleep(5 * time.Second)
		//			continue
		//		}
		//
		//		ne.RWMutex.RLock()
		//		defer ne.RWMutex.RUnlock()
		//
		//		close(b.BlockChan)
		//		b.RemoveActiveDataBox()
		//		break
		//	}
		//}
	}()

}

// Offline 模式下暂停\恢复任务
func (self *NodeEntity) PauseRecover() {
	switch self.Status() {
	case status.PAUSE:
		self.setStatus(status.RUN)
	case status.RUN:
		self.setStatus(status.PAUSE)
	}

	scheduler.PauseRecover()
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

// 服务器模式运行，必须在DataBoxPrepare()执行之后调用才可以成功添加任务
// 生成的任务与自身当前全局配置相同
func (self *NodeEntity) server() {
	// 标记结束
	defer func() {
		self.finishOnce.Do(func() { close(self.finish) })
	}()

	// 便利添加任务到库
	tasksNum, dataBoxsNum := self.addNewTask()

	if tasksNum == 0 {
		return
	}

	// 打印报告
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(" *                               —— 本次成功添加 %v 条任务，共包含 %v 条采集规则 ——", tasksNum, dataBoxsNum)
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)
	fmt.Println(" * ")
	fmt.Println(` *********************************************************************************************************************************** `)
	fmt.Println(" * ")
	fmt.Println(" *                               —— 本次成功添加 %v 条任务，共包含 %v 条采集规则 ——", tasksNum, dataBoxsNum)
	fmt.Println(" * ")
	fmt.Println(` *********************************************************************************************************************************** `)
}

// 服务器模式下，生成task并添加至库
func (self *NodeEntity) addNewTask() (tasksNum, dataBoxsNum int) {
	length := self.DataBoxQueue.Len()
	t := distribute.Task{}
	// 从配置读取字段
	self.setTask(&t)

	for i, sp := range self.DataBoxQueue.GetAll() {

		t.DataBoxs = append(t.DataBoxs, map[string]string{"name": sp.GetName(), "keyin": sp.GetKeyin()})
		dataBoxsNum++

		// 每十个databox存为一个任务
		if i > 0 && i%10 == 0 && length > 10 {
			// 存入
			one := t
			self.TaskBase.Push(&one)
			// logs.Log.App(" *     [新增任务]   详情： %#v", *t)

			tasksNum++

			// 清空databox
			t.DataBoxs = []map[string]string{}
		}
	}

	if len(t.DataBoxs) != 0 {
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
	self.DataBoxQueue.Reset()

	// 更改全局配置
	self.setAppConf(t)

	// 初始化databox队列
	for _, n := range t.DataBoxs {
		df := self.GetDataBoxByName(n["name"])
		if df == nil {
			continue
		}
		dfcopy := df.Copy()
		dfcopy.SetPausetime(t.Pausetime)
		if dfcopy.GetLimit() > 0 {
			dfcopy.SetLimit(t.Limit)
		} else {
			dfcopy.SetLimit(-1 * t.Limit)
		}
		if v, ok := n["keyin"]; ok {
			dfcopy.SetKeyin(v)
		}
		self.DataBoxQueue.Add(dfcopy)
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
