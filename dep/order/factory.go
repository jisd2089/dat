package order

import (
	logger "drcs/log"
	"drcs/settings"
	"fmt"
	"sync"
)

const (
	settings_xpath_filepath = "OrderFile"
	settings_xpath_expirs   = "OrderFileExpirS"
)

var (
	mutex               sync.Mutex
	defaultOrderManager OrderManager
)

func createOrderManager() (OrderManager, error) {
	settings := settings.GetCommonSettings()
	filePath := settings.ConfigFile.OrderFile
	if filePath == "" {
		logger.Error("OrderFile missing in setting:%s", settings_xpath_filepath)
		return nil, fmt.Errorf("配置缺失:%s", settings_xpath_filepath)
	}

	expirS := settings.ConfigFile.OrderFileExpireTime

	orderManager, err := NewOrderManagerXMLFile(filePath)
	if err != nil {
		logger.Error("get new orderManager from xml error ", err)
		return nil, err
	}

	expirManager := NewExpireOrderManager(expirS, orderManager)
	return expirManager, nil
}

// GetOrderManager OrderManager的工厂方法
func GetOrderManager() (OrderManager, error) {
	if defaultOrderManager != nil {
		return defaultOrderManager, nil
	}
	mutex.Lock()
	defer mutex.Unlock()
	if defaultOrderManager == nil {
		orderManager, err := createOrderManager()
		if err != nil {
			return nil, err
		}
		defaultOrderManager = orderManager
	}
	return defaultOrderManager, nil
}
