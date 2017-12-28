package dataman

/**
    Author: luzequan
    Created: 2017-12-27 13:35:22
*/
import (
	"runtime"
	"bytes"

	"dat/core/dataflow"
	"dat/core/interaction/request"
	"dat/core/interaction"
	"dat/runtime/cache"
	"dat/core/pipeline"
)

// 数据信使配送引擎
type (
	DataMan interface {
		Init(flow *dataflow.DataFlow) DataMan //初始化配送引擎
		Run()                                 //运行配送任务
	}
	dataMan struct {
		*dataflow.DataFlow  //执行的采集规则
		interaction.Carrier //全局公用的信息交互载体
		pipeline.Pipeline
		id int              // 信使ID
	}
)

func New(id int) DataMan {
	return &dataMan{
		id: id,
	}
}

func (m *dataMan) Init(f *dataflow.DataFlow) DataMan {
	return m
}

// 任务执行入口
func (m *dataMan) Run() {
	// 预先启动数据收集/输出管道
	//m.Pipeline.Start()

	// 运行处理协程
	c := make(chan bool)
	go func() {
		m.run()
		close(c)
	}()

	// 启动任务
	m.DataFlow.Start()

	<-c // 等待处理协程退出

	// 停止数据收集/输出管道
	//m.Pipeline.Stop()
}

func (m *dataMan) run() {
	//for {
	//	// 队列中取出一条请求并处理
	//	req := self.GetOne()
	//	if req == nil {
	//		// 停止任务
	//		if self.Spider.CanStop() {
	//			break
	//		}
	//		time.Sleep(20 * time.Millisecond)
	//		continue
	//	}
	//
	//	// 执行请求
	//	self.UseOne()
	//	go func() {
	//		defer func() {
	//			self.FreeOne()
	//		}()
	//		logs.Log.Debug(" *     Start: %v", req.GetUrl())
	//		self.Process(req)
	//	}()
	//
	//	// 随机等待
	//	self.sleep()
	//}
	//
	//// 等待处理中的任务完成
	//self.Spider.Defer()
}

// core processer
func (m *dataMan) Process(req *request.Request) {
	var (
		downUrl = req.GetUrl()
		df      = m.DataFlow
	)
	defer func() {
		if p := recover(); p != nil {
			if df.IsStopping() {
				// println("Process$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
				return
			}
			// 返回是否作为新的失败请求被添加至队列尾部
			if df.DoHistory(req, false) {
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

	// TODO execute http、kafka... communication
	var ctx = m.Carrier.Handle(df, req)
	//var ctx = self.Downloader.Download(sp, req) // download page

	if err := ctx.GetError(); err != nil {
		// 返回是否作为新的失败请求被添加至队列尾部
		if df.DoHistory(req, false) {
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
	df.DoHistory(req, true)

	// 统计成功页数
	cache.PageSuccCount()

	// 提示抓取成功
	//logs.Log.Informational(" *     Success: %v\n", downUrl)

	// 释放ctx准备复用
	dataflow.PutContext(ctx)
}
