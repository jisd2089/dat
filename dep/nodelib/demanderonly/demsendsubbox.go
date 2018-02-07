package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"dat/core/interaction/request"
	. "dat/core/databox"
	"fmt"
)

func init() {
	DEMSENDSUBBOX.Register()
}

var DEMSENDSUBBOX = &DataBox{
	Name:         "demsendsubbox",
	Description:  "demsendsubbox",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			fmt.Println("demsendsubbox Root start...")

		},

		Trunk: map[string]*Rule{
			"runcollision": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("runcollision start ...")

					parentBox := ctx.GetDataBox().ParentBox
					ctx.SetDataBox(parentBox).AddChanQueue(&request.DataRequest{
						Url:          ctx.DataRequest.Url,
						Method:       "POST",
						Parameters:   ctx.DataRequest.Parameters,
						Rule:         "collisionrslt",
						TransferType: request.FASTHTTP,
						Priority:     0,
						Bobject:      ctx.DataRequest.Bobject,
						Reloadable:   true,
					})
				},
			},
			"pushToSup": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("pushToSup start ...")

				},
			},
		},
	},
}
