package demanderonly

/**
    Author: luzequan
    Created: 2018-05-14 10:19:48
*/
import (
	"path"
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"fmt"
	"drcs/dep/management/util"
	"drcs/common/sftp"
	"drcs/runtime/status"
	"time"
	"strings"
)

func init() {
	ALGRCV.Register()
}

var ALGRCV = &DataBox{
	Name:        "algorithmreceive",
	Description: "algorithmreceive",
	RuleTree: &RuleTree{
		Root: algrcvRootFunc,

		Trunk: map[string]*Rule{
			"pushtoserver": {
				ParseFunc: pushToServerFunc,
			},
			"puttohdfs": {
				ParseFunc: putToHDFSFunc,
			},
			"dataready": {
				ParseFunc: datareadyFunc,
			},
			"runtask": {
				ParseFunc: runTaskFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func algrcvRootFunc(ctx *Context) {
	fmt.Println("algorithmreceive Root ...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "pushtoserver",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func pushToServerFunc(ctx *Context) {
	fmt.Println("algorithmreceive push to server...")

	filePath := ctx.GetDataBox().GetDataFilePath()
	dataFile := path.Base(filePath)
	dataFilePath := path.Dir(filePath)
	dataFileName := &util.DataFileName{}
	if err := dataFileName.ParseAndValidFileName(dataFile); err != nil {
		errEnd(ctx)
		return
	}

	fmt.Println(dataFilePath)

	fmt.Println(ctx.GetDataBox().GetDataFilePath())

	fsAddress := ctx.GetDataBox().FileServerAddress

	// 1. push local file to hadoop server
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

	ctx.GetDataBox().SetDataFilePath(path.Join(fsAddress.RemoteDir, dataFile))

	ctx.AddQueue(&request.DataRequest{
		Rule:         "puttohdfs",
		Method:       "PUT",
		TransferType: request.SFTP,
		FileCatalog:  fileCatalog,
		Reloadable:   true,
	})
}

func putToHDFSFunc(ctx *Context) {
	fmt.Println("algorithmreceive put to hdfs...")

	if ctx.GetResponse().StatusCode != 200 {
		errEnd(ctx)
		return
	}

	// hadoop server local file path
	filePath := ctx.GetDataBox().GetDataFilePath()
	dataFile := path.Base(filePath)
	dataFilePath := path.Dir(filePath)
	dataFileName := &util.DataFileName{}
	if err := dataFileName.ParseAndValidFileName(dataFile); err != nil {
		errEnd(ctx)
		return
	}

	fmt.Println(dataFilePath)

	fmt.Println(ctx.GetDataBox().GetDataFilePath())

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

	cmdParams := []string{}
	cmdParams = append(cmdParams, "spark-submit")
	cmdParams = append(cmdParams, "--class chinadep.Precollision")
	cmdParams = append(cmdParams, "--master yarn")
	cmdParams = append(cmdParams, "--deploy-mode cluster")
	cmdParams = append(cmdParams, "--queue sparkqueue")
	cmdParams = append(cmdParams, "--jars hdfs://deptest20:9000/user/aarontest/aaron-oozie/shell-spark-Precollision/lib/GenExidRealTime.jar")

	fmt.Println("NodeAddress: %s", ctx.GetDataBox().GetNodeAddress())
	ctx.AddQueue(&request.DataRequest{
		Rule:          "dataready",
		Method:        "CMD",
		TransferType:  request.SSH,
		FileCatalog:   fileCatalog,
		CommandParams: cmdParams,
		Reloadable:    true,
	})
}

func datareadyFunc(ctx *Context) {
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
		errEnd(ctx)
	}
}

func runTaskFunc(ctx *Context) {
	fmt.Println("end start ...")

}

func procEndFunc(ctx *Context) {
	fmt.Println("end start ...")

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
}

func errEnd(ctx *Context) {
	ctx.AddQueue(&request.DataRequest{
		Rule:         "end",
		TransferType: request.NONETYPE,
		Priority:     1,
		Reloadable:   true,
	})
}
