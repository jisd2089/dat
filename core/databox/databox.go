package databox

/**
    Author: luzequan
    Created: 2017-12-27 14:12:45
*/
import (
	"sync"
	"dat/runtime/status"
	"time"
	"math"
	"dat/common/util"
	"dat/core/scheduler"
	"dat/core/interaction/request"
	"dat/core/interaction/response"
	"mime/multipart"
)

const (
	KEYIN       = util.USE_KEYIN // 若使用DataBox.Keyin，则须在规则中设置初始值为USE_KEYIN
	LIMIT       = math.MaxInt64  // 如希望在规则中自定义控制Limit，则Limit初始值必须为LIMIT
	FORCED_STOP = "——主动终止DataBox——"
)

// 数据产品流
type (
	DataBox struct {
		// 以下字段由用户定义
		Name               string                                                      // 名称（应保证唯一性）
		Description        string                                                      // 描述
		DataFilePath       string                                                      // 数据文件地址
		DataFile           *multipart.FileHeader                                       // 数据文件内容
		NodeAddress        []*request.NodeAddress                                      // 交互节点地址
		Pausetime          int64                                                       // 随机暂停区间(50%~200%)，若规则中直接定义，则不被界面传参覆盖
		Limit              int64                                                       // 默认限制请求数，0为不限；若规则中定义为LIMIT，则采用规则的自定义限制方案
		Keyin              string                                                      // 自定义输入的配置信息，使用前须在规则中设置初始值为KEYIN
		EnableCookie       bool                                                        // 所有请求是否使用cookie记录
		NotDefaultField    bool                                                        // 是否禁止输出结果中的默认字段 Url/ParentUrl/DownloadTime
		Namespace          func(self *DataBox) string                                  // 命名空间，用于输出文件、路径的命名
		SubNamespace       func(self *DataBox, dataCell map[string]interface{}) string // 次级命名，用于输出文件、路径的命名，可依赖具体数据内容
		DetailCount        int                                                         // 明细条数
		TsfSuccCount       int                                                         // 流通成功明细条数
		BlockChan          chan bool                                                   // 用于ActiveDataBox阻塞，持续活跃
		StartWG            *sync.WaitGroup                                             // 启动成功通知
		RuleTree           *RuleTree                                                   // 定义具体的配送规则树
		OrigDataManId      int                                                         // 原始dataman id
		PairDataBoxId      int                                                         // 对接的databox id
		ActiveWG           *sync.WaitGroup                                             // 等待所有活动结束
		ChildBoxChan       chan *DataBox                                               // 子盒子通道
		ChildActiveBoxChan chan *DataBox                                               // 持续活跃子盒子通道
		IsParentBox        bool                                                        // 是否父databox
		ChildBox           *DataBox                                                    // child box
		ParentBox          *DataBox                                                    // parent box

		// 以下字段系统自动赋值
		id        int               // 自动分配的DataBoxQueue中的索引
		subName   string            // 由Keyin转换为的二级标识名
		reqMatrix *scheduler.Matrix // 请求矩阵
		timer     *Timer            // 定时器
		status    int               // 执行状态
		lock      sync.RWMutex
		once      sync.Once
	}
	// 数据产品流规则
	RuleTree struct {
		Root  func(*Context)   // 根节点(执行入口)
		Trunk map[string]*Rule // 节点散列表(业务规则过程)
	}
	Rule struct {
		ItemFields []string                                           // 结果字段列表(选填，写上可保证字段顺序) TODO 清点内容明细项
		ParseFunc  func(*Context)                                     // 内容解析函数
		AidFunc    func(*Context, map[string]interface{}) interface{} // 通用辅助函数
		FileFunc   func(*Context)                                     // 文件处理函数
		SyncFunc   func(*Context) *response.DataResponse              // 同步处理函数，带返回值
	}
)

/**
 *Action Function
 ****************************************************************************
 */

