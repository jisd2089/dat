package demanderonly

/**
    Author: luzequan
    Created: 2017-12-28 17:22:00
*/
import (
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"fmt"
	"drcs/common/sftp"
	"drcs/runtime/output"
	"drcs/dep/management/entity"
	"drcs/dep/management/constant"
	"time"
	"drcs/core/interaction/response"
)

func init() {
	SUPREC.Register()
}

var SUPREC = &DataBox{
	Name:         "suprec",
	Description:  "suprec",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			//fmt.Println("suprec Root start...")
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
					//fmt.Println("process start ...")
					batchRequestVo := ctx.DataRequest.Bobject.(*entity.BatchReqestVo)
					//fmt.Println("obj: ", batchRequestVo)
					writeType := output.CTW
					fileName := batchRequestVo.UserId + "_" + batchRequestVo.OrderId + "_" + batchRequestVo.IdType + "_" + batchRequestVo.BatchNo + "_" + batchRequestVo.FileNo + ".TARGET"
					content := ""

					switch batchRequestVo.ReqType {
					case constant.ReqType_Start:
						writeType = output.CTW
						content = batchRequestVo.Header + constant.LineTag

						ctx.DataResponse = &response.DataResponse{
							StatusCode: 200,
							ReturnCode: "000000",
						}

						//fmt.Println("write content$$$$$$$$$: ", content)
						ctx.GetDataBox().ActiveWG.Add(1)
						// 碰撞成功输出
						ctx.Output(map[string]interface{}{
							"FileName":     fileName,
							"LocalDir":     "D:/dds_receive",
							//"LocalDir":     "/home/ddsdev/data/test/sup/rec",
							"TargetFolder": constant.TargetFolder,
							"WriteType":    writeType,
							"Content":      content,
						})
					case constant.ReqType_Normal:
						writeType = output.WA
						content = batchRequestVo.Exid + constant.LineTag

						ctx.ExecDataReq(&request.DataRequest{TransferType: request.REDIS, Rule: "process",})

						fmt.Println("write content$$$$$$$$$: ", content)
						ctx.GetDataBox().ActiveWG.Add(1)
						// 碰撞成功输出
						ctx.Output(map[string]interface{}{
							"FileName":     fileName,
							"LocalDir":     "D:/dds_receive",
							//"LocalDir":     "/home/ddsdev/data/test/sup/rec",
							"TargetFolder": constant.TargetFolder,
							"WriteType":    writeType,
							"Content":      content,
						})
						// Redis碰撞
						//ctx.ExecDataReq(&request.DataRequest{
						//	Method:       "EXIST",
						//	TransferType: request.REDIS,
						//	Rule:         "process",
						//	PostData:     batchRequestVo.Exid,
						//})
						//if ctx.DataResponse.StatusCode == 200 && strings.EqualFold(ctx.DataResponse.ReturnCode, "000000") {
						//
						//}
					case constant.ReqType_End:
						ctx.GetDataBox().ActiveWG.Wait()

						//fmt.Println("ReqType_End ^^^^^^^^^^^^^^^^^^^^^^")

						fileCatalog := &sftp.FileCatalog{
							UserName:       "bdaas",
							Password:       `bdaas`,
							Host:           "10.101.12.11",
							Port:           22,
							TimeOut:        10 * time.Second,
							LocalDir:       "D:/dds_receive/target",
							//LocalDir:       "/home/ddsdev/data/test/sup/rec/target",
							LocalFileName:  fileName,
							RemoteDir:      "/home/bdaas/data/test/sup/rec",
							RemoteFileName: fileName,
						}
						ctx.ExecDataReq(&request.DataRequest{
							Rule:         "pushToSup",
							TransferType: request.SFTP,
							Method:       "PUT",
							FileCatalog:  fileCatalog,
							Reloadable:   true,
						})

						if ctx.DataResponse.StatusCode == 200 && ctx.DataResponse.ReturnCode == "000000" {
							fmt.Println("StopActiveBox ~~~~~~~~~~~~", ctx.GetDataBox().PairDataBoxId)

							ctx.GetDataBox().StopActiveBox()
						}

						//go ctx.ExecDataReq(&request.DataRequest{
						//	TransferType: request.SFTP,
						//	Method:       "CLOSE",
						//	FileCatalog:  fileCatalog,
						//	Reloadable:   true,
						//})
					}
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
