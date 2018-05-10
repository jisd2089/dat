package service

/**
    Author: luzequan
    Created: 2018-05-10 10:02:54
*/

import (
	"sync"
	"encoding/xml"

	"drcs/dep/agollo"
	"drcs/dep/or"
)

func init() {
	NewRouteService().Init()
}

type RouteService struct {
	lock       sync.RWMutex
}

func NewRouteService() *RouteService {
	return &RouteService{}
}

func (o *RouteService) Init() {

	initRouteConfig("D:/GoglandProjects/src/drcs/dep/or/route.properties")

}

func initRouteConfig(configDir string) {
	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		changesCnt := changeEvent.Changes["content"]
		value := changesCnt.NewValue

		switch changesCnt.ChangeType {
		case 0:
			orderRoute := &or.OrderRoute{}
			err := xml.Unmarshal([]byte(value), orderRoute)
			if err != nil {
			}

			orderRoute.LoadOrderRouteMap("ODN20161222000000071")

			//order.SetOrderInfos(orderInfoList)
		case 1:
			orderRoute := &or.OrderRoute{}
			err := xml.Unmarshal([]byte(value), orderRoute)
			if err != nil {
			}
			//order.SetOrderInfos(orderInfoList)
		}
	}
}