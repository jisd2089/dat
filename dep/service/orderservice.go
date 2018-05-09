package service

import (
	"sync"
	"encoding/xml"

	"drcs/dep/agollo"
	"drcs/dep/order"
)

/**
    Author: luzequan
    Created: 2018-05-08 17:08:17
*/

func init() {
	//NewOrderService().Init()
}

type OrderService struct {
	lock       sync.RWMutex
}

func NewOrderService() *OrderService {
	return &OrderService{
	}
}

func (o *OrderService) Init() {

	initOrderConfig("D:/GoglandProjects/src/drcs/dep/order/order.properties")

}

func initOrderConfig(configDir string) {
	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

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
	}
}