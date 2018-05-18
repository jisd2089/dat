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
	"drcs/common/sftp"
	"time"
	"strings"
)

func init() {
	ALGRCV.Register()
}

var relyDatas map[string]int

var ALGRCV = &DataBox{
	Name:        "algorithmreceive",
	Description: "algorithmreceive",
	RuleTree: &RuleTree{
		Root: algrcvRootFunc,

		Trunk: map[string]*Rule{
			"pushtoserver": {
				ParseFunc: pushToServerFunc, // 将算法产品推送到hadoop server
			},
			"datacheck": {
				ParseFunc: dataCheckFunc, // 检查hdfs数据文件是否准备就绪
			},
			"dataready": {
				ParseFunc: datareadyFunc, // hdfs数据文件准备就绪
			},
			"runtask": {
				ParseFunc: runTaskFunc, // 执行算法任务
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func algrcvRootFunc(ctx *Context) {
	fmt.Println("algorithmreceive Root ...")

	relyDatas = make(map[string]int)

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

	fmt.Println(dataFilePath + "&&" + ctx.GetDataBox().GetDataFilePath())

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

	// 算法文件路径重置为hadoop server上的路径
	ctx.GetDataBox().SetDataFilePath(path.Join(fsAddress.RemoteDir, dataFile))

	ctx.AddQueue(&request.DataRequest{
		Rule:         "datacheck",
		Method:       "PUT",
		TransferType: request.SFTP,
		FileCatalog:  fileCatalog,
		Reloadable:   true,
	})
}

func dataCheckFunc(ctx *Context) {
	fmt.Println("algorithmreceive data check ...")

	if ! (ctx.DataResponse.StatusCode == 200 && strings.EqualFold(ctx.DataResponse.ReturnCode, "000000")) {
		errEnd(ctx)
		return
	}

	fsAddress := ctx.GetDataBox().FileServerAddress

	fileCatalog := &sftp.FileCatalog{
		UserName: fsAddress.UserName,
		Password: fsAddress.Password,
		Host:     fsAddress.Host,
		Port:     fsAddress.Port,
		TimeOut:  time.Duration(fsAddress.TimeOut) * time.Second,
	}

	hdfsPaths := ctx.GetDataBox().Params

	ctx.GetDataBox().TsfSuccCount = len(hdfsPaths)

	for _, p := range hdfsPaths {

		cmdParam := `/usr/local/package/hadoop-2.7.3/bin/hdfs dfs -test -e ` + p + `
if [ $? -eq 0 ] ;then
    echo 'exist'
else
    echo 'Error! Directory is not exist'
fi
`

		ctx.AddQueue(&request.DataRequest{
			Rule:         "dataready",
			Method:       "STRING",
			TransferType: request.SSH,
			FileCatalog:  fileCatalog,
			CommandName:  cmdParam,
			Reloadable:   true,
		})
	}
}

func datareadyFunc(ctx *Context) {
	fmt.Println("algorithmreceive data ready...")

	cmdParam := ctx.GetResponse().BodyStr

	if (ctx.GetResponse().StatusCode == 200 && strings.EqualFold(strings.TrimSpace(string(ctx.GetResponse().Body)), "exist")) {

		relyDatas[cmdParam] = 1

	} else {

		time.Sleep(time.Duration(60) * time.Second)

		fsAddress := ctx.GetDataBox().FileServerAddress

		fileCatalog := &sftp.FileCatalog{
			UserName: fsAddress.UserName,
			Password: fsAddress.Password,
			Host:     fsAddress.Host,
			Port:     fsAddress.Port,
			TimeOut:  time.Duration(fsAddress.TimeOut) * time.Second,
		}

		ctx.AddQueue(&request.DataRequest{
			Rule:         "dataready",
			TransferType: request.NONETYPE,
			CommandName:  cmdParam,
			FileCatalog:  fileCatalog,
			Reloadable:   true,
		})
		return
	}

	if len(relyDatas) == ctx.GetDataBox().TsfSuccCount {
		ctx.AddQueue(&request.DataRequest{
			Rule:         "runtask",
			TransferType: request.NONETYPE,
			Reloadable:   true,
		})
	}
}

func runTaskFunc(ctx *Context) {
	fmt.Println("algorithmreceive run task...")

	inputPaths := ctx.GetDataBox().Params
	outputPath := ctx.GetDataBox().Param("hdfsOutputDir")

	cmdParams := []string{}
	cmdParams = append(cmdParams, "/usr/local/package/spark-2.2.1-bin-hadoop2.7/bin/spark-submit")
	cmdParams = append(cmdParams, "--class com.chinadep.spark.LogisticRegression")
	cmdParams = append(cmdParams, "--master yarn")
	cmdParams = append(cmdParams, "--deploy-mode cluster")
	cmdParams = append(cmdParams, "--queue sparkqueue")
	cmdParams = append(cmdParams, ctx.GetDataBox().GetDataFilePath())

	for _, i := range inputPaths {
		cmdParams = append(cmdParams, i)
	}

	cmdParams = append(cmdParams, outputPath)

	fmt.Println(cmdParams)

	fsAddress := ctx.GetDataBox().FileServerAddress

	// 1. 从sftp服务器（需方dmp服务器）拉取文件到节点服务器本地
	fileCatalog := &sftp.FileCatalog{
		UserName: fsAddress.UserName,
		Password: fsAddress.Password,
		Host:     fsAddress.Host,
		Port:     fsAddress.Port,
		TimeOut:  time.Duration(fsAddress.TimeOut) * time.Second,
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
