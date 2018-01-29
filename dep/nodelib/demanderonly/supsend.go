package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"dat/core/interaction/request"
	. "dat/core/databox"
	"fmt"
	"dat/common/sftp"
	"time"
	"path"
)

func init() {
	SUPSEND.Register()
}

var SUPSEND = &DataBox{
	Name:        "supsend",
	Description: "supsend",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: supsendRootFunc,

		Trunk: map[string]*Rule{
			"pushfile": {
				ParseFunc:pushfileFunc,
			},
			"notifydem": {
				ParseFunc: notifydemFunc,
			},
		},
	},
}

func supsendRootFunc(ctx *Context) {
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
		LocalDir:       "/home/ddsdev/data/test/sup/send",
		LocalFileName:  fileName,
		RemoteDir:      remoteDir,
		RemoteFileName: fileName,
	}
	ctx.AddQueue(&request.DataRequest{
		Method:       "GET",
		FileCatalog:  fileCatalog,
		Rule:         "pushfile",
		TransferType: request.SFTP,
		Reloadable:   true,
	})
}

func pushfileFunc(ctx *Context) {
	fmt.Println("pushfile start ...")
	// 2. 从本地节点服务器通过sftp方式推送至dem节点服务器
	fileCatalog := ctx.DataRequest.FileCatalog
	fileName := fileCatalog.LocalFileName

	fmt.Println("push file name:", fileName)

	ctx.AddQueue(&request.DataRequest{
		Url:          `http://10.101.12.17:8899/api/dem/rec`,
		Rule:         "notifydem",
		TransferType: request.HTTP,
		Method:       "PostFile",
		PostData:     fileName,
		Reloadable:   true,
	})
}

func notifydemFunc(ctx *Context) {
	fmt.Println("notifydem start ...")
	// 3. 通知dem节点服务器，继续往下执行
	//ctx.AddQueue(&request.DataRequest{
	//	Url:          "",
	//	Rule:         "notifydem",
	//	TransferType: request.NONETYPE,
	//})
}
