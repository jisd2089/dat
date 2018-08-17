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
	"path/filepath"
	"fmt"
)

type RouteService struct {
	lock       sync.RWMutex
}

func NewRouteService() *RouteService {
	return &RouteService{}
}

func (o *RouteService) Init() {
	path := filepath.Join(SettingPath, "route.properties")
	go initRouteConfig(filepath.Clean(path))
}

func initRouteConfig(configDir string) {
	newAgollo := agollo.NewAgollo(configDir)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		fmt.Println("initRouteConfig")

		changesCnt := changeEvent.Changes["content"]
		value := changesCnt.NewValue

		switch changesCnt.ChangeType {
		case 0:
			orderRoute := &or.OrderRoute{}
			err := xml.Unmarshal([]byte(value), orderRoute)
			if err != nil {
			}

			orderRoute.LoadOrderRouteMap("JON20180516000000431")

			//order.SetOrderInfos(orderInfoList)
		case 1:
			orderRoute := &or.OrderRoute{}
			err := xml.Unmarshal([]byte(value), orderRoute)
			if err != nil {
			}
			//order.SetOrderInfos(orderInfoList)
			orderRoute.LoadOrderRouteMap("JON20180516000000431")
		}
	}
}

func initRouteCfg(config *agollo.AppConfig, jobId string) {
	newAgollo := agollo.NewAgolloByConfig(config)
	go newAgollo.Start()

	event := newAgollo.ListenChangeEvent()
	for {
		changeEvent := <-event

		fmt.Println("initRouteConfig")

		changesCnt := changeEvent.Changes["content"]
		value := changesCnt.NewValue

		switch changesCnt.ChangeType {
		case 0:
			orderRoute := &or.OrderRoute{}
			err := xml.Unmarshal([]byte(value), orderRoute)
			if err != nil {
			}

			orderRoute.LoadOrderRouteMap(jobId)

			//order.SetOrderInfos(orderInfoList)
		case 1:
			orderRoute := &or.OrderRoute{}
			err := xml.Unmarshal([]byte(value), orderRoute)
			if err != nil {
			}
			//order.SetOrderInfos(orderInfoList)
			orderRoute.LoadOrderRouteMap(jobId)
		}
	}
}

func (o *RouteService) InitByJobId(jobId string) {
	config := &agollo.AppConfig{}
	config.AppId = "DLS"
	config.Cluster = "default"
	config.Ip = "10.101.12.29:8085"
	config.NamespaceName = fmt.Sprintf("route_%s.xml", jobId)
	go initRouteCfg(config, jobId)
}