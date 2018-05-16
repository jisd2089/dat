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
	"time"
	"strings"
	"os/exec"
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
		Rule:         "dataready",
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

	hdfsInputDir := ctx.GetDataBox().Param("hdfsInputDir")

	cmdParams := []string{}
	cmdParams = append(cmdParams, "/usr/local/package/hadoop-2.7.3/bin/hdfs")
	cmdParams = append(cmdParams, "dfs")
	cmdParams = append(cmdParams, "-put")
	cmdParams = append(cmdParams, filePath)
	cmdParams = append(cmdParams, hdfsInputDir)

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
	fmt.Println("algorithmreceive data ready ...")

	if ctx.DataResponse.StatusCode == 200 && strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {

		hdfsPath := ctx.GetDataBox().Param("dataPath")

		//hdfsInputDir := ctx.GetDataBox().Param("hdfsInputDir")
		//
		//filePath := ctx.GetDataBox().GetDataFilePath()
		//dataFile := path.Base(filePath)
		//
		//hdfsPath := path.Join(hdfsInputDir, dataFile)

		fsAddress := ctx.GetDataBox().FileServerAddress

		// 1. push local file to hadoop server
		fileCatalog := &sftp.FileCatalog{
			UserName:       fsAddress.UserName,
			Password:       fsAddress.Password,
			Host:           fsAddress.Host,
			Port:           fsAddress.Port,
			TimeOut:        time.Duration(fsAddress.TimeOut) * time.Second,
		}

		cmdParam := `/usr/local/package/hadoop-2.7.3/bin/hdfs dfs -test -e ` + hdfsPath + `
if [ $? -eq 0 ] ;then
    echo 'exist'
else
    echo 'Error! Directory is not exist'
fi
`

		ctx.AddQueue(&request.DataRequest{
			Rule:         "runtask",
			Method:       "STRING",
			TransferType: request.SSH,
			FileCatalog:  fileCatalog,
			CommandName:  cmdParam,
			Reloadable:   true,
		})

	} else {
		errEnd(ctx)
	}
}

func runTaskFunc(ctx *Context) {
	fmt.Println("algorithmreceive run task...")

	if !(ctx.GetResponse().StatusCode == 200 && strings.EqualFold(strings.TrimSpace(string(ctx.GetResponse().Body)), "exist")) {

		time.Sleep(time.Duration(60) * time.Second)

		ctx.AddQueue(&request.DataRequest{
			Rule:          "dataready",
			TransferType:  request.NONETYPE,
			Reloadable:    true,
		})
		return
	}

	cmdParams := []string{}
	cmdParams = append(cmdParams, "/usr/local/package/spark-2.2.1-bin-hadoop2.7/bin/spark-submit")
	cmdParams = append(cmdParams, "--class chinadep.Precollision")
	cmdParams = append(cmdParams, "--master yarn")
	cmdParams = append(cmdParams, "--deploy-mode cluster")
	cmdParams = append(cmdParams, "--queue sparkqueue")
	cmdParams = append(cmdParams, "--jars hdfs://deptest20:9000/user/aarontest/aaron-oozie/shell-spark-Precollision/lib/GenExidRealTime.jar")

	fsAddress := ctx.GetDataBox().FileServerAddress

	// 1. 从sftp服务器（需方dmp服务器）拉取文件到节点服务器本地
	fileCatalog := &sftp.FileCatalog{
		UserName:       fsAddress.UserName,
		Password:       fsAddress.Password,
		Host:           fsAddress.Host,
		Port:           fsAddress.Port,
		TimeOut:        time.Duration(fsAddress.TimeOut) * time.Second,
	}

	ctx.AddQueue(&request.DataRequest{
		Rule:          "end",
		Method:        "SLICE",
		TransferType:  request.SSH,
		FileCatalog:   fileCatalog,
		CommandParams: cmdParams,
		Reloadable:    true,
	})
}


func isHdfsDirExists(ctx *Context) bool {

	hdfsInputDir := ctx.GetDataBox().Param("hdfsInputDir")

	filePath := ctx.GetDataBox().GetDataFilePath()
	dataFile := path.Base(filePath)

	hdfsPath := path.Join(hdfsInputDir, dataFile)

	c := `/usr/local/package/hadoop-2.7.3/bin/hdfs dfs -test -e ` + hdfsPath + `
if [ $? -eq 0 ] ;then
    echo 'exist'
else
    echo 'Error! Directory is not exist'
fi
`
	cmd := exec.Command("sh", "-c", c)
	out, err := cmd.Output()
	if err != nil || !strings.EqualFold(string(out), "exist") {
		return false
	}

	return true
}


