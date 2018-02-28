package formatter

import "dds/log/yagrus"

func Initialize() {
	// 注册到yagrus
	yagrus.RegisgerFormatter("PatternLayout", NewPatternLayout)
}
