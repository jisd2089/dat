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
	DEMRECBIG.Register()
}

var DEMRECBIG = &DataBox{
	Name:         "demrecbig",
	Description:  "demrecbig",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: rootBigFunc,

		Trunk: map[string]*Rule{
			"verify": {
				ParseFunc: verifyFunc,
			},
			"pushdem": {
				ParseFunc: pushBigDemFunc,
			},
			"pushDone": {
				ParseFunc: pushBigDoneFunc,
			},
		},
	},
}

func rootBigFunc(ctx *Context) {
	fmt.Println("demrec Root start ...")

	//dataFile := ctx.GetDataBox().DataFile
	//
	//fmt.Println("file name : ", dataFile.Filename)
	//
	////targetFileDir := "/home/ddsdev/data/test/dem/rec"
	//targetFileDir := "D:/dds_receive/tmp"
	//targetFilePath := path.Join(targetFileDir, dataFile.Filename)

	ctx.AddQueue(&request.DataRequest{
		Rule:         "pushdem",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func pushBigDemFunc(ctx *Context) {
	fmt.Println("demrec pushdem start ...")

	targetFilePath := ctx.GetDataBox().DataFilePath
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

func pushBigDoneFunc(ctx *Context) {
	fmt.Println("demrec pushDoneFunc start ...", ctx.GetDataBox().GetId())

}
