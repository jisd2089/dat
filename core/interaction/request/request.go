package request

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"drcs/common/util"
	"drcs/common/sftp"
	"mime/multipart"
	"gopkg.in/redis.v5"
)

// DataRequest represents object waiting for being crawled.
type DataRequest struct {
	DataBox       string                // 规则名，自动设置，禁止人为填写
	DataBoxId     int                   // databox id
	TransferType  string                // 传输类型
	Url           string                // 目标URL，必须设置
	Rule          string                // 用于解析响应的规则节点名，必须设置
	Method        string                // GET POST POST-M HEAD
	Header        http.Header           // 请求头信息
	EnableCookie  bool                  // 是否使用cookies，在DataBox的EnableCookie设置
	Parameters    []byte                // 传参
	CommandName   string                // command命令名称
	CommandParams []string              // command参数
	Bobject       interface{}           // 业务参数
	PostData      string                // POST values
	DialTimeout   time.Duration         // 创建连接超时 dial tcp: i/o timeout
	ConnTimeout   time.Duration         // 连接状态超时 WSARecv tcp: i/o timeout
	TryTimes      int                   // 尝试下载的最大次数
	RetryPause    time.Duration         // 下载失败后，下次尝试下载的等待时间
	RedirectTimes int                   // 重定向的最大次数，为0时不限，小于0时禁止重定向
	Temp          Temp                  // 临时数据
	TempIsJson    map[string]bool       // 将Temp中以JSON存储的字段标记为true，自动设置，禁止人为填写
	Priority      int                   // 指定调度优先级，默认为0（最小优先级为0）
	Reloadable    bool                  // 是否允许重复该链接下载
	FileCatalog   *sftp.FileCatalog     // SFTP使用
	RedisOptions  *redis.Options        // Redis 连接参数
	DataFile      *multipart.FileHeader // http传输文件
	//Surfer下载器内核ID
	//0为Surf高并发下载器，各种控制功能齐全
	//1为PhantomJS下载器，特点破防力强，速度慢，低并发
	DownloaderID int
	TimeOutCh    chan string // 超时channel

	proxy  string //当用户界面设置可使用代理IP时，自动设置代理
	unique string //ID
	lock   sync.RWMutex
}

const (
	DefaultDialTimeout = 2 * time.Minute // 默认请求服务器超时
	DefaultConnTimeout = 2 * time.Minute // 默认下载超时
	DefaultTryTimes    = 3               // 默认最大下载次数
	DefaultRetryPause  = 2 * time.Second // 默认重新下载前停顿时长
)

const (
	SURF_ID    = 0 // 默认的surf下载内核（Go原生），此值不可改动
	PHANTOM_ID = 1 // 备用的phantomjs下载内核，一般不使用（效率差，头信息支持不完善）

	HTTP      = "HTTP"
	FASTHTTP  = "FASTHTTP"
	SFTP      = "SFTP"
	REDIS     = "REDIS"
	DATABOX   = "DATABOX"
	NONETYPE  = "NONETYPE"
	FILETYPE  = "FILETYPE"
	SHELLTYPE = "SHELLTYPE"
)

