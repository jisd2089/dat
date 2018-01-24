package dataman

/**
    Author: luzequan
    Created: 2017-12-27 13:35:22
*/
import (
	"runtime"
	"bytes"
	"math/rand"

	"dat/core/databox"
	"dat/core/interaction/request"
	"dat/core/interaction"
	"dat/runtime/cache"
	"dat/core/pipeline"
	"time"


	"fmt"
	"dat/runtime/status"
)

// 数据信使配送引擎
type (
	DataMan interface {
		MiniInit(box *databox.DataBox) DataMan       // 最小化初始化配送信使引擎
		Init(box *databox.DataBox) DataMan           // 初始化配送信使引擎
		Run()                                        // 运行配送任务
		SyncRun()                                    // 同步运行配送任务
		RunRequest(obj interface{}) *databox.Context // 执行DataRequest
		Stop()                                       // 主动终止
		CanStop() bool                               // 能否终止
		GetId() int                                  // 获取引擎ID
		GetPipeline() pipeline.Pipeline              // 获取管道
	}
	dataMan struct {
		*databox.DataBox    //执行的采集规则
		interaction.Carrier //全局公用的信息交互载体
		pipeline.Pipeline   // 拆包与核验管道
		id     int          // 信使ID
		pause  [2]int64     //[请求间隔的最短时长,请求间隔的增幅时长]
		status int          // 信使状态
	}
)

func New(id int) DataMan {
	return &dataMan{
		id:      id,
		status:  status.RUN,
		Carrier: interaction.CrossHandler,
	}
}

func (m *dataMan) Init(f *databox.DataBox) DataMan {
	m.DataBox = f.ReqmatrixInit()
	m.Pipeline = pipeline.New(f)
	m.pause[0] = f.Pausetime / 2
	if m.pause[0] > 0 {
		m.pause[1] = m.pause[0] * 3
	} else {
		m.pause[1] = 1
	}
	return m
}

func (m *dataMan) MiniInit(b *databox.DataBox) DataMan {
	m.DataBox = b
	return m
}

// 任务执行入口
func (m *dataMan) Run() {
	// 预先启动数据拆包/核验管道
	m.Pipeline.Start()

	// 运行处理协程
	c := make(chan bool)
	go func() {
		m.run()
		close(c)
	}()

	// 启动任务
	m.DataBox.Start()

	<-c // 等待处理协程退出

	// 停止数据拆包/核验管道
	m.Pipeline.Stop()
}

// 任务同步执行入口
func (m *dataMan) SyncRun() {
	// 预先启动数据拆包/核验管道
	m.Pipeline.Start()

	m.DataBox.BlockChan = make(chan bool)

	// 启动任务
	m.DataBox.Start()

	// 持续活跃DataBox的原始dataman id
	m.DataBox.OrigDataManId = m.id

	// 启动成功后加入活跃队列
	m.DataBox.AddActiveList()

	m.DataBox.StartWG.Done()

	<-m.DataBox.BlockChan // 等待处理协程退出

	//dataResp := m.syncRun()

	// 停止数据拆包/核验管道
	m.Pipeline.Stop()
	fmt.Println("SyncRun stop ...")
}

func (m *dataMan) run() {
	for {
		// 队列中取出一条请求并处理
		req := m.GetOne()
		if req == nil {
			// 停止任务
			time.Sleep(10 * time.Millisecond)
			if m.DataBox.CanStop() {
				fmt.Println("DataBox Can Stop@@@@@@@@@@@@@@@@@@@@@@@@@@")
				break
			}
			continue
		}

		// 执行请求
		m.UseOne()
		go func() {
			defer func() {
				m.FreeOne()
			}()
			//logs.Log.Debug(" *     Start: %v", req.GetUrl())
			m.Process(req)
		}()

		// 随机等待
		m.sleep()
	}

	// 等待处理中的任务完成
	m.DataBox.Defer()
}

func (m *dataMan) RunRequest(obj interface{}) *databox.Context {
	// 队列中取出一条请求并处理
	var req *request.DataRequest
	for {
		req = m.GetOne()
		if req == nil {
			m.DataBox.Start()
			//// 停止任务
			//if m.DataBox.CanStop() {
			//	return nil
			//}
			//return nil
		} else {
			break
		}
	}

	// 执行请求
	m.UseOne()
	defer func() {
		m.FreeOne()
	}()
	//logs.Log.Debug(" *     Start: %v", req.GetUrl())
	req.Bobject = obj
	context := m.SyncProcess(req)

	// 等待处理中的任务完成
	//m.DataBox.Defer()

	return context
}

func (m *dataMan) CanStop() bool {
	if m.status == status.STOP {
		return true
	}
	return false
}

// 临时dataman终止
func (m *dataMan) Stop() {
	// 主动崩溃DataBox运行协程
	m.status = status.STOP

	//m.Carrier.Close()

	//m.DataBox.Stop()
	//m.Pipeline.Stop()
}

