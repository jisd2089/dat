package log

import (
	"drcs/log/appender"
	"drcs/log/formatter"
	"drcs/log/logs"
	"drcs/settings"
)

const (
	defaultConfigPath   = "logrus.yaml"
	setting_xpath_path = "Log.ConfigPath"
)

// Initialize  模块初始化方法
func Initialize() {
	appender.Initialize()
	formatter.Initialize()


	se := settings.GetCommomSettings()
	configPath := se.Log.ConfigPath
	if configPath == "" {
		Info("log initialize, use default configuration file: %s", defaultConfigPath)
		configPath = defaultConfigPath
	}
	err := logs.Initialize(configPath)
	if err != nil {
		Warn("log initialize error, configuration file: %s", configPath, err)
	}
}