// 发送请求前的准备工作，设置一系列默认值
// DataRequest.Url与Request.Rule必须设置
// DataRequest.DataBox无需手动设置(由系统自动设置)
// DataRequest.EnableCookie在DataBox字段中统一设置，规则请求中指定的无效
// 以下字段有默认值，可不设置:
// DataRequest.Method默认为GET方法;
// DataRequest.DialTimeout默认为常量DefaultDialTimeout，小于0时不限制等待响应时长;
// DataRequest.ConnTimeout默认为常量DefaultConnTimeout，小于0时不限制下载超时;
// DataRequest.TryTimes默认为常量DefaultTryTimes，小于0时不限制失败重载次数;
// DataRequest.RedirectTimes默认不限制重定向次数，小于0时可禁止重定向跳转;
// DataRequest.RetryPause默认为常量DefaultRetryPause;
// DataRequest.DownloaderID指定下载器ID，0为默认的Surf高并发下载器，功能完备，1为PhantomJS下载器，特点破防力强，速度慢，低并发。
func (self *DataRequest) Prepare() error {
	// 确保url正确，且和Response中Url字符串相等
	URL, err := url.Parse(self.Url)
	if err != nil {
		return err
	} else {
		self.Url = URL.String()
	}

	if self.Method == "" {
		self.Method = "GET"
	} else {
		self.Method = strings.ToUpper(self.Method)
	}

	if self.Header == nil {
		self.Header = make(http.Header)
	}

	if self.DialTimeout < 0 {
		self.DialTimeout = 0
	} else if self.DialTimeout == 0 {
		self.DialTimeout = DefaultDialTimeout
	}

	if self.ConnTimeout < 0 {
		self.ConnTimeout = 0
	} else if self.ConnTimeout == 0 {
		self.ConnTimeout = DefaultConnTimeout
	}

	if self.TryTimes == 0 {
		self.TryTimes = DefaultTryTimes
	}

	if self.RetryPause <= 0 {
		self.RetryPause = DefaultRetryPause
	}

	if self.Priority < 0 {
		self.Priority = 0
	}

	// TODO
	if self.DownloaderID < SURF_ID || self.DownloaderID > PHANTOM_ID {
		self.DownloaderID = SURF_ID
	}

	if self.TempIsJson == nil {
		self.TempIsJson = make(map[string]bool)
	}

	if self.Temp == nil {
		self.Temp = make(Temp)
	}
	return nil
}

// 反序列化
func UnSerialize(s string) (*DataRequest, error) {
	req := &DataRequest{}
	return req, json.Unmarshal([]byte(s), req)
}

// 序列化
func (self *DataRequest) Serialize() string {
	for k, v := range self.Temp {
		self.Temp.set(k, v)
		self.TempIsJson[k] = true
	}
	b, _ := json.Marshal(self)
	return strings.Replace(util.Bytes2String(b), `\u0026`, `&`, -1)
}

// 请求的唯一识别码
func (self *DataRequest) Unique() string {
	if self.unique == "" {
		block := md5.Sum([]byte(self.DataBox + self.Rule + self.Url + self.Method ))
		self.unique = hex.EncodeToString(block[:])
	}
	return self.unique
}

// 获取副本
func (self *DataRequest) Copy() *DataRequest {
	reqcopy := &DataRequest{}
	b, _ := json.Marshal(self)
	json.Unmarshal(b, reqcopy)
	return reqcopy
}

// 获取Url
func (self *DataRequest) GetUrl() string {
	return self.Url
}

// 获取Http请求的方法名称 (注意这里不是指Http GET方法)
func (self *DataRequest) GetMethod() string {
	return self.Method
}

// 设定Http请求方法的类型
func (self *DataRequest) SetMethod(method string) *DataRequest {
	self.Method = strings.ToUpper(method)
	return self
}

func (self *DataRequest) SetUrl(url string) *DataRequest {
	self.Url = url
	return self
}

func (self *DataRequest) GetReferer() string {
	return self.Header.Get("Referer")
}

func (self *DataRequest) SetReferer(referer string) *DataRequest {
	self.Header.Set("Referer", referer)
	return self
}

func (self *DataRequest) GetBobject() interface{} {
	return self.Bobject
}

func (self *DataRequest) GetPostData() string {
	return self.PostData
}

func (self *DataRequest) GetHeader() http.Header {
	return self.Header
}

func (self *DataRequest) SetHeader(key, value string) *DataRequest {
	self.Header.Set(key, value)
	return self
}

func (self *DataRequest) AddHeader(key, value string) *DataRequest {
	self.Header.Add(key, value)
	return self
}

func (self *DataRequest) GetEnableCookie() bool {
	return self.EnableCookie
}

func (self *DataRequest) SetEnableCookie(enableCookie bool) *DataRequest {
	self.EnableCookie = enableCookie
	return self
}

func (self *DataRequest) GetCookies() string {
	return self.Header.Get("Cookie")
}

