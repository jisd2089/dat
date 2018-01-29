package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"path"
	"dat/core/interaction/request"
	. "dat/core/databox"
	"fmt"
	"encoding/json"
	"bufio"
	"strings"
	"io"
	"dat/dep/management/constant"
	"dat/dep/management/util"
	"dat/dep/management/entity"
	"os"
	"strconv"
	"dat/common/sftp"
	"dat/runtime/output"
	"dat/runtime/status"
	"time"
)

func init() {
	DEMSEND.Register()
}

var DEMSEND = &DataBox{
	Name:        "demsend",
	Description: "demsend",
	RuleTree: &RuleTree{
		Root: demsendRootFunc,

		Trunk: map[string]*Rule{
			"start": {
				ParseFunc: startFunc,
			},
			"process": {
				ParseFunc: processFunc,
			},
			"normal": {
				ParseFunc: normalFunc,
			},
			"collision": {
				ParseFunc: collisionFunc,
			},
			"end": {
				ParseFunc: endFunc,
			},
			"endreslt": {
				ParseFunc: endresltFunc,
			},
		},
	},
}

func demsendRootFunc(ctx *Context) {
	fmt.Println("demsend Root ...")

	memId := "000001"
	filePath := ctx.GetDataBox().GetDataFilePath()
	dataFile := path.Base(filePath)
	dataFilePath := path.Dir(filePath)
	dataFileName := &util.DataFileName{}
	if err := dataFileName.ParseAndValidFileName(dataFile); err != nil {
		return
	}

	paramBatch := entity.BatchReqestVo{}
	paramBatch.SerialNo = util.NewSeqUtil().GenBusiSerialNo(memId)
	paramBatch.OrderId = dataFileName.JobId
	paramBatch.FileNo = dataFileName.FileNo
	paramBatch.IdType = dataFileName.IdType
	paramBatch.TimeStamp = util.GetTimestampString()
	paramBatch.BatchNo = dataFileName.BatchNo
	paramBatch.UserId = memId
	paramBatch.MaxDelay = strconv.Itoa(10)
	paramBatch.ReqType = ""
	paramBatch.TaskId = ""
	paramBatch.Header = ""
	paramBatch.Exid = ""

	fmt.Println(ctx.GetDataBox().GetDataFilePath())

	// 1. 从sftp服务器（需方dmp服务器）拉取文件到节点服务器本地
	fileCatalog := &sftp.FileCatalog{
		UserName:       "bdaas",
		Password:       `bdaas`,
		Host:           "10.101.12.11",
		Port:           22,
		TimeOut:        10 * time.Second,
		LocalDir:       "/home/ddsdev/data/test/dem/send",
		LocalFileName:  dataFile,
		RemoteDir:      dataFilePath,
		RemoteFileName: dataFile,
	}

	fmt.Println("NodeAddress: %s", ctx.GetDataBox().GetNodeAddress())
	ctx.AddQueue(&request.DataRequest{
		Rule:         "start",
		Method:       "GET",
		TransferType: request.SFTP,
		FileCatalog:  fileCatalog,
		Bobject:      paramBatch,
		Reloadable:   true,
	})
}

func startFunc(ctx *Context) {
	fmt.Println("start rule...")
	rows := 0
	paramBatch := ctx.DataRequest.Bobject.(entity.BatchReqestVo)
	addressList := ctx.GetDataBox().GetNodeAddress()

	dataFilePath := path.Join(ctx.DataRequest.FileCatalog.LocalDir, ctx.DataRequest.FileCatalog.LocalFileName)
	ctx.GetDataBox().DataFilePath = dataFilePath

	f, err := os.Open(dataFilePath)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	buf := bufio.NewReader(f)

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if err == io.EOF {
			//fmt.Println("file end ###############################")
			break
		}
		if err != nil {
			break
		}
		if rows == 0 { // 返回第一行头记录
			rows ++
			paramBatch.Header = line
			paramBatch.ReqType = entity.ReqType_Start
			paramBatch.DataBoxId = ctx.DataRequest.DataBoxId
			data, err := json.Marshal(paramBatch)
			if err != nil {
				break
			}
			for _, addr := range addressList {
				ctx.AddQueue(&request.DataRequest{
					Url:          addr.GetUrl(),
					Rule:         "process",
					TransferType: request.FASTHTTP,
					Priority:     1,
					Bobject:      paramBatch,
					Reloadable:   true,
					Parameters:   data,
				})
			}
			break
		}
	}
}

func processFunc(ctx *Context) {
	fmt.Println("process ...")
	addressList := ctx.GetDataBox().GetNodeAddress()

	continueFlag := true

	for _, addr := range addressList {
		uri := addr.GetUrl()
		if strings.EqualFold(uri, ctx.DataRequest.Url) && ctx.DataResponse.StatusCode == 200 {
			addr.Connectable = true
			fmt.Println("uri: %s, connectable: true", uri)
		} else if strings.EqualFold(uri, ctx.DataRequest.Url) && ctx.DataResponse.StatusCode != 200 {
			addr.Connectable = false
			fmt.Println("uri: %s, connectable: false", uri)
		}

		if !addr.Connectable && addr.RetryTimes == 0 {
			continueFlag = false
		}
	}

	if continueFlag {
		//fmt.Println("normaltype pending...")
		ctx.AddQueue(&request.DataRequest{
			Rule:         "normal",
			TransferType: request.NONETYPE,
			Priority:     1,
			Bobject:      ctx.DataRequest.Bobject,
			Reloadable:   true,
		})
	}
}

