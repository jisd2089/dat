package main

/**
    Author: luzequan
    Created: 2018-08-13 10:03:09
*/
import (
	"drcs/core/interaction/request"
	. "drcs/core/databox"
	"drcs/runtime/status"
	logger "drcs/log"
)

var PLUGIN = &DataBox{
	Name:        "plugin_test",
	Description: "plugin_test",
	RuleTree: &RuleTree{
		Root: pluginRootFunc,

		Trunk: map[string]*Rule{
			"end": {
				ParseFunc: procEndFunc,
			},
		},
	},
}

func pluginRootFunc(ctx *Context) {
	logger.Info("pluginRootFunc start")

	ctx.AddChanQueue(&request.DataRequest{
		Rule:         "end",
		Method:       "GET",
		TransferType: request.NONETYPE,
		Reloadable:   true,
	})
}

func procEndFunc(ctx *Context) {
	logger.Info("end start")

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
}