func (self *DataRequest) SetCookies(cookie string) *DataRequest {
	self.Header.Set("Cookie", cookie)
	return self
}

func (self *DataRequest) GetDialTimeout() time.Duration {
	return self.DialTimeout
}

func (self *DataRequest) GetConnTimeout() time.Duration {
	return self.ConnTimeout
}

func (self *DataRequest) GetTryTimes() int {
	return self.TryTimes
}

func (self *DataRequest) GetRetryPause() time.Duration {
	return self.RetryPause
}

func (self *DataRequest) GetProxy() string {
	return self.proxy
}

func (self *DataRequest) SetProxy(proxy string) *DataRequest {
	self.proxy = proxy
	return self
}

func (self *DataRequest) GetRedirectTimes() int {
	return self.RedirectTimes
}

func (self *DataRequest) GetRuleName() string {
	return self.Rule
}

func (self *DataRequest) SetRuleName(ruleName string) *DataRequest {
	self.Rule = ruleName
	return self
}

func (self *DataRequest) GetDataBoxName() string {
	return self.DataBox
}

func (self *DataRequest) SetDataBoxName(dataBoxName string) *DataRequest {
	self.DataBox = dataBoxName
	return self
}

func (self *DataRequest) IsReloadable() bool {
	return self.Reloadable
}

func (self *DataRequest) SetReloadable(can bool) *DataRequest {
	self.Reloadable = can
	return self
}

// 获取临时缓存数据
// defaultValue 不能为 interface{}(nil)
func (self *DataRequest) GetTemp(key string, defaultValue interface{}) interface{} {
	if defaultValue == nil {
		panic("*DataRequest.GetTemp()的defaultValue不能为nil，错误位置：key=" + key)
	}
	self.lock.RLock()
	defer self.lock.RUnlock()

	if self.Temp[key] == nil {
		return defaultValue
	}

	if self.TempIsJson[key] {
		return self.Temp.get(key, defaultValue)
	}

	return self.Temp[key]
}

func (self *DataRequest) GetTemps() Temp {
	return self.Temp
}

func (self *DataRequest) SetTemp(key string, value interface{}) *DataRequest {
	self.lock.Lock()
	self.Temp[key] = value
	delete(self.TempIsJson, key)
	self.lock.Unlock()
	return self
}

func (self *DataRequest) SetTemps(temp map[string]interface{}) *DataRequest {
	self.lock.Lock()
	self.Temp = temp
	self.TempIsJson = make(map[string]bool)
	self.lock.Unlock()
	return self
}

func (self *DataRequest) GetPriority() int {
	return self.Priority
}

func (self *DataRequest) SetPriority(priority int) *DataRequest {
	self.Priority = priority
	return self
}

func (self *DataRequest) GetDownloaderID() int {
	return self.DownloaderID
}

func (self *DataRequest) SetDownloaderID(id int) *DataRequest {
	self.DownloaderID = id
	return self
}

func (self *DataRequest) GetTransferType() string {
	return self.TransferType
}

func (self *DataRequest) SetTransferType(transferType string) *DataRequest {
	self.TransferType = transferType
	return self
}

func (self *DataRequest) SetParameters(params []byte) *DataRequest {
	self.Parameters = params
	return self
}

func (self *DataRequest) GetParameters() []byte {
	return self.Parameters
}

func (dq *DataRequest) GetFileCatalog() *sftp.FileCatalog {
	return dq.FileCatalog
}

func (dq *DataRequest) GetDataFile() *multipart.FileHeader {
	return dq.DataFile
}

func (dq *DataRequest) GetRedisOptions() *redis.Options {
	return dq.RedisOptions
}

func (dq *DataRequest) GetCommandName() string {
	return dq.CommandName
}

func (dq *DataRequest) GetCommandParams() []string {
	return dq.CommandParams
}

func (self *DataRequest) MarshalJSON() ([]byte, error) {
	for k, v := range self.Temp {
		if self.TempIsJson[k] {
			continue
		}
		self.Temp.set(k, v)
		self.TempIsJson[k] = true
	}
	b, err := json.Marshal(*self)
	return b, err
}
