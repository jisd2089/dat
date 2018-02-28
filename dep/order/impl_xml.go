package order

import (
	"encoding/xml"
	"io/ioutil"
	"sync/atomic"
	logger "dds/log"
)

// OrderManagerXMLFile 基于XML文件的OrderManager实现
type OrderManagerXMLFile struct {
	filePath string
	value    *atomic.Value
}

func (manager *OrderManagerXMLFile) GetOrderInfo() *OrderInfo {
	orderInfo := manager.value.Load().(*OrderInfo)
	return orderInfo
}

// Update 重新读取XML文件，更新订单信息
func (manager *OrderManagerXMLFile) Update() error {
	orderInfo, err := parseXMLFile(manager.filePath)
	if err != nil {
		return err
	}
	manager.value.Store(orderInfo)
	return nil
}

// NewOrderManagerXMLFile 创建OrderManagerXMLFile实例
func NewOrderManagerXMLFile(filePath string) (*OrderManagerXMLFile, error) {
	orderInfo, err := parseXMLFile(filePath)
	if err != nil {
		return nil, err
	}

	manager := &OrderManagerXMLFile{}
	manager.filePath = filePath
	manager.value = &atomic.Value{}
	manager.value.Store(orderInfo)
	return manager, nil
}

func parseXMLFile(filePath string) (*OrderInfo, error) {
	fd, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error("read order info:%s error", filePath, err)
		return nil, err
	}

	var orderInfoXML orderInfoXML
	err = xml.Unmarshal(fd, &orderInfoXML)
	if err != nil {
		logger.Error("unmarshal order info:%s error", filePath, err)
		return nil, err
	}

	return orderInfoXML2OrderInfo(&orderInfoXML), nil
}

func orderInfoXML2OrderInfo(orderInfoXML *orderInfoXML) *OrderInfo {
	taskInfoMap := make(map[string]*TaskInfo)

	for _, taskInfoXML := range orderInfoXML.TaskInfoXMLList.TaskInfoXMLs {
		taskInfo := &TaskInfo{
			TaskID:          taskInfoXML.TaskID,
			SupMemID:        taskInfoXML.SupMemID,
			DemMemID:        taskInfoXML.DemMemID,
			ConnObjCatCd:    taskInfoXML.ConnObjCatCd,
			ConnObjNo:       taskInfoXML.ConnObjNo,
			ConnObjID:       taskInfoXML.ConnObjID,
			PrdtIDCd:        taskInfoXML.PrdtIDCd,
			ValuationModeCd: taskInfoXML.ValuationModeCd,
			ValuationPrice:  taskInfoXML.ValuationPrice,
			NeedCache:       taskInfoXML.NeedCache,
			CacheTime:       taskInfoXML.CacheTime,
			FeeCalDim:       taskInfoXML.FeeCalDim,
			EvalScore:       taskInfoXML.EvalScore,
			SvcType:         taskInfoXML.SvcType,
		}
		taskInfoMap[taskInfo.TaskID] = taskInfo
		logger.Info("init taskinfo: %+v", taskInfo)

	}
	return &OrderInfo{
		TaskInfoMap: taskInfoMap,
	}
}

type taskInfoXML struct {
	TaskID          string  `xml:"taskId"`
	SupMemID        string  `xml:"supMemId"`
	DemMemID        string  `xml:"demMemId"`
	ConnObjCatCd    string  `xml:"connObjCatCd"`
	ConnObjNo       string  `xml:"connObjNo"`
	ConnObjID       string  `xml:"connObjId"`
	PrdtIDCd        string  `xml:"prdtIdCd"`
	ValuationModeCd string  `xml:"valuationModeCd"`
	ValuationPrice  float64 `xml:"valuationPrice"`
	NeedCache       int     `xml:"needCache"`
	CacheTime       int     `xml:"cacheTime"`
	FeeCalDim       int     `xml:"feeCalDim"`
	EvalScore       int     `xml:"evalScore"`
	SvcType         string  `xml:"svcType"`
}

type taskInfoXMLList struct {
	TaskInfoXMLs []taskInfoXML `xml:"order_dtl_info"`
}

type orderInfoXML struct {
	TaskInfoXMLList taskInfoXMLList `xml:"order_dtl_list"`
}
