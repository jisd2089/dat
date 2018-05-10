package order

/**
    Author: luzequan
    Created: 2018-05-08 17:27:55
*/
var (
	orderInfoList *OrderInfoList
	orderInfoMap map[string]*OrderData
)

type OrderInfoList struct {
	Head  *Head    `xml:"head"`
	Order []*Order `xml:"order_info"`
}

type Head struct {
	FileName          string `xml:"fileName"`
	FileCreateTime    string `xml:"fileCreateTime"`
	FileCreateTimeStr string `xml:"fileCreateTimeStr"`
}

type Order struct {
	OrderId         string           `xml:"orderId"`
	MemberRole      string           `xml:"memberRole"`
	JobId           string           `xml:"jobId"`
	OrderDetailList *OrderDetailList `xml:"order_dtl_list"`
}

type OrderDetailList struct {
	OrderDetailInfo []*OrderDetailInfo `xml:"order_dtl_info"`
}

type OrderDetailInfo struct {
	TaskId           string `xml:"taskId"`
	SupMemId         string `xml:"supMemId"`
	DemMemId         string `xml:"demMemId"`
	ConnObjCatCd     string `xml:"connObjCatCd"`
	ConnObjNo        string `xml:"connObjNo"`
	ConnObjId        string `xml:"connObjId"`
	ConnObjIdVersion string `xml:"connObjIdVersion"`
	PrdtIdCd         string `xml:"prdtIdCd"`
	ValuationModeCd  string `xml:"valuationModeCd"`
	ValuationPrice   string `xml:"valuationPrice"`
	NeedCache        string `xml:"needCache"`
	CacheTime        string `xml:"cacheTime"`
	FeeCalDim        string `xml:"feeCalDim"`
	EvalScore        string `xml:"evalScore"`
	SvcType          string `xml:"svcType"`
}

type OrderData struct {
	TaskInfoMapById        map[string]*OrderDetailInfo //以TaskId作为key
	TaskInfoMapByConnObjID map[string]*OrderDetailInfo //以ConnObjID作为key
}

func SetOrderInfos(orderInfos *OrderInfoList) *OrderInfoList {
	orderInfoList = orderInfos
	SetOrderInfoMap(orderInfos)
	return orderInfoList
}

func GetOrderInfos() *OrderInfoList {
	return orderInfoList
}

func SetOrderInfoMap(orderInfos *OrderInfoList) {
	if orderInfoMap == nil {
		orderInfoMap = make(map[string]*OrderData)
	}

	for _, o := range orderInfos.Order {
		taskInfoMapById := make(map[string]*OrderDetailInfo)
		taskInfoMapByConnObjID := make(map[string]*OrderDetailInfo)
		for _, d := range o.OrderDetailList.OrderDetailInfo {
			taskInfoMapById[d.TaskId] = d
			taskInfoMapByConnObjID[d.ConnObjId] = d
		}
		orderInfoMap[o.JobId] = &OrderData{
			taskInfoMapById,
			taskInfoMapByConnObjID,
		}
	}
}

func GetOrderInfoMap() map[string]*OrderData {
	return orderInfoMap
}