// 添加自身到数据流产品菜单
func (self DataBox) Register() *DataBox {
	self.status = status.STOPPED
	return Species.Add(&self)
}

// 根据名称获取child box
func (b *DataBox) GetChildBoxByName(name string) *DataBox {
	return Species.GetByName(name)
}

// 添加自身到活跃DataBox列表
func (self DataBox) AddActiveList() *DataBox {
	self.status = status.RUN
	return Activites.Add(&self)
}

// 从活跃DataBox列表移除
func (self DataBox) RemoveActiveDataBox() *DataBoxActivites {
	self.status = status.STOP
	return Activites.Remove(&self)
}

// 停止活跃databox
func (b *DataBox) StopActiveBox() {
	b.lock.RLock()
	defer b.lock.RUnlock()

	close(b.BlockChan)
	b.RemoveActiveDataBox()
}

// 数据流产品开始穿越
func (self *DataBox) Start() {
	defer func() {
		if p := recover(); p != nil {
			self.status = status.STOP
			//logs.Log.Error(" *     Panic  [root]: %v\n", p)
		}
		self.lock.Lock()
		self.status = status.RUN
		self.lock.Unlock()
	}()
	// 执行业务规则入口Function
	self.RuleTree.Root(GetContext(self, nil))
}

// 主动崩溃DataBox运行协程
func (self *DataBox) Stop() {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.status == status.STOP {
		return
	}
	self.status = status.STOP
	// 取消所有定时器
	if self.timer != nil {
		self.timer.drop()
		self.timer = nil
	}
}

func (b *DataBox) ExecTsfSuccCount() {
	b.lock.RLock()
	defer b.lock.RUnlock()
	b.TsfSuccCount ++
}

// 若已主动终止任务，则崩溃DataBox协程
func (self *DataBox) tryPanic() {
	if self.IsStopping() {
		panic(FORCED_STOP)
	}
}

// 退出任务前收尾工作
func (self *DataBox) Defer() {
	// 取消所有定时器
	if self.timer != nil {
		self.timer.drop()
		self.timer = nil
	}
	// 等待处理中的请求完成
	self.reqMatrix.Wait()
	// 更新失败记录
	self.reqMatrix.TryFlushFailure()
}

// 是否输出默认添加的字段 Url/ParentUrl/DownloadTime
func (self *DataBox) OutDefaultField() bool {
	return !self.NotDefaultField
}

/**
 * Private Function
 ***************************************************************************************
 */

// 获取流任务ID
func (self *DataBox) GetId() int {
	return self.id
}

// 设置流任务ID
func (self *DataBox) SetId(id int) {
	self.id = id
}

// 获取数据流产品名称
func (self *DataBox) GetName() string {
	return self.Name
}

// 指定规则的获取结果的字段名列表
func (self *DataBox) GetItemFields(rule *Rule) []string {
	return rule.ItemFields
}

// 返回结果字段名的值
// 不存在时返回空字符串
func (self *DataBox) GetItemField(rule *Rule, index int) (field string) {
	if index > len(rule.ItemFields)-1 || index < 0 {
		return ""
	}
	return rule.ItemFields[index]
}

// 返回结果字段名的其索引
// 不存在时索引为-1
func (self *DataBox) GetItemFieldIndex(rule *Rule, field string) (index int) {
	for idx, v := range rule.ItemFields {
		if v == field {
			return idx
		}
	}
	return -1
}

// 为指定Rule动态追加结果字段名，并返回索引位置
// 已存在时返回原来索引位置
func (self *DataBox) UpsertItemField(rule *Rule, field string) (index int) {
	for i, v := range rule.ItemFields {
		if v == field {
			return i
		}
	}
	rule.ItemFields = append(rule.ItemFields, field)
	return len(rule.ItemFields) - 1
}

// 获取databox二级标识名
func (self *DataBox) GetSubName() string {
	self.once.Do(func() {
		self.subName = self.GetKeyin()
		if len([]rune(self.subName)) > 8 {
			self.subName = util.MakeHash(self.subName)
		}
	})
	return self.subName
}

