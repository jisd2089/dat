package demanderonly

/**
    Author: luzequan
    Created: 2018-01-16 15:37:57
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
	"dat/dep/management/entity"
	"os"
)

func init() {
	DEMSPLIT.Register()
}

var DEMSPLIT = &DataBox{
	Name:        "demsplit",
	Description: "demsplit",
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			dataFile := path.Base(ctx.GetDataBox().GetDataFilePath())

			f, err := os.Open(dataFile)
			defer f.Close()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			buf := bufio.NewReader(f)

			headerLine := ""
			rows := 0
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
					headerLine = line
				} else {

				}
			}
			fmt.Println("headerlint", headerLine)

			fmt.Println("NodeAddress: %s", ctx.GetDataBox().GetNodeAddress())
			ctx.AddQueue(&request.DataRequest{
				Rule:         "split",
				TransferType: request.DATABOX,
			})
		},

		Trunk: map[string]*Rule{
			"split": {
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
							paramBatch.DataBoxId = ctx.DataRequest.DataBoxId
							data, err := json.Marshal(paramBatch)
							if err != nil {
								break
							}
							for _, addr := range addressList {
								ctx.AddQueue(&request.DataRequest{
									Url:          addr.GetUrl(),
									Rule:         "process",
									TransferType: request.HTTP,
									Priority:     1,
									Bobject:      paramBatch,
									Reloadable:   true,
									Parameters:   data,
								})
							}
							break
						}
					}
				},
			},
		},
	},
}