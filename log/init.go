package log

import (
	"drcs/log/appender"
	"drcs/log/formatter"
	"drcs/log/logs"
	"drcs/settings"
	"drcs/common/util"
)

const (
	defaultConfigPath   = "D:/GoglandProjects/src/drcs/log/logrus.yaml"
	setting_xpath_path = "Log.ConfigPath"
)

// Initialize  模块初始化方法
func Initialize() {
	appender.Initialize()
	formatter.Initialize()


	se := settings.GetCommonSettings()
	configPath := se.Log.ConfigPath

	if configPath == "" || !util.IsFileExists(configPath) {
		Info("log initialize, use default configuration file: %s", defaultConfigPath)
		configPath = defaultConfigPath
	}

	err := logs.Initialize(configPath)
	if err != nil {
		Warn("log initialize error, configuration file: %s", configPath, err)
	}
}
