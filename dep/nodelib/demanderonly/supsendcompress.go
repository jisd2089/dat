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
	SUPSENDCOMPRESS.Register()
}

var SUPSENDCOMPRESS = &DataBox{
	Name:         "supsendcompress",
	Description:  "supsendcompress",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: supsendCompressRootFunc,

		Trunk: map[string]*Rule{
			"compressfile": {
				ParseFunc: compressfileFunc,
			},
			"pushcompressfile": {
				ParseFunc: pushcompressfileFunc,
			},
			"notifydemcompress": {
				ParseFunc: notifydemcompressFunc,
			},
		},
	},
}

func supsendCompressRootFunc(ctx *Context) {
	fmt.Println("supsend compress Root start...")

	filePath := ctx.GetDataBox().DataFilePath

	fileName := path.Base(filePath)

	remoteDir := path.Dir(filePath)

	fmt.Println("supsend pull file name ...", fileName)

	// 1. 从sftp服务器（供方dmp服务器）拉取文件到节点服务器本地
	fileCatalog := &sftp.FileCatalog{
		UserName: "bdaas",
		Password: `bdaas`,
		Host:     "10.101.12.11",
		Port:     22,
		TimeOut:  10 * time.Second,
		//LocalDir:       "/home/ddsdev/data/test/sup/send",
		LocalDir:       "D:/dds_send/batch",
		LocalFileName:  fileName,
		RemoteDir:      remoteDir,
		RemoteFileName: fileName,
	}
	ctx.AddQueue(&request.DataRequest{
		Method:       "GET",
		FileCatalog:  fileCatalog,
		Rule:         "compressfile",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func compressfileFunc(ctx *Context) {
	fmt.Println("compress file start ...")

	fileCatalog := ctx.DataRequest.FileCatalog
	localDir := fileCatalog.LocalDir
	fileName := fileCatalog.LocalFileName
	filePath := path.Join(localDir, fileName)

	fmt.Println("compress file name:", filePath)

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "pushcompressfile",
		Method:       "COMPRESS",
		TransferType: request.FILETYPE,
		Priority:     1,
		Reloadable:   true,
		FileCatalog:  fileCatalog,
	})
}

func pushcompressfileFunc(ctx *Context) {
	fmt.Println("pushfile compress start ...")

	// 2. 从本地节点服务器通过sftp方式推送至dem节点服务器
	fileCatalog := ctx.DataRequest.FileCatalog
	fileName := fileCatalog.LocalFileName

	fmt.Println("push file name:", fileName)

	ctx.AddChanQueue(&request.DataRequest{
		Url:          `http://127.0.0.1:8899/api/dem/rec`,
		//Url:          `http://10.101.12.17:8899/api/dem/rec`,
		Rule:         "notifydemcompress",
		TransferType: request.FASTHTTP,
		Method:       "PostFile",
		PostData:     fileName,
		Reloadable:   true,
	})

}

func notifydemcompressFunc(ctx *Context) {
	fmt.Println("notifydem compress start ...")
	// 3. 通知dem节点服务器，继续往下执行
	//ctx.AddQueue(&request.DataRequest{
	//	Url:          "",
	//	Rule:         "notifydem",
	//	TransferType: request.NONETYPE,
	//})

	fmt.Println("notifydem ok ...")
	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
}
