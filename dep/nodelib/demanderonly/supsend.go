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
	"os"
	"dat/runtime/status"
)

func init() {
	SUPSEND.Register()
}

var SUPSEND = &DataBox{
	Name:         "supsend",
	Description:  "supsend",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: supsendRootFunc,

		Trunk: map[string]*Rule{
			"splitfile": {
				ParseFunc: splitfileFunc,
			},
			"pushfile": {
				ParseFunc: pushfileFunc,
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
		//LocalDir:       "/home/ddsdev/data/test/sup/send",
		LocalDir:       "D:/dds_send",
		LocalFileName:  fileName,
		RemoteDir:      remoteDir,
		RemoteFileName: fileName,
	}
	ctx.AddQueue(&request.DataRequest{
		Method:       "GET",
		FileCatalog:  fileCatalog,
		Rule:         "splitfile",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func splitfileFunc(ctx *Context) {
	fmt.Println("split file start ...")

	fileCatalog := ctx.DataRequest.FileCatalog
	localDir := fileCatalog.LocalDir
	fileName := fileCatalog.LocalFileName
	filePath := path.Join(localDir, fileName)
	tmpPath := path.Join(localDir, "tmp")
	os.MkdirAll(tmpPath, 0777)
	tmpFilePrefix := tmpPath + "/" + fileName + "__"

	fmt.Println("split file name:", filePath, tmpFilePrefix)
	cmdParams := []string{"-l", "1000000", filePath, tmpFilePrefix}

	ctx.AddChanQueue(&request.DataRequest{
		Rule:          "pushfile",
		TransferType:  request.SHELLTYPE,
		Priority:      1,
		CommandName:   "split",
		Reloadable:    true,
		CommandParams: cmdParams,
		FileCatalog:   fileCatalog,
	})
}

func pushfileFunc(ctx *Context) {
	fmt.Println("pushfile start ...")
	// 2. 从本地节点服务器通过sftp方式推送至dem节点服务器
	fileCatalog := ctx.DataRequest.FileCatalog
	//fileName := fileCatalog.LocalFileName

	localDir := fileCatalog.LocalDir
	tmpPath := path.Join(localDir, "tmp")

	fmt.Println("pushfile open path ...", tmpPath)

	dir, err := os.Open(tmpPath)
	if err != nil {
		fmt.Println("open file path error: ", err)
		return
	}

	names, err := dir.Readdirnames(0)
	if err != nil {
		fmt.Println("Readdirnames error: ", err)
		return
	}

	ctx.GetDataBox().DetailCount = len(names)

	fmt.Println("push file count: ", ctx.GetDataBox().DetailCount)

	for _, fileN := range names {
		fmt.Println("push file name:", fileN)

		ctx.AddChanQueue(&request.DataRequest{
			Url:          `http://10.101.12.17:8899/api/dem/rec`,
			Rule:         "notifydem",
			TransferType: request.HTTP,
			Method:       "PostFile",
			PostData:     "tmp/" + fileN,
			Reloadable:   true,
		})
	}
}

func notifydemFunc(ctx *Context) {
	fmt.Println("notifydem start ...")
	// 3. 通知dem节点服务器，继续往下执行
	//ctx.AddQueue(&request.DataRequest{
	//	Url:          "",
	//	Rule:         "notifydem",
	//	TransferType: request.NONETYPE,
	//})
	ctx.GetDataBox().ExecTsfSuccCount()
	fmt.Println("TsfSuccCount: ", ctx.GetDataBox().TsfSuccCount)

	if ctx.GetDataBox().TsfSuccCount == ctx.GetDataBox().DetailCount {

		fmt.Println("notifydem ok ...")
		defer ctx.GetDataBox().SetStatus(status.STOP)
	}

}