// 安全返回指定规则
func (self *DataBox) GetRule(ruleName string) (*Rule, bool) {
	rule, found := self.RuleTree.Trunk[ruleName]
	return rule, found
}

// 返回指定规则
func (self *DataBox) MustGetRule(ruleName string) *Rule {
	return self.RuleTree.Trunk[ruleName]
}

// 返回规则树
func (self *DataBox) GetRules() map[string]*Rule {
	return self.RuleTree.Trunk
}

// 获取DataBox描述
func (self *DataBox) GetDescription() string {
	return self.Description
}

// 获取自定义配置信息
func (self *DataBox) GetKeyin() string {
	return self.Keyin
}

// 设置自定义配置信息
func (self *DataBox) SetKeyin(keyword string) {
	self.Keyin = keyword
}

// 获取自定义DataFilePath配置信息
func (b *DataBox) GetDataFilePath() string {
	return b.DataFilePath
}

// 设置自定义DataFilePath配置信息
func (b *DataBox) SetDataFilePath(path string) {
	b.DataFilePath = path
}

// 获取自定义NodeAddress配置信息
func (b *DataBox) GetNodeAddress() []*request.NodeAddress {
	return b.NodeAddress
}

// 设置自定义DataFilePath配置信息
func (b *DataBox) SetNodeAddress(addrs []*request.NodeAddress) {
	b.NodeAddress = addrs
}

// 获取采集上限
// <0 表示采用限制请求数的方案
// >0 表示采用规则中的自定义限制方案
func (self *DataBox) GetLimit() int64 {
	return self.Limit
}

// 设置采集上限
// <0 表示采用限制请求数的方案
// >0 表示采用规则中的自定义限制方案
func (self *DataBox) SetLimit(max int64) {
	self.Limit = max
}

// 控制所有请求是否使用cookie
func (self *DataBox) GetEnableCookie() bool {
	return self.EnableCookie
}

// 自定义暂停时间 pause[0]~(pause[0]+pause[1])，优先级高于外部传参
// 当且仅当runtime[0]为true时可覆盖现有值
func (self *DataBox) SetPausetime(pause int64, runtime ...bool) {
	if self.Pausetime == 0 || len(runtime) > 0 && runtime[0] {
		self.Pausetime = pause
	}
}

// 设置定时器
// @id为定时器唯一标识
// @bell==nil时为倒计时器，此时@tol为睡眠时长
// @bell!=nil时为闹铃，此时@tol用于指定醒来时刻（从now起遇到的第tol个bell）
func (self *DataBox) SetTimer(id string, tol time.Duration, bell *Bell) bool {
	if self.timer == nil {
		self.timer = newTimer()
	}
	return self.timer.set(id, tol, bell)
}

// 启动定时器，并返回定时器是否可以继续使用
func (self *DataBox) RunTimer(id string) bool {
	if self.timer == nil {
		return false
	}
	return self.timer.sleep(id)
}

