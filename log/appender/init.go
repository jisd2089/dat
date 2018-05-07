package appender

import "drcs/log/logs"

func Initialize() {
	// 注册到yagrus
	logs.RegisgerAppenders("NormalFile", NewNormalFileAppender)
	logs.RegisgerAppenders("RotationFile", NewRotateFileAppender)
	logs.RegisgerAppenders("RollingFile", NewRollingFileAppender)
}
