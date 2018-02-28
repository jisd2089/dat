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
				ParseFunc: verifyFunc,
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

	fmt.Println("file name : ", dataFile.Filename)

	//targetFileDir := "/home/ddsdev/data/test/dem/rec"
	targetFileDir := "D:/dds_receive/tmp"
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
		UserName:       "bdaas",
		Password:       `bdaas`,
		Host:           "10.101.12.11",
		Port:           22,
		TimeOut:        10 * time.Second,
		LocalDir:       targetFileDir,
		LocalFileName:  targetFileName,
		RemoteDir:      "/home/bdaas/data/test/dem/rec",
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
	fmt.Println("demrec pushDoneFunc start ...", ctx.GetDataBox().GetId())

}


func verifyFunc(ctx *Context) {
	fmt.Println("demrec verify start ...")
	// 2. 将接收到的反馈文件推送至需方dmp
	ctx.AddQueue(&request.DataRequest{
		FileCatalog:  &sftp.FileCatalog{},
		Rule:         "pushdem",
		TransferType: request.SFTP,
	})
}