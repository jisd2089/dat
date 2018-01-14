package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"dat/core/interaction/request"
	. "dat/core/databox"
	"fmt"
	"dat/common/sftp"
)

func init() {
	DEMREC.Register()
}

var DEMREC = &DataBox{
	Name:        "demrec",
	Description: "demrec",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			fmt.Println("demrec Root start ...")
			// 1. 校验md5
			ctx.AddQueue(&request.DataRequest{
				Rule:         "verify",
				TransferType: request.NONETYPE,
			})
		},

		Trunk: map[string]*Rule{
			"verify": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("demrec verify start ...")
					// 2. 将接收到的反馈文件推送至需方dmp
					ctx.AddQueue(&request.DataRequest{
						FileCatalog:  &sftp.FileCatalog{},
						Rule:         "pushdem",
						TransferType: request.SFTP,
					})
				},
			},
			"pushdem": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("demrec pushdem start ...")
					// 3. 记录业务日志

				},
			},
		},
	},
}
