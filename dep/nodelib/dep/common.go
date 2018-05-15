package demanderonly

import (
	. "drcs/core/databox"

	"fmt"
	"drcs/runtime/status"
	"os"
	"drcs/core/interaction/request"
)

/**
    Author: luzequan
    Created: 2018-05-15 19:30:06
*/
func procEndFunc(ctx *Context) {
	fmt.Println("end start ...")

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
}

func errEnd(ctx *Context) {
	ctx.AddQueue(&request.DataRequest{
		Rule:         "end",
		TransferType: request.NONETYPE,
		Priority:     1,
		Reloadable:   true,
	})
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}