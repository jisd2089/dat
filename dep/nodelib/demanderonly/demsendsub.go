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
	"sync"
)

func init() {
	DEMSENDSUB.Register()
}

var DEMSENDSUB = &DataBox{
	Name:        "demsendsub",
	Description: "demsendsub",
	IsParentBox: true,
	RuleTree: &RuleTree{
		Root: demsendSubRootFunc,

		Trunk: map[string]*Rule{
			"start": {
				ParseFunc: startSubFunc,
			},
			"process": {
				ParseFunc: processSubFunc,
			},
			"normal": {
				ParseFunc: normalSubFunc,
			},
			"collisionrslt": {
				ParseFunc: collisionSubFunc,
			},
			"end": {
				ParseFunc: endSubFunc,
			},
			"endreslt": {
				ParseFunc: endResltSubFunc,
			},
		},
	},
}

func demsendSubRootFunc(ctx *Context) {
	fmt.Println("demsendsub Root ...")

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
		UserName: "bdaas",
		Password: `bdaas`,
		Host:     "10.101.12.11",
		Port:     22,
		TimeOut:  10 * time.Second,
		//LocalDir:       "/home/ddsdev/data/test/dem/send",
		LocalDir:       "D:/dds_send/tmp",
		LocalFileName:  dataFile,
		RemoteDir:      dataFilePath,
		RemoteFileName: dataFile,
	}

	fmt.Println("NodeAddress: %s", ctx.GetDataBox().GetNodeAddress())
	ctx.AddQueue(&request.DataRequest{
		Rule:         "start",
		Method:       "GET",
		TransferType: request.NONETYPE,
		FileCatalog:  fileCatalog,
		Bobject:      paramBatch,
		Reloadable:   true,
	})
}

func startSubFunc(ctx *Context) {
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

		if err == io.EOF || err != nil {
			//fmt.Println("file end ###############################")
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
					Method:       "POST",
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

func processSubFunc(ctx *Context) {
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

func normalSubFunc(ctx *Context) {
	fmt.Println("normaltype ...")

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("normaltype recover error: ", err)
			ctx.GetDataBox().CloseRequestChan()
		}
	}()

	childBox := ctx.GetDataBox().GetChildBoxByName("demsendsubbox")

	childBoxAct := childBox.Copy()
	wg := &sync.WaitGroup{}
	childBoxAct.ParentBox = ctx.GetDataBox()
	childBoxAct.StartWG = wg

	wg.Add(1)
	ctx.GetDataBox().ChildActiveBoxChan <- childBoxAct
	close(ctx.GetDataBox().ChildBoxChan)

	wg.Wait()

	rows := 0
	paramBatch := ctx.DataRequest.Bobject.(entity.BatchReqestVo)
	addressList := ctx.GetDataBox().GetNodeAddress()

	ctx.GetDataBox().ChildBox = childBoxAct

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

		if err == io.EOF || err != nil {

			ctx.GetDataBox().ParentBox.DetailCount = rows
			defer ctx.GetDataBox().CloseRequestChan()
			//fmt.Println("file end ###############################", rows, line)
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
			//ctx.AddChanQueue(&request.DataRequest{
			//	Url:          addressList[0].GetUrl(),
			//	Method:       "POST",
			//	Parameters:   data,
			//	Rule:         "collision",
			//	TransferType: request.FASTHTTP,
			//	Priority:     0,
			//	Bobject:      paramBatch,
			//	Reloadable:   true,
			//})
			ctx.SetDataBox(childBoxAct).AddChanQueue(&request.DataRequest{
				Url:          addressList[0].GetUrl(),
				Method:       "POST",
				Parameters:   data,
				Rule:         "runcollision",
				TransferType: request.NONETYPE,
				Priority:     1,
				Bobject:      paramBatch,
				Reloadable:   true,
			})
		}
	}
}

func collisionSubFunc(ctx *Context) {
	fmt.Println("collision result start.................")

	addressList := ctx.GetDataBox().GetNodeAddress()
	currentUrl := ctx.DataRequest.GetUrl()

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
						Rule:         "collisionrslt",
						Method:       "POST",
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
				"FileName": path.Base(ctx.GetDataBox().GetDataFilePath()) + ".SUCCESS",
				//"LocalDir":     "/home/ddsdev/data/test/dem/send",
				"LocalDir":     "D:/dds_send/success",
				"TargetFolder": constant.SuccessFolder,
				"WriteType":    output.CTWR,
				"Content":      paramBatch.Exid + "\n",
			})
		}
		fmt.Println("TsfSuccCount ", ctx.GetDataBox().TsfSuccCount)

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

func endSubFunc(ctx *Context) {
	fmt.Println("end start ...")

	ctx.GetDataBox().ChildBox.StopActiveBox()

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
			Method:       "POST",
			TransferType: request.FASTHTTP,
			Rule:         "endreslt",
			Priority:     1,
			Reloadable:   true,
			Parameters:   data,
		})
	}
}

func endResltSubFunc(ctx *Context) {
	fmt.Println("end reslt start ...")
	defer ctx.GetDataBox().SetStatus(status.STOP)

}
