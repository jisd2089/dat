package demanderonly

/**
    Author: luzequan
    Created: 2018-01-16 15:37:57
*/
import (
	//"path"
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"fmt"
	"os"
	"path"
)

func init() {
	DEMSPLIT.Register()
}

var DEMSPLIT = &DataBox{
	Name:        "demsplit",
	Description: "demsplit",
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			fmt.Println("demsplit root start...")

			filePath := ctx.GetDataBox().DataFilePath

			dir, err := os.Open(path.Dir(filePath))

			if err != nil {
				fmt.Println("open file path error: ", err)
				return
			}

			names, err := dir.Readdirnames(0)
			if err != nil {
				fmt.Println("Readdirnames error: ", err)
				return
			}

			fmt.Println(names)

			// //"-l", "500", filePath, filePath + "_"
			//cmdParams := []string{"-l", "500", filePath, filePath + "_"}
			//
			//ctx.AddQueue(&request.DataRequest{
			//	Rule:          "split",
			//	TransferType:  request.SHELLTYPE,
			//	Priority:      1,
			//	CommandName:   "split",
			//	Reloadable:    true,
			//	CommandParams: cmdParams,
			//})
		},

		Trunk: map[string]*Rule{
			"split": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("split start ...")

					filePath := ctx.GetDataBox().DataFilePath

					cmdParams := []string{filePath + "_*"}

					ctx.AddQueue(&request.DataRequest{
						Rule:          "list",
						TransferType:  request.SHELLTYPE,
						Priority:      1,
						CommandName:   "ls",
						Reloadable:    true,
						CommandParams: cmdParams,
					})
				},
			},
			"list": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("list start ...")

					fmt.Println("list response: ", ctx.DataResponse.BodyStr)
				},
			},
		},
	},
}
