package appender

import "dds/log/yagrus"

func Initialize() {
	// 注册到yagrus
	yagrus.RegisgerAppenders("NormalFile", NewNormalFileAppender)
	yagrus.RegisgerAppenders("RotationFile", NewRotateFileAppender)
	yagrus.RegisgerAppenders("RollingFile", NewRollingFileAppender)
}
