package service

import (
	"sync"
	"encoding/xml"

	"drcs/dep/agollo"
	"drcs/dep/order"
	"path/filepath"
	"fmt"
)

/**
    Author: luzequan
    Created: 2018-05-08 17:08:17
*/
type OrderService struct {
	lock    sync.RWMutex
	orderCh chan bool
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (o *OrderService) Init() {
	o.orderCh = make(chan bool, 1)

	path := filepath.Join(SettingPath, "order.properties")
	go o.initOrderConfig(filepath.Clean(path))

	// 初始化order route
	select {
	case ret := <-o.orderCh:
		fmt.Println(" order route init", ret)

		for _, o := range order.GetOrderInfos().Order {
			NewRouteService().InitByJobId(o.JobId)
		}

		break
	}
}

func (o *OrderService) initOrderConfig(configDir string) {
	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		fmt.Println("initOrderConfig")

		changesCnt := changeEvent.Changes["content"]
		value := changesCnt.NewValue

		switch changesCnt.ChangeType {
		case 0:
			orderInfoList := &order.OrderInfoList{}
			err := xml.Unmarshal([]byte(value), orderInfoList)
			if err != nil {
			}

			order.SetOrderInfos(orderInfoList)
		case 1:
			orderInfoList := &order.OrderInfoList{}
			err := xml.Unmarshal([]byte(value), orderInfoList)
			if err != nil {
			}
			order.SetOrderInfos(orderInfoList)
		}
		//bytes, _ := json.Marshal(changeEvent)
		//fmt.Println("event:", string(bytes))
		o.orderCh <- true
	}
}
