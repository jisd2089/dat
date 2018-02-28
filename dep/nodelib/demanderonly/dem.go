package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"fmt"
	"strconv"
	"drcs/core/interaction/response"
)

func init() {
	DEM.Register()
}

var DEM = &DataBox{
	Name:        "demtest",
	Description: "demtest",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			fmt.Println(ctx)
			ctx.AddQueue(&request.DataRequest{
				Url:          "http://www.inderscience.com/info/inarticletoc.php?jcode=ijguc&year=2016&vol=7&issue=1",
				Rule:         "ruleTest",
				TransferType: request.HTTP,
			})
		},

		Trunk: map[string]*Rule{
			"ruleTest": {
				ParseFunc: func(ctx *Context) {
					fmt.Println("(((((((((((((((((")
					for i := 1; i < 10; i++ {
						ctx.AddQueue(&request.DataRequest{
							Url:          "http://www.inderscience.com/info/inarticletoc.php?jcode=ijguc&year=2016&vol=7&issue=" + strconv.Itoa(i),
							Rule:         "ruleTest2",
							TransferType: request.HTTP,
						})
					}
				},
			},
			"ruleTest2": {
				ParseFunc: func(ctx *Context) {
					fmt.Println(")))))))))))))))))))")
					//fmt.Println(string(ctx.DataResponse.GetBody()))
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


