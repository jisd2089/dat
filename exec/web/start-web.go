package web

import (
	"flag"
	"drcs/core"

	_ "drcs/dep/service"

	//"drcs/dep/job"
	"drcs/dep/service"
)

var (
	ip          *string
	port        *int
	addr        string
	dataBoxMenu []map[string]string

	SettingPath *string
)



// 获取外部参数
func Flag() {
	flag.String("b ******************************************** only for web ******************************************** -b", "", "")
	// web服务器IP与端口号
	ip = flag.String("b_ip", "0.0.0.0", "   <Web Server IP>")
	port = flag.Int("b_port", 9090, "   <Web Server Port>")

	SettingPath = flag.String("c", "/dep/go/conf/dem-settings.yaml", "setting config path")
}

// 执行入口
func Run() {

	service.NewNodeService().Init()

	assetnode.AssetNodeEntity.Init()

	assetnode.AssetNodeEntity.Run()

	httpServer := &HttpServer{}
	httpServer.Run()
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

/**
 **********************************************************************************************************
 */

func RunDem() {

	//jobs.InitDem()

	Run()
}

func RunSup() {

	//jobs.InitSup()

	Run()
}
