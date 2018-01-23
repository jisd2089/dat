package transfer

import (
	"dat/core/interaction/response"
	"os"
	"io"
)

/**
    Author: luzequan
    Created: 2018-01-22 14:30:31
*/

type FileTransfer struct {}

func NewFileTransfer() Transfer {
	return &FileTransfer{}
}

// 封装fasthttp服务
func (ft *FileTransfer) ExecuteMethod(req Request) Response {

	targetFilePath := req.GetPostData()
	dataFile := req.GetDataFile()

	targetFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer targetFile.Close()
	if err != nil {

	}

	dataFileContent, err := dataFile.Open()
	defer dataFileContent.Close()
	if err != nil {

	}

	io.Copy(targetFile, dataFileContent)

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: "000000",
	}
}

func (ft *FileTransfer) Close() {

}