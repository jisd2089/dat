package web

import (
	"flag"
	"drcs/core"
	"drcs/dep/service"
	"fmt"
	"runtime"
)

var (
	ip          *string
	port        *int
	addr        string
	dataBoxMenu []map[string]string

	settingPath *string
)

func init() {
	Flag()
	service.SettingPath = *settingPath
}

// 获取外部参数
func Flag() {
	flag.String("b ******************************************** only for web ******************************************** -b", "", "")
	// web服务器IP与端口号
	ip = flag.String("b_ip", "0.0.0.0", "   <Web Server IP>")
	port = flag.Int("b_port", 9090, "   <Web Server Port>")

	settingPath = flag.String("c", "D:/GoglandProjects/src/drcs/settings/properties", "setting config path")

	flag.Parse()
	fmt.Println("flag settingPath: ", *settingPath)
}

// 执行入口
func Run() {

	runtime.GOMAXPROCS(8)

	service.NewNodeService().Init()

	assetnode.AssetNodeEntity.Init().Run()

	httpServer := &HttpServer{}
	httpServer.Run()
}

func RunTest(port int) {

	service.NewNodeService().Init()

	assetnode.AssetNodeEntity.Init().Run()

	httpServer := &HttpServer{}
	httpServer.RunTest(port)
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
