package demanderonly

/**
    Author: luzequan
    Created: 2018-01-16 15:37:57
*/
import (
	. "drcs/core/databox"
	"fmt"
	"drcs/core/interaction/request"
	"drcs/runtime/status"
)

func init() {
	DEMCHILD.Register()
}

var DEMCHILD = &DataBox{
	Name:        "demchild",
	Description: "demchild",
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			fmt.Println("demchild start...")
			ctx.GetDataBox().StartWG.Done()
		},
		Trunk: map[string]*Rule{
			"child": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("child rule start ...")
					defer ctx.GetDataBox().SetStatus(status.STOP)
					defer ctx.GetDataBox().CloseRequestChan()

					parentBox := ctx.GetDataBox().ParentBox
					ctx.SetDataBox(parentBox).AddChanQueue(&request.DataRequest{
						Rule:         "parent",
						TransferType: request.NONETYPE,
						Priority:     1,
						Reloadable:   true,
					})
				},
			},
		},
	},
}
