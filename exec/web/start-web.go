package web

import (
	"flag"
	"dat/core"
	//"log"
	//"net/http"
	//"os"
	//"os/exec"
	//"runtime"
	//"strconv"
	//"time"
	//

)

var (
	ip           *string
	port         *int
	addr         string
	dataBoxMenu []map[string]string
)

// 获取外部参数
func Flag() {
	flag.String("b ******************************************** only for web ******************************************** -b", "", "")
	// web服务器IP与端口号
	ip = flag.String("b_ip", "0.0.0.0", "   <Web Server IP>")
	port = flag.Int("b_port", 9090, "   <Web Server Port>")
}

// 执行入口
func Run(port int) {
	//assetnode.AssetNodeEntity.Init()

	assetnode.AssetNodeEntity.Run()

	httpServer := &HttpServer{}
	httpServer.Run(port)
}

func appInit() {
	//app.LogicApp.SetLog(Lsc).SetAppConf("Mode", cache.Task.Mode)

	dataBoxMenu = func() (dfmenu []map[string]string) {
		// 获取databox家族
		for _, df := range assetnode.AssetNodeEntity.GetDataBoxLib() {
			dfmenu = append(dfmenu, map[string]string{"name": df.GetName(), "description": df.GetDescription()})
		}
		return dfmenu
	}()
}
