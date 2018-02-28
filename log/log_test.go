package log

import (
	"fmt"
	"testing"
	"time"
)

func TestInfo(t *testing.T) {
	Initialize()
	Info("abcdef")
	Info("sddddddddddddddddddddddddddddddddddabcdef")
	// logger.Info("abcdef :%s", "fffff", fmt.Errorf("this is the reason"))

	appLogger := GetLogger("app")
	//appLogger.Info("abcdef")
	//appLogger.Info("abcdef :%s", "fffff", fmt.Errorf("this is the reason"))

	for {
		appLogger.Info("abcdef")
		appLogger.Info("abcdef :%s", "fffff", fmt.Errorf("this is the reason"))
		//if time.Now().Unix() - curT > 60{
		//	break
		//}
		time.Sleep(10000000000)
	}
}
