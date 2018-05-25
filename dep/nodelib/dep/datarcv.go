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
	DATARCV.Register()
}

var DATARCV = &DataBox{
	Name:        "datareceive",
	Description: "datareceive",
	RuleTree: &RuleTree{
		Root: datarcvRootFunc,

		Trunk: map[string]*Rule{
			"existhdfs": {
				ParseFunc: existHdfsFileFunc,
			},
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
		Rule:         "existhdfs",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func existHdfsFileFunc(ctx *Context) {
	fmt.Println("datareceive push data to server...")

	hdfsInputDir := ctx.GetDataBox().Param("hdfsInputDir")

	filePath := ctx.GetDataBox().GetDataFilePath()
	dataFile := path.Base(filePath)

	hdfsPath := path.Join(hdfsInputDir, dataFile)

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

	cmdParam := `/usr/local/package/hadoop-2.7.3/bin/hdfs dfs -test -e ` + hdfsPath + `
if [ $? -eq 0 ] ;then
    echo 'exist'
else
    echo 'Error! Directory is not exist'
fi
`

	ctx.AddQueue(&request.DataRequest{
		Rule:         "pushtoserver",
		Method:       "STRING",
		TransferType: request.SSH,
		FileCatalog:  fileCatalog,
		CommandName:  cmdParam,
		Reloadable:   true,
	})
}

func pushDataToServerFunc(ctx *Context) {
	fmt.Println("datareceive push data to server...")

	if ctx.GetResponse().StatusCode == 200 && strings.EqualFold(strings.TrimSpace(string(ctx.GetResponse().Body)), "exist") {
		errEnd(ctx)
		return
	}

	filePath := ctx.GetDataBox().GetDataFilePath()
	dataFile := path.Base(filePath)
	dataFilePath := path.Dir(filePath)

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

	hdfsInputDir := ctx.GetDataBox().Param("hdfsInputDir")

	cmdParams := []string{}
	cmdParams = append(cmdParams, "/usr/local/package/hadoop-2.7.3/bin/hdfs")
	cmdParams = append(cmdParams, "dfs")
	cmdParams = append(cmdParams, "-put")
	cmdParams = append(cmdParams, filePath)
	cmdParams = append(cmdParams, hdfsInputDir)

	ctx.AddQueue(&request.DataRequest{
		Rule:          "end",
		Method:        "SLICE",
		TransferType:  request.SSH,
		FileCatalog:   fileCatalog,
		CommandParams: cmdParams,
		Reloadable:    true,
	})
}
