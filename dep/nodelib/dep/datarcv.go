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
	"encoding/json"
	"bufio"
	"strings"
	"io"
	"drcs/dep/management/util"
	"drcs/dep/management/entity"
	"os"
	"strconv"
	"drcs/common/sftp"
	"time"
)

func init() {
	DATARCV.Register()
}

var DATARCV = &DataBox{
	Name:        "datareceive",
	Description: "datareceive",
	RuleTree: &RuleTree{
		Root: datarcvRootFunc,

		Trunk: map[string]*Rule{
			"pushtoserver": {
				ParseFunc: pushDataToServerFunc,
			},
			"puttohdfs": {
				ParseFunc: putDataToHDFSFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func datarcvRootFunc(ctx *Context) {
	fmt.Println("datareceive Root...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "pushtoserver",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func pushDataToServerFunc(ctx *Context) {
	fmt.Println("datareceive push data to server...")
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

func putDataToHDFSFunc(ctx *Context) {
	fmt.Println("datareceive put data to hdfs...")

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


