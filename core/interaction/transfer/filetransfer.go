package transfer

import (
	"dat/core/interaction/response"
	"os"
	"io"
	cf "dat/common/compress"
	"path"
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

	switch req.GetMethod() {
	case "WRITE":
		writeFile(req)
	case "COMPRESS":
		compressFile(req)
	case "UNCOMPRESS":
		uncompressFile(req)
	}

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: "000000",
	}
}

func writeFile(req Request) {
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
}

func compressFile(req Request) {
	fileCat := req.GetFileCatalog()
	localFilePath := path.Join(fileCat.LocalDir, fileCat.LocalFileName)

	targetFilePath := path.Join(path.Join(fileCat.LocalDir, "compress"), fileCat.LocalFileName + ".tar.gz")

	localFile, err := os.OpenFile(localFilePath, os.O_RDWR, 0644)
	if err != nil {

	}

	files := make([]*os.File, 1)
	files[0] = localFile

	err = cf.Compress(files, targetFilePath)
	if err != nil {

	}

}

func uncompressFile(req Request) {

}

func (ft *FileTransfer) Close() {

}