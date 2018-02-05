package demanderonly

/**
    Author: luzequan
    Created: 2018-01-16 15:37:57
*/
import (
	//"path"
	"dat/core/interaction/request"
	. "dat/core/databox"
	"fmt"
	"sync"
	"dat/runtime/status"
)

func init() {
	DEMSUB.Register()
}

var DEMSUB = &DataBox{
	Name:        "demsub",
	Description: "demsub",
	IsParentBox: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			fmt.Println("dem parent start...")

			childBox := ctx.GetDataBox().GetChildBoxByName("demchild")
			childBoxAct := childBox.Copy()
			wg := &sync.WaitGroup{}
			childBoxAct.ParentBox = ctx.GetDataBox()
			childBoxAct.StartWG = wg

			wg.Add(1)
			ctx.GetDataBox().ChildBoxChan <- childBoxAct
			close(ctx.GetDataBox().ChildBoxChan)
			wg.Wait()

			ctx.SetDataBox(childBoxAct).AddChanQueue(&request.DataRequest{
				Rule:         "child",
				TransferType: request.NONETYPE,
				Priority:     1,
				Reloadable:   true,
			})
		},

		Trunk: map[string]*Rule{
			"parent": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("parent rule start ...")
					defer ctx.GetDataBox().SetStatus(status.STOP)
					defer ctx.GetDataBox().CloseRequestChan()
				},
			},
		},
	},
}