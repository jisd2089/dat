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
)

func init() {
	SUPREC.Register()
}

var SUPREC = &DataBox{
	Name:        "suprec",
	Description: "suprec",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			fmt.Println("suprec Root start...")
			//ctx.AddQueue(&request.DataRequest{
			//	Rule:         "process",
			//	TransferType: request.NONETYPE,
			//	Priority:     1,
			//	//Bobject:      paramBatch,
			//	Reloadable:   true,
			//})
		},

		Trunk: map[string]*Rule{
			"process": {
				SyncFunc: func(ctx *Context) *response.DataResponse {
					fmt.Println("process start ...")
					fmt.Println("obj: ", ctx.DataRequest.Bobject.(string))
					return &response.DataResponse{StatusCode: 200, ReturnCode: "000000", ReturnMsg: "成功"}
				},
			},
			"ruleTest2": {
				ParseFunc: func(ctx *Context) {
					fmt.Println(")))))))))))))))))))")
					//fmt.Println(string(ctx.DataResponse.GetBody()))
					ctx.AddQueue(&request.DataRequest{
						Rule:         "process",
						TransferType: request.NONETYPE,
						Priority:     1,
						//Bobject:      paramBatch,
						Reloadable:   true,
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
