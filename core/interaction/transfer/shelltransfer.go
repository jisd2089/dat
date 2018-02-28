package transfer

import (
	"drcs/core/interaction/response"
	"os/exec"
	"bytes"
	"fmt"
)

/**
    Author: luzequan
    Created: 2018-02-01 11:04:31
*/
type ShellTransfer struct {
}

func NewShellTransfer() Transfer {
	return &ShellTransfer{}
}

// 封装shell服务
func (st *ShellTransfer) ExecuteMethod(req Request) Response {

	//retCode := "000000"
	//cmd := exec.Command("split", "-l", "500", filePath, filePath + "_")
	cmd := exec.Command(req.GetCommandName(), req.GetCommandParams()...)

	fmt.Println("ShellTransfer exec: ", cmd.Args)

	//读取io.Writer类型的cmd.Stdout，再通x过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()

	if err != nil {
		fmt.Println("dem split file err: ", err)
		//retCode = "000001"
	}

	fmt.Println("dem split stdout: ", out.String())

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: "000000",
		//BodyStr:    out.String(),
	}
}

func (ft *ShellTransfer) Close() {

}