// 返回一个自身复制品
func (self *DataBox) Copy() *DataBox {
	ghost := &DataBox{}
	ghost.Name = self.Name
	ghost.subName = self.subName

	ghost.RuleTree = &RuleTree{
		Root:  self.RuleTree.Root,
		Trunk: make(map[string]*Rule, len(self.RuleTree.Trunk)),
	}
	for k, v := range self.RuleTree.Trunk {
		ghost.RuleTree.Trunk[k] = &Rule{}

		ghost.RuleTree.Trunk[k].ItemFields = make([]string, len(v.ItemFields))
		copy(ghost.RuleTree.Trunk[k].ItemFields, v.ItemFields)

		ghost.RuleTree.Trunk[k].ParseFunc = v.ParseFunc
		ghost.RuleTree.Trunk[k].AidFunc = v.AidFunc
	}

	ghost.Description = self.Description
	ghost.Pausetime = self.Pausetime
	ghost.EnableCookie = self.EnableCookie
	ghost.Limit = self.Limit
	ghost.Keyin = self.Keyin
	ghost.NodeAddress = self.NodeAddress
	ghost.DataFilePath = self.DataFilePath
	ghost.DataFile = self.DataFile

	ghost.NotDefaultField = self.NotDefaultField
	ghost.Namespace = self.Namespace
	ghost.SubNamespace = self.SubNamespace

	ghost.timer = self.timer
	ghost.status = self.status
	ghost.StartWG = self.StartWG
	ghost.PairDataBoxId = self.PairDataBoxId
	ghost.IsParentBox = self.IsParentBox
	ghost.ChildBoxChan = self.ChildBoxChan
	ghost.ChildActiveBoxChan = self.ChildActiveBoxChan

	return ghost
}

func (b *DataBox) Refresh() *DataBox {

	b.Description = ""
	b.Pausetime = 0
	b.EnableCookie = false
	b.Limit = 0
	b.Keyin = ""
	b.NodeAddress = nil
	b.DataFilePath = ""
	b.DataFile = nil

	b.NotDefaultField = false
	b.Namespace = nil
	b.SubNamespace = nil

	b.status = status.RUN
	b.StartWG = nil
	b.PairDataBoxId = 0
	b.IsParentBox = false
	return b
}

// DataRequest矩阵初始化
func (self *DataBox) ReqmatrixInit() *DataBox {
	if self.Limit < 0 {
		self.reqMatrix = scheduler.AddMatrix(self.GetName(), self.GetSubName(), self.Limit)
		self.SetLimit(0)
	} else {
		self.reqMatrix = scheduler.AddMatrix(self.GetName(), self.GetSubName(), math.MinInt64)
	}
	return self
}

// 返回是否作为新的失败请求被添加至队列尾部
func (self *DataBox) DoHistory(req *request.DataRequest, ok bool) bool {
	return self.reqMatrix.DoHistory(req, ok)
}

func (self *DataBox) RequestPush(req *request.DataRequest) {
	self.reqMatrix.Push(req)
}

func (self *DataBox) RequestPull() *request.DataRequest {
	r := self.reqMatrix.Pull()
	if r != nil {
		r.DataBoxId = self.GetId()
	}
	return r
}

func (self *DataBox) RequestPushChan(req *request.DataRequest) {
	self.reqMatrix.PushChan(req)
}

func (self *DataBox) RequestPullChan() *request.DataRequest {
	return self.reqMatrix.PullChan()
}

func (self *DataBox) RequestChan() chan *request.DataRequest {
	return self.reqMatrix.RequestChan()
}

func (self *DataBox) CloseRequestChan() {
	self.reqMatrix.CloseReqChan()
}

func (self *DataBox) IsRequestEmpty() bool {
	return self.reqMatrix.IsEmpty()
}

func (self *DataBox) AddressPush(addr *request.NodeAddress) {
	self.reqMatrix.PushAddr(addr)
}

func (self *DataBox) AddressPull() *request.NodeAddress {
	return self.reqMatrix.PullAddr()
}

func (self *DataBox) RequestUse() {
	self.reqMatrix.Use()
}

func (self *DataBox) RequestFree() {
	self.reqMatrix.Free()
}

func (self *DataBox) RequestLen() int {
	return self.reqMatrix.Len()
}

func (self *DataBox) TryFlushSuccess() {
	self.reqMatrix.TryFlushSuccess()
}

func (self *DataBox) TryFlushFailure() {
	self.reqMatrix.TryFlushFailure()
}

func (self *DataBox) CanStop() bool {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.status != status.STOPPED && self.reqMatrix.CanStop() && self.status != status.RUN
}

func (b *DataBox) SetStatus(status int) {
	b.status = status
}

func (self *DataBox) IsStopping() bool {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.status == status.STOP
}
