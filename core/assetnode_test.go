package assetnode

/**
    Author: luzequan
    Created: 2018-01-02 11:18:25
*/
import (
	"testing"
	"runtime"
)

func init() {
	// 开启最大核心数运行
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestInit(t *testing.T) {

}
