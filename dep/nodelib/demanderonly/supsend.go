package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"dat/core/interaction/request"
	. "dat/core/databox"
	"fmt"
	"dat/core/interaction/response"
	"dat/common/sftp"
)

func init() {
	SUPSEND.Register()
}

var SUPSEND = &DataBox{
	Name:        "supsend",
	Description: "supsend",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			fmt.Println("supsend Root start...")

			// 1. 从sftp服务器（供方dmp服务器）拉取文件到节点服务器本地
			fileCatalog := &sftp.FileCatalog{}
			ctx.AddQueue(&request.DataRequest{
				FileCatalog:  fileCatalog,
				Rule:         "pushfile",
				TransferType: request.SFTP,
			})
		},

		Trunk: map[string]*Rule{
			"pushfile": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("pushfile start ...")
					// 2. 从本地节点服务器通过sftp方式推送至dem节点服务器
					fileCatalog := &sftp.FileCatalog{}
					ctx.AddQueue(&request.DataRequest{
						FileCatalog:  fileCatalog,
						Rule:         "notifydem",
						TransferType: request.SFTP,
					})
				},
			},
			"notifydem": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("notifydem start ...")
					// 3. 通知dem节点服务器，继续往下执行
					ctx.AddQueue(&request.DataRequest{
						Url:          "",
						Rule:         "notifydem",
						TransferType: request.HTTP,
					})
				},
			},
			"ruleTest3": {
				SyncFunc: func(ctx *Context) *response.DataResponse {
					fmt.Println(")))))))))))))))))))")

					//ctx.GetDataBox().SyncProcess(ctx.DataRequest)
					//fmt.Println(string(ctx.DataResponse.GetBody()))
					dResponse := &response.DataResponse{}
					dResponse.StatusCode = 303
					return dResponse
				},
			},
			"ruleTest4": {
				SyncFunc: func(ctx *Context) *response.DataResponse {
					fmt.Println(")))))))))))))))))))")

					return nil
				},
			},
		},
	},
}
