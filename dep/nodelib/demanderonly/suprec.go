package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"dat/core/interaction/request"
	. "dat/core/databox"
	"fmt"
	"dat/runtime/output"
	"dat/dep/management/entity"
	"strings"
	"dat/dep/management/constant"
)

func init() {
	SUPREC.Register()
}

var SUPREC = &DataBox{
	Name:        "suprec",
	Description: "suprec",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			fmt.Println("suprec Root start...")
			ctx.AddQueue(&request.DataRequest{
				Rule:         "process",
				TransferType: request.NONETYPE,
				Priority:     1,
				//Bobject:      paramBatch,
				Reloadable: true,
			})
		},

		Trunk: map[string]*Rule{
			"process": {
				ItemFields: []string{
					"FileName",
					"LocalDir",
					"TargetFolder",
					"WriteType",
					"Content",
				},
				ParseFunc: func(ctx *Context) {
					fmt.Println("process start ...")
					batchRequestVo := ctx.DataRequest.Bobject.(*entity.BatchReqestVo)
					fmt.Println("obj: ", batchRequestVo)
					writeType := output.CTW
					fileName := batchRequestVo.UserId + "_" + batchRequestVo.OrderId + "_" + batchRequestVo.IdType + "_" + batchRequestVo.BatchNo + "_" + batchRequestVo.FileNo + ".TARGET"
					content := ""

					// Redis碰撞
					ctx.ExecDataReq(&request.DataRequest{TransferType: request.REDIS,Rule: "process",})
					if ctx.DataResponse.StatusCode == 200 && strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
						switch batchRequestVo.ReqType {
						case constant.ReqType_Start:
							writeType = output.CTW
							content = batchRequestVo.Header + constant.LineTag
						case constant.ReqType_Normal:
							writeType = output.WA
							content = batchRequestVo.Exid + constant.LineTag
						}

						// 碰撞成功输出
						ctx.Output(map[string]interface{}{
							"FileName":     fileName,
							"LocalDir":     "D:/output",
							"TargetFolder": constant.TargetFolder,
							"WriteType":    writeType,
							"Content":      content,
						})
					}
				},
			},
		},
	},
}
