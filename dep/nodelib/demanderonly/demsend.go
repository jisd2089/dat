package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	//"os"
	"path"
	"dat/core/interaction/request"
	. "dat/core/databox"
	"fmt"
	"dat/core/interaction/response"
	"bufio"
	"strings"
	"io"
	"dat/dep/management/util"
	"dat/dep/management/entity"
	"os"
)

func init() {
	DEMSEND.Register()
}

var DEMSEND = &DataBox{
	Name:        "demsend",
	Description: "demsend",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			dataFile := path.Base(ctx.GetDataBox().GetDataFilePath())
			dataFileName := &util.DataFileName{}
			if err := dataFileName.ParseAndValidFileName(dataFile); err != nil {
			}

			paramBatch := entity.BatchReqestVo{}
			paramBatch.SerialNo = ""
			paramBatch.ReqType = ""
			paramBatch.OrderId = dataFileName.JobId
			paramBatch.FileNo = dataFileName.FileNo
			paramBatch.IdType = dataFileName.IdType
			paramBatch.TimeStamp = ""
			paramBatch.BatchNo = dataFileName.BatchNo
			paramBatch.UserId = ""
			paramBatch.Exid = ""
			paramBatch.TaskId = ""
			paramBatch.MaxDelay = ""
			paramBatch.Header = ""
			paramBatch.Exid = ""

			fmt.Println(ctx)

			fmt.Println(ctx.GetDataBox().GetDataFilePath())

			fmt.Println("NodeAddress: %s", ctx.GetDataBox().GetNodeAddress())
			ctx.AddQueue(&request.DataRequest{
				Rule:         "start",
				TransferType: request.NONETYPE,
				Priority:     1,
				Bobject:      paramBatch,
			})
		},

		Trunk: map[string]*Rule{
			"start": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("start ...")
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
							fmt.Println("file end ###############################")
							break
						}
						if err != nil {
							break
						}
						if rows == 0 { // 返回第一行头记录
							rows ++
							paramBatch.Header = line
							paramBatch.ReqType = entity.ReqType_Start
							for _, addr := range addressList {
								ctx.AddQueue(&request.DataRequest{
									Url:          addr.IP + addr.URL,
									Rule:         "process",
									TransferType: request.NONETYPE,
									Priority:     1,
									Bobject:      paramBatch,
									Reloadable:   true,
								})
							}
							break
						}
					}
				},
			},
			"process": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("process ...")
					addressList := ctx.GetDataBox().GetNodeAddress()

					continueFlag := true

					for _, addr := range addressList {
						uri := addr.IP + addr.URL
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
						fmt.Println("normaltype pending...")
						ctx.AddQueue(&request.DataRequest{
							Rule:         "normal",
							TransferType: request.NONETYPE,
							Priority:     1,
							Bobject:      ctx.DataRequest.Bobject,
							Reloadable:   true,
						})
					}

				},
			},
			"normal": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("normaltype ...")
					rows := 0
					paramBatch := ctx.DataRequest.Bobject.(entity.BatchReqestVo)
					addressList := ctx.GetDataBox().GetNodeAddress()
					url := addressList[0].IP + addressList[0].URL

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
							fmt.Println("file end ###############################", rows, line)
							break
						}
						if err != nil {
							break
						}
						rows ++
						if rows != 1 { // 跳过第一行头记录
							fmt.Println("rownum: ", rows, line)
							paramBatch.Exid = line
							paramBatch.ReqType = entity.ReqType_Normal
							ctx.AddQueue(&request.DataRequest{
								Url:          url,
								Rule:         "collision",
								TransferType: request.NONETYPE,
								Priority:     0,
								Bobject:      paramBatch,
								Reloadable:   true,
							})
						}
					}
				},
			},
			"collision": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("collision start.................")
					fmt.Println("detail count ", ctx.GetDataBox().DetailCount)

					//fmt.Println(ctx.DataRequest.Bobject.(entity.BatchReqestVo))
					addressList := ctx.GetDataBox().GetNodeAddress()
					currentUrl := ctx.DataRequest.GetUrl()

					fmt.Println(addressList)
					fmt.Println(currentUrl)
					fmt.Println(ctx.DataResponse.ReturnCode)

					if !strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
						for i, addr := range addressList {
							uri := addr.IP + addr.URL
							if strings.EqualFold(uri, currentUrl) {
								nextUrl := addressList[i+1].IP + addressList[i+1].URL
								ctx.AddQueue(&request.DataRequest{
									Url:          nextUrl,
									Rule:         "collision",
									TransferType: request.NONETYPE,
									Priority:     0,
									Bobject:      ctx.DataRequest.Bobject,
									Reloadable:   true,
								})
							}
						}
					} else {
						ctx.GetDataBox().TsfSuccCount ++
						fmt.Println("TsfSuccCount ", ctx.GetDataBox().TsfSuccCount)
					}

					if ctx.GetDataBox().TsfSuccCount == ctx.GetDataBox().DetailCount -1 {
						ctx.AddQueue(&request.DataRequest{
							//Url:          nextUrl,
							Rule:         "end",
							TransferType: request.NONETYPE,
							Priority:     0,
							Bobject:      ctx.DataRequest.Bobject,
						})
					}

					//ctx.GetDataBox().SyncProcess(ctx.DataRequest)
					//fmt.Println(string(ctx.DataResponse.GetBody()))
				},
			},
			"end": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("end start ...")

					//ctx.GetDataBox().SyncProcess(ctx.DataRequest)
					//fmt.Println(string(ctx.DataResponse.GetBody()))
				},
			},
			"collision1": {
				SyncFunc: func(ctx *Context) *response.DataResponse {
					fmt.Println(")))))))))))))))))))")

					return nil
				},
			},
		},
	},
}
