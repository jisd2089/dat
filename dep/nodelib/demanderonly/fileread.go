package demanderonly

/**
    Author: luzequan
    Created: 2018-02-02 13:27:36
*/
import (
	. "drcs/core/databox"
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
	"drcs/core/interaction/request"
	"drcs/runtime/status"
)

func init() {
	FILEREAD.Register()
}

var FILEREAD = &DataBox{
	Name:        "fileread",
	Description: "fileread",
	RuleTree: &RuleTree{
		Root: readFileRootFunc,

		Trunk: map[string]*Rule{
			"read": {
				ParseFunc: readFunc,
			},
		},
	},
}

func readFileRootFunc(ctx *Context) {
	fmt.Println("readFileRoot start ...")

	filePath := ctx.GetDataBox().GetDataFilePath()

	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	buf := bufio.NewReader(f)

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if err == io.EOF || err != nil {
			break
		}

		ctx.AddChanQueue(&request.DataRequest{
			Url:          "",
			Rule:         "read",
			TransferType: request.NONETYPE,
			Priority:     0,
			PostData:     line,
			Reloadable:   true,
		})
		break
	}
}

func readFunc(ctx *Context) {
	fmt.Println("read start...", ctx.DataRequest.PostData)

	defer ctx.GetDataBox().SetStatus(status.STOP)
	defer ctx.GetDataBox().CloseRequestChan()
}
