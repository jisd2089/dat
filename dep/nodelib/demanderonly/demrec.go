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
	"path"
	"time"
)

func init() {
	DEMREC.Register()
}

var DEMREC = &DataBox{
	Name:         "demrec",
	Description:  "demrec",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: rootFunc,

		Trunk: map[string]*Rule{
			"verify": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("demrec verify start ...")
					// 2. 将接收到的反馈文件推送至需方dmp
					ctx.AddQueue(&request.DataRequest{
						FileCatalog:  &sftp.FileCatalog{},
						Rule:         "pushdem",
						TransferType: request.SFTP,
					})
				},
			},
			"pushdem": {
				ParseFunc: pushdemFunc,
			},
			"pushDone": {
				ParseFunc: pushDoneFunc,
			},
		},
	},
}

func rootFunc(ctx *Context) {
	fmt.Println("demrec Root start ...")

	dataFile := ctx.GetDataBox().DataFile

	targetFileDir := "D:/input/SOURCE"
	targetFilePath := path.Join(targetFileDir, dataFile.Filename)

	ctx.AddQueue(&request.DataRequest{
		Rule:         "pushdem",
		TransferType: request.FILETYPE,
		DataFile:     ctx.GetDataBox().DataFile,
		PostData:     targetFilePath,
		Reloadable:   true,
	})
}

func pushdemFunc(ctx *Context) {
	fmt.Println("demrec pushdem start ...")

	targetFilePath := ctx.DataRequest.PostData
	targetFileDir := path.Dir(targetFilePath)
	targetFileName := path.Base(targetFilePath)

	fmt.Println("demrec pushdem filePath ...", targetFilePath)

	fileCatalog := &sftp.FileCatalog{
		UserName:       "ddsdev",
		Password:       `[BSR3+uLe\U*o^vy`,
		Host:           "10.101.12.17",
		Port:           22,
		TimeOut:        10 * time.Second,
		LocalDir:       targetFileDir,
		LocalFileName:  targetFileName,
		RemoteDir:      "/home/ddsdev/data/test/input",
		RemoteFileName: targetFileName,
	}

	ctx.AddQueue(&request.DataRequest{
		Method:       "PUT",
		FileCatalog:  fileCatalog,
		Rule:         "pushDone",
		TransferType: request.SFTP,
		Reloadable:   true,
	})
}

func pushDoneFunc(ctx *Context) {
	fmt.Println("demrec pushDoneFunc start ...")
}
