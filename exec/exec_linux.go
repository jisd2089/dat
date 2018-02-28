package exec

import (
	"drcs/exec/web" // web版
)

func run(port int) {
	//exec.Command("/bin/sh", "-c", "title", config.FULL_NAME).Start()
	//
	//// 选择运行界面
	//switch which {
	//case "cmd":
	//	cmd.Run()
	//
	//case "web":
	//	fallthrough
	//default:
	//	ctrl := make(chan os.Signal, 1)
	//	signal.Notify(ctrl, os.Interrupt, os.Kill)
	//	go web.Run()
	//	<-ctrl
	//}
	web.Run(port)
}

func runNode(role string) {
	switch role {
	case "dem":
		web.RunDem()
	case "sup":
		web.RunSup()
	}
}
