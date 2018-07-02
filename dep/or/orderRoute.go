package or

import (
	"drcs/dep/errors"
	logger "drcs/log"
	"drcs/dep/order"
	"drcs/settings"
	"drcs/runtime/cache"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

var (
	OrderRoutePolicyMap map[string]*OrderRoutePolicy
)

func init() {
	OrderRoutePolicyMap = make(map[string]*OrderRoutePolicy)
}

type Head struct {
	RouteMethod   RouteMethod    `xml:"routeMethod"`
	SvcTimeOut    int            `xml:"svcTimeOut"`
	PreTimeOut    int            `xml:"preTimeOut"`
	Batch         int            `xml:"batch"`
	MaxDelay      int            `xml:"maxDelay"`
	StopFlag      StopFlag       `xml:"stopFlag"`
	NeedCache     CacheFlag      `xml:"needCache"`
	NeedKafka     KafkaFlag      `xml:"needKafka"`
	CacheTime     int            `xml:"cacheTime"`
	SyncFlag      SyncFlag       `xml:"syncFlag"`
	CacheCodeList *CacheCodeList `xml:"cacheCodeList"`
}

type CacheCodeList struct {
	CacheCodes []*CacheCode `xml:"code"`
}

type CacheCode struct {
	Time string `xml:"time,attr"`
	Code string `xml:",chardata"`
}

type SvcInfo struct {
	MemId      string       `xml:"memId"`
	TaskIdList *TaskIdInfos `xml:"taskIdList"`
}

type TaskIdInfos struct {
	TaskIdInfo []*TaskIdInfo `xml:"taskId"`
}

type TaskIdInfo struct {
	KeyQryType string `xml:"keyQryType,attr"`
	TaskId     string `xml:",chardata"`
}

type SvcList struct {
	SvcInfo []*SvcInfo `xml:"svc_info"`
}

type SinglePolicy struct {
	MemId string `xml:"memId"`
}

type StaticPolicy struct {
	CallList *Calllist `xml:"call_list"`
}

type DynamicPolicy struct {
	PolicyType int `xml:"policyType"`
}

type BroadcastPolicy struct {
	CallList *Calllist `xml:"call_list"`
}

type Calllist struct {
	MemId []string `xml:"memId"`
}

type OrderRoute struct {
	XMLName         xml.Name         `xml:"route_info"`
	Head            *Head            `xml:"head"`
	SvcList         *SvcList         `xml:"svc_list"`
	SinglePolicy    *SinglePolicy    `xml:"single_policy"`
	StaticPolicy    *StaticPolicy    `xml:"static_policy"`
	DynamicPolicy   *DynamicPolicy   `xml:"dynamic_policy"`
	BroadcastPolicy *BroadcastPolicy `xml:"broadcast_policy"`

	mu sync.RWMutex
}

// 解析之后的订单路由结构体
type OrderRoutePolicy struct {
	JobId            string
	RouteMethod      RouteMethod
	SvcTimeOut       int
	PreTimeOut       int
	Batch            int
	MaxDelay         int
	StopFlag         StopFlag
	NeedCache        CacheFlag
	NeedKafka        KafkaFlag
	CacheTime        int
	SyncFlag         SyncFlag
	PolicyType       int
	Calllist         []string
	MemTaskIdMap     map[string]string
	MemConnObjIdMap  map[string]string
	MemKeyQryTypeMap map[string][]*TaskIdInfo
	CacheCodeList    []*CacheCode
}

func (o *OrderRoute) LoadOrderRouteFromXml(jobId string) (err error) {

	xmlName := fmt.Sprintf("order_route_%s.xml", jobId)
	xmlDir := settings.GetCommonSettings().Conf.XmlDir
	xmlPath := path.Join(xmlDir, xmlName)
	xmlFile, err := ioutil.ReadFile(xmlPath)
	if err != nil {
		logger.Error("read order route %s error ", xmlPath, err)
		panic(err)
	}

	err = xml.Unmarshal(xmlFile, o)
	if err != nil {
		logger.Error("Unmarshal order route %s error ", xmlPath, err)
		panic(err)
	}
	return
}

func (this *OrderRoute) GetConnObjIdByTaskId(taskId string) string {
	orderManager, err := order.GetOrderManager()
	if err != nil {
		fmt.Println("GetOrderMananger err :", err)
		return ""
	}
	orderInfo := orderManager.GetOrderInfo()
	orderDtl := orderInfo.GetTaskInfo(taskId)
	if orderDtl == nil {
		logger.Debug("while init order route taskId:%s not found in order info", taskId)
		return ""
	}

	return orderDtl.ConnObjID
}

func (this OrderRoute) GetConnObjListStr(taskIdStr string) string {
	s := make([]string, 0)
	for _, taskId := range strings.Split(taskIdStr, "|@|") {
		s = append(s, this.GetConnObjIdByTaskId(taskId))
	}

	return strings.Join(s, "|@|")
}

func (or *OrderRoute) getConnObjListStr(jobId string, taskIdInfos []*TaskIdInfo) string {
	s := make([]string, 0)
	for _, taskIdInfo := range taskIdInfos {
		s = append(s, or.getConnObjIdByTaskId(jobId, taskIdInfo.TaskId))
	}

	return strings.Join(s, "|@|")
}

func (or *OrderRoute) getConnObjIdByTaskId(jobId string, taskId string) string {
	orderData, ok := order.GetOrderInfoMap()[jobId]
	if !ok {
		return ""
	}
	orderInfo, ok := orderData.TaskInfoMapById[taskId]
	if !ok {
		return ""
	}

	return orderInfo.ConnObjId
}

// 旧版本根据orderId选择路由策略，新版本根据jobId（工单号）选择路由策略 ————lzq 2017.12.07
func (this OrderRoute) GetOrderRoute(orderId string) (OrderRoutePolicy, *errors.MeanfulError) {
	var policy OrderRoutePolicy
	memCache := cache.GetMemCacheInstance()
	cache_key := fmt.Sprintf("orderRoute_%s", orderId)
	order_route := memCache.GetMemCache(cache_key)
	if order_route != nil {
		logger.Debug("order_route  find orderRoute_%s", orderId)
		return order_route.(OrderRoutePolicy), nil
	} else {
		logger.Debug("order_route cannot find orderRoute_%s", orderId)
		return policy, errors.RawNew(ErrorPanic, "order_route cannot find")
	}
}

func (this OrderRoute) LoadOrderRouteXml(policy OrderRoutePolicy, orderRoute OrderRoute, orderId string) (OrderRoutePolicy, *errors.MeanfulError) {

	xmlName := fmt.Sprintf("order_route_%s.xml", orderId)
	xmlDir := settings.GetCommonSettings().Conf.XmlDir
	xmlPath := path.Join(xmlDir, xmlName)
	xmlFile, err := ioutil.ReadFile(xmlPath)
	if err != nil {
		logger.Error("read order route %s error ", xmlPath, err)
		panic(err)
	}

	err = xml.Unmarshal(xmlFile, &orderRoute)
	if err != nil {
		logger.Error("Unmarshal order route %s error ", xmlPath, err)
		panic(err)
	}

	memTaskIdMap := make(map[string]string)
	memConnObjIdMap := make(map[string]string)
	memKeyQryTypeMap := make(map[string][]*TaskIdInfo)
	cacheCodeList := make([]*CacheCode, 0)
	for _, svcInfo := range orderRoute.SvcList.SvcInfo {
		taskIdArr := make([]string, 0)
		memKeyQryTypeMap[svcInfo.MemId] = svcInfo.TaskIdList.TaskIdInfo
		for _, taskIdInfo := range svcInfo.TaskIdList.TaskIdInfo {
			taskIdArr = append(taskIdArr, taskIdInfo.TaskId)
		}
		taskIdStr := strings.Join(taskIdArr, "|@|")
		memTaskIdMap[svcInfo.MemId] = taskIdStr
		memConnObjIdMap[svcInfo.MemId] = this.GetConnObjListStr(taskIdStr)
	}

	for _, cacheCode := range orderRoute.Head.CacheCodeList.CacheCodes {
		cacheCodeList = append(cacheCodeList, cacheCode)
	}

	callist := make([]string, 0)
	routeMethod := orderRoute.Head.RouteMethod
	if routeMethod == 1 {
		callist = append(callist, orderRoute.SinglePolicy.MemId)
	} else if routeMethod == 2 {
		fmt.Println("routeMethod :", routeMethod)
		for _, memId := range orderRoute.StaticPolicy.CallList.MemId {
			callist = append(callist, memId)
		}
	} else {
		panic("route method TO-DO")
	}

	policy.JobId = orderId
	policy.RouteMethod = orderRoute.Head.RouteMethod
	policy.SvcTimeOut = orderRoute.Head.SvcTimeOut
	policy.PreTimeOut = orderRoute.Head.PreTimeOut
	policy.Batch = orderRoute.Head.Batch
	policy.MaxDelay = orderRoute.Head.MaxDelay
	policy.StopFlag = orderRoute.Head.StopFlag
	policy.NeedCache = orderRoute.Head.NeedCache
	policy.NeedKafka = orderRoute.Head.NeedKafka
	policy.CacheTime = orderRoute.Head.CacheTime
	policy.SyncFlag = orderRoute.Head.SyncFlag
	policy.PolicyType = orderRoute.DynamicPolicy.PolicyType
	policy.Calllist = callist
	policy.MemTaskIdMap = memTaskIdMap
	policy.MemConnObjIdMap = memConnObjIdMap
	policy.MemKeyQryTypeMap = memKeyQryTypeMap
	policy.CacheCodeList = cacheCodeList

	timeout := settings.GetCommonSettings().Conf.XmlReloadTime
	memCache := cache.GetMemCacheInstance()
	cache_key := fmt.Sprintf("orderRoute_%s", orderId)
	memCache.SetMemCache(cache_key, policy, timeout)
	logger.Info("saveing policy while load oder route: %+v", policy)
	return policy, nil
}

func (or *OrderRoute) LoadOrderRouteMap(jobId string) (*errors.MeanfulError) {
	fmt.Println("set order route config")
	memTaskIdMap := make(map[string]string)
	memConnObjIdMap := make(map[string]string)
	memKeyQryTypeMap := make(map[string][]*TaskIdInfo)
	cacheCodeList := make([]*CacheCode, 0)
	for _, svcInfo := range or.SvcList.SvcInfo {
		taskIdArr := make([]string, 0)
		taskIdList := svcInfo.TaskIdList.TaskIdInfo
		memKeyQryTypeMap[svcInfo.MemId] = taskIdList
		for _, taskIdInfo := range taskIdList {
			taskIdArr = append(taskIdArr, taskIdInfo.TaskId)
		}
		taskIdStr := strings.Join(taskIdArr, "|@|")
		memTaskIdMap[svcInfo.MemId] = taskIdStr
		memConnObjIdMap[svcInfo.MemId] = or.getConnObjListStr(jobId, taskIdList)
	}

	for _, cacheCode := range or.Head.CacheCodeList.CacheCodes {
		cacheCodeList = append(cacheCodeList, cacheCode)
	}

	callist := make([]string, 0)
	routeMethod := or.Head.RouteMethod
	if routeMethod == 1 {
		callist = append(callist, or.SinglePolicy.MemId)
	} else if routeMethod == 2 {
		for _, memId := range or.StaticPolicy.CallList.MemId {
			callist = append(callist, memId)
		}
	} else {
		logger.Error("route method %s error ", routeMethod)
		return errors.RawNew("999999", "route method error " + string(routeMethod))
	}

	policy := &OrderRoutePolicy{}
	policy.JobId = jobId
	policy.RouteMethod = or.Head.RouteMethod
	policy.SvcTimeOut = or.Head.SvcTimeOut
	policy.PreTimeOut = or.Head.PreTimeOut
	policy.Batch = or.Head.Batch
	policy.MaxDelay = or.Head.MaxDelay
	policy.StopFlag = or.Head.StopFlag
	policy.NeedCache = or.Head.NeedCache
	policy.NeedKafka = or.Head.NeedKafka
	policy.CacheTime = or.Head.CacheTime
	policy.SyncFlag = or.Head.SyncFlag
	policy.PolicyType = or.DynamicPolicy.PolicyType
	policy.Calllist = callist
	policy.MemTaskIdMap = memTaskIdMap
	policy.MemConnObjIdMap = memConnObjIdMap
	policy.MemKeyQryTypeMap = memKeyQryTypeMap
	policy.CacheCodeList = cacheCodeList

	or.mu.RLock()
	defer or.mu.RUnlock()
	OrderRoutePolicyMap[jobId] = policy

	return nil
}

//func GetOrderRoutePolicyMap() map[string]*OrderRoutePolicy {
//	return OrderRoutePolicyMap
//}

//初始化目录下所有OrderRouteFile
func InitOrderRouteFile() {

	xmlPath := settings.GetCommonSettings().Conf.XmlDir
	if xmlPath == "" {
		logger.Error("init orderRoute file empty ")
	}
	prefix := ORDER_ROUTE_PREFIX
	suffix := ORDER_ROUTE_SUFFIX
	orderIdArr := make([]string, 0)

	files, _, err := ListDir(xmlPath, prefix, suffix)
	if err != nil {
		logger.Error("go through order route dir error :", err)
		panic(err)
	}
	for i := range files {
		file := strings.TrimSuffix(strings.TrimPrefix(files[i], prefix), suffix)
		orderIdArr = append(orderIdArr, file)
	}
	for i := range orderIdArr {
		var policy OrderRoutePolicy
		var orderRoute OrderRoute
		OrderRoute{}.LoadOrderRouteXml(policy, orderRoute, orderIdArr[i])
	}
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, prefix string, suffix string) (files []string, filePaths []string, err error) {
	files = make([]string, 0, 10)
	filePaths = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasPrefix(fi.Name(), prefix) && strings.HasSuffix(fi.Name(), suffix) { //匹配文件
			files = append(files, fi.Name())
			filePaths = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, filePaths, nil
}