func normalFunc(ctx *Context) {
	fmt.Println("normaltype ...")

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("normaltype recover error: ", err)
			ctx.GetDataBox().CloseRequestChan()
		}
	}()
	rows := 0
	paramBatch := ctx.DataRequest.Bobject.(entity.BatchReqestVo)
	addressList := ctx.GetDataBox().GetNodeAddress()

	f, err := os.Open(ctx.GetDataBox().GetDataFilePath())
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	buf := bufio.NewReader(f)

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if err == io.EOF {
			ctx.GetDataBox().DetailCount = rows
			defer ctx.GetDataBox().CloseRequestChan()
			//fmt.Println("file end ###############################", rows, line)
			break
		}
		if err != nil {
			break
		}
		rows ++
		if rows != 1 { // 跳过第一行头记录
			//fmt.Println("rownum: ", rows, line)
			paramBatch.Exid = line
			paramBatch.ReqType = entity.ReqType_Normal
			paramBatch.DataBoxId = ctx.DataRequest.DataBoxId
			data, err := json.Marshal(paramBatch)
			if err != nil {
				break
			}
			ctx.AddChanQueue(&request.DataRequest{
				Url:          addressList[0].GetUrl(),
				Parameters:   data,
				Rule:         "collision",
				TransferType: request.FASTHTTP,
				Priority:     0,
				Bobject:      paramBatch,
				Reloadable:   true,
			})
		}
	}
}

func collisionFunc(ctx *Context) {
	//fmt.Println("collision start.................")
	//fmt.Println("detail count ", ctx.GetDataBox().DetailCount)
	//fmt.Println("collision response.................", ctx.DataResponse.StatusCode, ctx.DataResponse.ReturnCode)

	addressList := ctx.GetDataBox().GetNodeAddress()
	currentUrl := ctx.DataRequest.GetUrl()

	//fmt.Println(addressList)
	//fmt.Println(currentUrl)
	//fmt.Println(ctx.DataResponse.ReturnCode)

	if ctx.DataResponse.StatusCode == 200 {

		if !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
			for i, addr := range addressList {
				if strings.EqualFold(addr.GetUrl(), currentUrl) {
					if i+1 >= len(addressList) { // no hit
						ctx.GetDataBox().ExecTsfSuccCount()
						break
					}
					if !addressList[i+1].Connectable {
						continue
					}
					nextUrl := addressList[i+1].GetUrl()
					ctx.AddQueue(&request.DataRequest{
						Url:          nextUrl,
						Rule:         "collision",
						TransferType: request.FASTHTTP,
						Priority:     0,
						Bobject:      ctx.DataRequest.Bobject,
						Reloadable:   true,
						Parameters:   ctx.DataRequest.Parameters,
					})
				}
			}
		} else {
			//ctx.GetDataBox().TsfSuccCount ++
			ctx.GetDataBox().ExecTsfSuccCount()

			paramBatch := ctx.DataRequest.Bobject.(entity.BatchReqestVo)

			ctx.Output(map[string]interface{}{
				"FileName":     path.Base(ctx.GetDataBox().GetDataFilePath()) + ".SUCCESS",
				"LocalDir":     "/home/ddsdev/data/test/dem/send",
				"TargetFolder": constant.SuccessFolder,
				"WriteType":    output.CTWR,
				"Content":      paramBatch.Exid + "\n",
			})
		}
		//fmt.Println("TsfSuccCount ", ctx.GetDataBox().TsfSuccCount)

		if ctx.GetDataBox().TsfSuccCount == ctx.GetDataBox().DetailCount-1 {

			fmt.Println("dem send end ************************", ctx.GetDataBox().GetId())

			ctx.AddQueue(&request.DataRequest{
				Rule:         "end",
				TransferType: request.NONETYPE,
				Priority:     0,
				Bobject:      ctx.DataRequest.Bobject,
				Reloadable:   true,
			})
		}
	}
}

func endFunc(ctx *Context) {
	fmt.Println("end start ...")

	defer ctx.GetDataBox().SetStatus(status.STOP)

	paramBatch := ctx.DataRequest.Bobject.(entity.BatchReqestVo)
	addressList := ctx.GetDataBox().GetNodeAddress()

	paramBatch.ReqType = entity.ReqType_End
	paramBatch.DataBoxId = ctx.DataRequest.DataBoxId
	data, err := json.Marshal(paramBatch)
	if err != nil {
		return
	}
	for _, addr := range addressList {
		ctx.AddQueue(&request.DataRequest{
			Url:          addr.GetUrl(),
			TransferType: request.FASTHTTP,
			Rule:         "endreslt",
			Priority:     1,
			Reloadable:   true,
			Parameters:   data,
		})
	}
}

func endresltFunc(ctx *Context) {
	fmt.Println("end reslt start ...")

}
