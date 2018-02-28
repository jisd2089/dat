package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"fmt"
	"drcs/common/sftp"
	"time"
	"path"
	"drcs/runtime/status"
)

func init() {
	SUPSENDNOSPLIE.Register()
}

var SUPSENDNOSPLIE = &DataBox{
	Name:         "supsendnotsplit",
	Description:  "supsendnotsplit",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: supsendnotsplitRootFunc,

		Trunk: map[string]*Rule{
			"pushfile": {
				ParseFunc: pushFileNotSplitFunc,
			},
			"notifydem": {
				ParseFunc: notifyDemNotSplitFunc,
			},
		},
	},
}

func supsendnotsplitRootFunc(ctx *Context) {
	fmt.Println("supsend Root start...")

	filePath := ctx.GetDataBox().DataFilePath

	fileName := path.Base(filePath)

	remoteDir := path.Dir(filePath)

	fmt.Println("supsend pull file name ...", fileName)

	// 1. 从sftp服务器（供方dmp服务器）拉取文件到节点服务器本地
	fileCatalog := &sftp.FileCatalog{
		UserName:       "bdaas",
		Password:       `bdaas`,
		Host:           "10.101.12.11",
		Port:           22,
		TimeOut:        10 * time.Second,
		//LocalDir:       "/home/ddsdev/data/test/sup/send",
		LocalDir:       "D:/dds_send/tmp",
		LocalFileName:  fileName,
		RemoteDir:      remoteDir,
		RemoteFileName: fileName,
	}
	ctx.AddQueue(&request.DataRequest{
		Method:       "GET",
		FileCatalog:  fileCatalog,
		Rule:         "pushfile",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func pushFileNotSplitFunc(ctx *Context) {
	fmt.Println("pushfile start ...")
	// 2. 从本地节点服务器通过sftp方式推送至dem节点服务器
	fileCatalog := ctx.DataRequest.FileCatalog
	fileName := fileCatalog.LocalFileName

	ctx.AddChanQueue(&request.DataRequest{
		Url:          `http://127.0.0.1:8899/api/dem/rec`,
		//Url:          `http://10.101.12.17:8899/api/dem/rec`,
		Rule:         "notifydem",
		TransferType: request.FASTHTTP,
		Method:       "PostFile",
		PostData:     fileName,
		Reloadable:   true,
	})
}

func notifyDemNotSplitFunc(ctx *Context) {
	fmt.Println("notifydem start ...")

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
	// 3. 通知dem节点服务器，继续往下执行
	//ctx.AddQueue(&request.DataRequest{
	//	Url:          "",
	//	Rule:         "notifydem",
	//	TransferType: request.NONETYPE,
	//})
}
