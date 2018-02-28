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
	"drcs/runtime/status"
)

func init() {
	DEMRECUNCOMPRESS.Register()
}

var DEMRECUNCOMPRESS = &DataBox{
	Name:         "demrecuncompress",
	Description:  "demrecuncompress",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: demrecUncompressRootFunc,

		Trunk: map[string]*Rule{
			"pushuncompressdem": {
				ParseFunc: pushuncompressdemFunc,
			},
			"pushuncompressdone": {
				ParseFunc: pushUncompressDoneFunc,
			},
		},
	},
}

func demrecUncompressRootFunc(ctx *Context) {
	fmt.Println("demrec uncompress Root start ...")

	dataFilePath := ctx.GetDataBox().DataFilePath

	fmt.Println("file path name : ", dataFilePath)

	fileName := path.Base(dataFilePath)

	//targetFileDir := "/home/ddsdev/data/test/dem/rec"
	//targetFileDir := "D:/dds_receive/tmp"
	//targetFilePath := path.Join(targetFileDir, fileName)

	fileCatalog := &sftp.FileCatalog{
		LocalDir:       "D:/dds_receive/tmp",
		LocalFileName:  fileName,
		RemoteDir:      "D:/dds_receive/tmp/uncompress",
		RemoteFileName: fileName,
	}
	ctx.AddQueue(&request.DataRequest{
		Rule:         "pushuncompressdem",
		TransferType: request.FILETYPE,
		Method:       "UNCOMPRESS",
		FileCatalog:  fileCatalog,
		Reloadable:   true,
	})
}

func pushuncompressdemFunc(ctx *Context) {
	fmt.Println("demrec uncompress pushdem start ...")

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
		Rule:         "pushuncompressdone",
		TransferType: request.SFTP,
		Reloadable:   true,
	})
}

func pushUncompressDoneFunc(ctx *Context) {
	fmt.Println("demrec pushDoneFunc start ...", ctx.GetDataBox().GetId())
	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
}
