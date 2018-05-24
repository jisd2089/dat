package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"path"
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"fmt"
	"drcs/common/sftp"
	"time"
	"strings"
)

func init() {
	ALGSEND.Register()
}

var ALGSEND = &DataBox{
	Name:        "algorithmsend",
	Description: "algorithmsend",
	RuleTree: &RuleTree{
		Root: algsendRootFunc,

		Trunk: map[string]*Rule{
			"pull": {
				ParseFunc: pullFunc,
			},
			"process": {
				ParseFunc: processAlgFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func algsendRootFunc(ctx *Context) {
	fmt.Println("algorithmsend Root ...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "pull",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func pullFunc(ctx *Context) {
	fmt.Println("start rule...")

	filePath := ctx.GetDataBox().GetDataFilePath()
	dataFile := path.Base(filePath)
	dataFilePath := path.Dir(filePath)

	fmt.Println(dataFilePath + "@" + ctx.GetDataBox().GetDataFilePath())

	fsAddress := ctx.GetDataBox().FileServerAddress

	// 1. 从sftp服务器（需方dmp服务器）拉取文件到节点服务器本地
	fileCatalog := &sftp.FileCatalog{
		UserName:       fsAddress.UserName,
		Password:       fsAddress.Password,
		Host:           fsAddress.Host,
		Port:           fsAddress.Port,
		TimeOut:        time.Duration(fsAddress.TimeOut) * time.Second,
		LocalDir:       fsAddress.LocalDir,
		LocalFileName:  dataFile,
		RemoteDir:      fsAddress.RemoteDir,
		RemoteFileName: dataFile,
	}

	ctx.GetDataBox().SetDataFilePath(path.Join(fsAddress.LocalDir, dataFile))

	fmt.Println("NodeAddress: %s", ctx.GetDataBox().GetNodeAddress())
	ctx.AddQueue(&request.DataRequest{
		Rule:         "process",
		Method:       "GET",
		TransferType: request.SFTP,
		FileCatalog:  fileCatalog,
		Reloadable:   true,
	})
}

func processAlgFunc(ctx *Context) {
	fmt.Println("process ...")

	if ctx.DataResponse.StatusCode == 200 && strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
		addressList := ctx.GetDataBox().GetNodeAddress()

		addr := addressList[0]

		ctx.AddQueue(&request.DataRequest{
			Rule:         "end",
			Url:          addr.GetUrl(),
			Method:       "POSTFILE",
			TransferType: request.FASTHTTP,
			Priority:     1,
			PostData:     ctx.GetDataBox().DataFilePath,
			Reloadable:   true,
		})
	} else {
		ctx.AddQueue(&request.DataRequest{
			Rule:         "end",
			TransferType: request.NONETYPE,
			Priority:     1,
			Reloadable:   true,
		})
	}
}