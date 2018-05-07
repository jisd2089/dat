package formatter

import "drcs/log/logs"

func Initialize() {
	// 注册到yagrus
	logs.RegisgerFormatter("PatternLayout", NewPatternLayout)
}