// core processer
func (m *dataMan) Process(req *request.DataRequest) {
	var (
		//downUrl = req.GetUrl()
		b = m.DataBox
	)
	defer func() {
		if p := recover(); p != nil {
			if b.IsStopping() {
				// println("Process$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
				return
			}
			// 返回是否作为新的失败请求被添加至队列尾部
			if b.DoHistory(req, false) {
				// 统计失败数
				cache.PageFailCount()
			}
			// 提示错误
			stack := make([]byte, 4<<10) //4KB
			length := runtime.Stack(stack, true)
			start := bytes.Index(stack, []byte("/src/runtime/panic.go"))
			stack = stack[start:length]
			start = bytes.Index(stack, []byte("\n")) + 1
			stack = stack[start:]
			if end := bytes.Index(stack, []byte("\ngoroutine ")); end != -1 {
				stack = stack[:end]
			}
			stack = bytes.Replace(stack, []byte("\n"), []byte("\r\n"), -1)
			//logs.Log.Error(" *     Panic  [process][%s]: %s\r\n[TRACE]\r\n%s", downUrl, p, stack)
		}
	}()

	// TODO execute http、kafka、protocolbuffer... communication
	var ctx = m.Carrier.Handle(b, req)
	//var ctx = self.Downloader.Download(sp, req) // download page

	if err := ctx.GetError(); err != nil {
		// 返回是否作为新的失败请求被添加至队列尾部
		if b.DoHistory(req, false) {
			// 统计失败数
			cache.PageFailCount()
		}
		// 提示错误
		//logs.Log.Error(" *     Fail  [download][%v]: %v\n", downUrl, err)
		return
	}

	// 过程处理，提炼数据
	ctx.Parse(req.GetRuleName())

	// 该条请求文件结果存入pipeline
	for _, f := range ctx.PullFiles() {
		if m.Pipeline.CollectFile(f) != nil {
			break
		}
	}
	// 该条请求文本结果存入pipeline
	for _, item := range ctx.PullItems() {
		if m.Pipeline.CollectData(item) != nil {
			break
		}
	}

	// 处理成功请求记录
	b.DoHistory(req, true)

	// 统计成功页数
	cache.PageSuccCount()

	// 提示抓取成功
	//logs.Log.Informational(" *     Success: %v\n", downUrl)

	// 释放ctx准备复用
	databox.PutContext(ctx)
}

// core processer
func (m *dataMan) SyncProcess(req *request.DataRequest) *databox.Context {
	var (
		//downUrl = req.GetUrl()
		b = m.DataBox
	)
	defer func() {
		if p := recover(); p != nil {
			if b.IsStopping() {
				// println("Process$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
				return
			}
			// 返回是否作为新的失败请求被添加至队列尾部
			if b.DoHistory(req, false) {
				// 统计失败数
				cache.PageFailCount()
			}
			// 提示错误
			stack := make([]byte, 4<<10) //4KB
			length := runtime.Stack(stack, true)
			start := bytes.Index(stack, []byte("/src/runtime/panic.go"))
			stack = stack[start:length]
			start = bytes.Index(stack, []byte("\n")) + 1
			stack = stack[start:]
			if end := bytes.Index(stack, []byte("\ngoroutine ")); end != -1 {
				stack = stack[:end]
			}
			stack = bytes.Replace(stack, []byte("\n"), []byte("\r\n"), -1)
			//logs.Log.Error(" *     Panic  [process][%s]: %s\r\n[TRACE]\r\n%s", downUrl, p, stack)
		}
	}()

	// TODO execute http、kafka、protocolbuffer... communication
	var ctx = m.Carrier.Handle(b, req)
	//var ctx = self.Downloader.Download(sp, req) // download page

	if err := ctx.GetError(); err != nil {
		// 返回是否作为新的失败请求被添加至队列尾部
		if b.DoHistory(req, false) {
			// 统计失败数
			cache.PageFailCount()
		}
		// 提示错误
		//logs.Log.Error(" *     Fail  [download][%v]: %v\n", downUrl, err)
		return nil
	}

	// 过程处理，提炼数据
	ctx.Parse(req.GetRuleName())

	// 该条请求文件结果存入pipeline
	//for _, f := range ctx.PullFiles() {
	//	if om.GetPipeline().CollectFile(f) != nil {
	//		break
	//	}
	//}
	//// 该条请求文本结果存入pipeline
	//for _, item := range ctx.PullItems() {
	//	if om.GetPipeline().CollectData(item) != nil {
	//		break
	//	}
	//}

	// 处理成功请求记录
	b.DoHistory(req, true)

	// 统计成功页数
	cache.PageSuccCount()

	// 提示抓取成功
	//logs.Log.Informational(" *     Success: %v\n", downUrl)

	// 释放ctx准备复用
	//defer databox.PutContext(ctx)
	return ctx
}

// 从调度读取一个请求
func (m *dataMan) GetOne() *request.DataRequest {
	return m.DataBox.RequestPull()
}

func (m *dataMan) IsRequestEmpty() bool {
	return m.DataBox.IsRequestEmpty()
}

//从调度使用一个资源空位
func (m *dataMan) UseOne() {
	m.DataBox.RequestUse()
}

//从调度释放一个资源空位
func (m *dataMan) FreeOne() {
	m.DataBox.RequestFree()
}

// 常用基础方法
func (m *dataMan) sleep() {
	sleeptime := m.pause[0] + rand.Int63n(m.pause[1])
	time.Sleep(time.Duration(sleeptime) * time.Millisecond)
}

func (m *dataMan) SetId(id int) {
	m.id = id
}

func (m *dataMan) GetId() int {
	return m.id
}

func (m *dataMan) GetPipeline() pipeline.Pipeline {
	return m.Pipeline
}
