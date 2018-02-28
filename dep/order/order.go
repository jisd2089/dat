package order

// OrderManager 订单管理器接口
type OrderManager interface {
	// 获取订单信息
	GetOrderInfo() *OrderInfo
}

// OrderInfo 订单信息结构体
type OrderInfo struct {
	TaskInfoMap map[string]*TaskInfo  //以TaskId作为key
	TaskInfoMapByConnObjID map[string]*TaskInfo //以ConnObjID作为key
}

func (orderInfo *OrderInfo) GetTaskInfo(taskID string) *TaskInfo {
	return orderInfo.TaskInfoMap[taskID]
}

func (orderInfo *OrderInfo) ConvertToMapByConnObjID() {
	taskInfoMapByConnObjID := make(map[string]*TaskInfo)

	for _, taskInfo := range orderInfo.TaskInfoMap {
		taskInfoMapByConnObjID[taskInfo.ConnObjID] = taskInfo
	}
	orderInfo.TaskInfoMapByConnObjID = taskInfoMapByConnObjID
}

// TaskInfo 任务信息结构体
type TaskInfo struct {
	TaskID          string
	SupMemID        string
	DemMemID        string
	ConnObjCatCd    string
	ConnObjNo       string
	ConnObjID       string
	PrdtIDCd        string
	ValuationModeCd string
	ValuationPrice  float64
	NeedCache       int
	CacheTime       int
	FeeCalDim       int
	EvalScore       int
	SvcType         string
}

// UpdatableOrderManager 对外提供更新方法的OrderManager接口扩展
type UpdatableOrderManager interface {
	OrderManager
	Update() error
}
