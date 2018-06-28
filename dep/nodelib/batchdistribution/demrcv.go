package batchdistribution

/**
    Author: luzequan
    Created: 2018-06-25 18:08:48
*/
import (
	//"path"
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	//"drcs/common/sftp"
	//"drcs/runtime/status"
	"fmt"
	//"strings"
	//"time"
)

func init() {
	BATCH_DEM_RCV.Register()
}

var BATCH_DEM_RCV = &DataBox{
	Name:        "batch_dem_rcv",
	Description: "batch_dem_rcv",
	RuleTree: &RuleTree{
		Root: batchDemRcvRootFunc,

		Trunk: map[string]*Rule{
			"checkMD5": {
				ParseFunc: checkMD5Func,
			},
			"pushToServer": {
				ParseFunc: pushToServerFunc,
			},
			"sendRecord": {
				ParseFunc: sendRcvRecordFunc,
			},
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func batchDemRcvRootFunc(ctx *Context) {
	fmt.Println("batchDemRcvRootFunc ...")

	ctx.AddQueue(&request.DataRequest{
		Rule:         "checkMD5",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}