package common

/**
    Author: luzequan
    Created: 2018-08-03 09:40:43
*/


type CommonRequest struct {
	DemMemId string
	SupMemId string

	PreTimeOut      int
	AuthMode        string
	Reqsign         string
	JobId           string
	SerialNo        string
	IdType          string
	DmpSerialNo     string
	BusiSerialNo    string
	BusiReqInfo     map[string]interface{}
	BusiInfoStr     string
	BusiReqInfoHash string
	TimeStamp       string
	TaskIdStr       string
	Timeout         int
	//TaskIdInfo      []xml.TaskIdInfo
	ConnObjStr   string
	SvcStartTime int
	TranHash     string
	SignBlock1   string
	ServiceId    string
	UnitPrice    float64
}

type CommonResponse struct {
	SuccFlag bool
	ResCode  string
	ResMsg   string
	ErrCode  string
	ErrMsg   string
	//FlowStatus   busilog.FlowStatus
	SupMemId     string
	TaskId       string
	BusiSerialNo string
	DmpSerialNo  string
	TimeStamp    string
	PubResInfo   map[string]interface{}
	BusiResInfo  map[string]interface{}
	BusiInfoStr  string
	SvcStartTime int
	SignBlock1   string
	SignMemId2   string
	SignBlock2   string
	Chargflag    bool
}

type CommonRequestData struct {
	PubReqInfo PubReqInfo             `json:"pubReqInfo"`
	BusiInfo   map[string]interface{} `json:"busiInfo"`
}

type PubReqInfo struct {
	MemId        string `json:"memId"`
	SerialNo     string `json:"serialNo"`
	JobId        string `json:"jobId"`
	AuthMode     string `json:"authMode"`
	TimeStamp    string `json:"timeStamp"`
	ReqSign      string `json:"reqSign"`
}

type PubAnsInfo struct {
	SerialNo     string `json:"serialNo"`
	BusiSerialNo string `json:"busiSerialNo"`
	ResCode      string `json:"resCode"`
	ResMsg       string `json:"resMsg"`
	TimeStamp    string `json:"timeStamp"`
}

type ResponseInfo struct {
	PubResInfo  *PubResInfo            `json:"PubResInfo"`
	BusiResInfo map[string]interface{} `json:"BusiResInfo"`
}

type PubResInfo struct {
	ResCode    string `json:"resCode"`
	ResMsg     string `json:"resMsg"`
	Chargeflag string `json:"chargflag"`
}

/*
业务日志相关数据结构和接口
*/

// RecordType 记录类型
type RecordType string

const (
	// RecordTypeSingle 单笔
	RecordTypeSingle RecordType = "1"
	// RecordTypeBatch 批量
	RecordTypeBatch RecordType = "2"
)

// FlowStatus 流程状态
type FlowStatus string

const (
	// FlowStatusDemReq 需方请求
	FlowStatusDemReq FlowStatus = "01"
	// FlowStatusSupSucc 供方成功
	FlowStatusSupSucc FlowStatus = "11"
	// FlowStatusSupFail 供方失败
	FlowStatusSupFail FlowStatus = "12"
	// FlowStatusDemSucc 需方服务查询成功
	FlowStatusDemSucc FlowStatus = "21"
	// FlowStatusDemFail 需方服务查询失败
	FlowStatusDemFail FlowStatus = "22"
	// FlowStatusDemCache 需方缓存查询成功
	GetCacheSucc FlowStatus = "24"
	// 数据已缓存
	DataCached FlowStatus = "23"
)

const (
	// ErrCodeSucc 成功时的错误码
	ErrCodeSucc string = "000000"

	//需方返回代码
	//解析请求失败
	ErrorUnmarshalReq = "021001"
	//秘钥初始化参数错误
	ErrorEncryptionInit = "021002"
	//请求供方服务失败
	ErrorPostSup = "021003"
	//请求供方服务返回数据无法序列化
	ErrorUnmarshalPb = "021004"
	//需方处理错误
	ErrorPanic = "021999"

	// no data 查询结果为空
	ErrorNoData = "030002"
	//获取传输任务的订单信息错误
	ErrorOrderInfo = "030001"
	//获取供方服务URL错误
	ErrorSupUrl = "031002"
	// 异步调用 直接返回code
	ErrorAsync = "031006"
	//数据已经缓存
	ErrorDataCached = "001000"
	//余额申请错误
	ErrorBalance = "042000"
	//Authentication error
	ErrorDMPAuthentication = "042001"
	//EXID error
	ErrorEXID = "042002"
	//supplier error
	SUP_AUTH_ERROR = "042004"
)

// 4.2 中心返回代码列表
/**
999999		其他未定义错误，具体信息参见返回信息
000000		成功
020001		报文格式错误
020002		服务不存在
021000		未上送appkey
021001		appkey校验失败
021002		公钥不存在
021003		签名校验失败
032001		请求数据异常
032999		请求失败
*/
const (
	CenterCodeOther      = "999999" //其他未定义错误，具体信息参见返回信息
	CenterCodeSucc       = "000000" //成功
	CenterCodeFormat     = "020001" //报文格式错误
	CenterCodeNoService  = "020002" //服务不存在
	CenterCodeNoAppkey   = "021000" //未上送appkey
	CenterCodeFailAppkey = "021001" //appkey校验失败
	CenterCodeNoPubkey   = "021002" //公钥不存在
	CenterCodeFailSign   = "021003" //签名校验失败
	CenterCodeReqData    = "032001" //请求数据异常
	CenterCodeReqFail    = "032999" //请求失败
	CenterCodeNoAccess   = "033000" //无接口访问权限
	CenterCodeNoMoney    = "033001" //无接口访问权限,欠费

	CenterCodeTestNoHit  = "030002" // 未查找到该条数据
)

var centerCodeText = map[string]string{
	CenterCodeOther:      "其他未定义错误，具体信息参见返回信息",
	CenterCodeSucc:       "成功",
	CenterCodeFormat:     "报文格式错误",
	CenterCodeNoService:  "服务不存在",
	CenterCodeNoAppkey:   "未上送appkey",
	CenterCodeFailAppkey: "appkey校验失败",
	CenterCodeNoPubkey:   "公钥不存在",
	CenterCodeFailSign:   "签名校验失败",
	CenterCodeReqData:    "请求数据异常",
	CenterCodeReqFail:    "请求失败",
	CenterCodeNoAccess:   "无接口访问权限",
	CenterCodeNoMoney:    "无接口访问权限,欠费",

	CenterCodeTestNoHit:  "未查找到该条数据",
}

//根据中心code返回text内容
func GetCenterCodeText(code string) string {
	str, ok := centerCodeText[code]
	if ok {
		return str
	}
	return centerCodeText[CenterCodeOther]
}