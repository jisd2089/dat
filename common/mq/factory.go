package mq

import (
	logger "drcs/log"
	"drcs/settings"
	"fmt"
	"sync"
	"drcs/common/mq/posixmq"
)

const (
	settings_xpath_mqpath = "busilog.MQPath"
)

var (
	mutex           sync.Mutex
	defaultRecordor Recordor
)

func createRecordor() (Recordor, error) {
	settings := settings.GetCommonSettings()
	mqPath := settings.BusiLog.MQPath

	if mqPath == "" {
		logger.Error(" %s missing in setting", settings_xpath_mqpath)
		return nil, fmt.Errorf("配置缺失:%s", settings_xpath_mqpath)
	}

	//fmt.Println("mqPath: ", mqPath)
	recordor, err := posixmq.New(mqPath)
	if err != nil {
		logger.Error("create posixmq error ", err)
		return nil, err
	}
	return recordor, nil
}

func GetRecordor() (Recordor, error) {
	if defaultRecordor != nil {
		return defaultRecordor, nil
	}
	mutex.Lock()
	defer mutex.Unlock()
	if defaultRecordor == nil {
		recordor, err := createRecordor()
		if err != nil {
			return nil, err
		}
		defaultRecordor = recordor
	}
	return defaultRecordor, nil
}
