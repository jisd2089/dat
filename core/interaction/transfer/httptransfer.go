package transfer

import (
	"fmt"
	"mime/multipart"
	"os"
	"net"
	"bytes"
	"io"
	"net/http"
	"dat/core/interaction/response"
	"path"
	"time"
)

/**
    Author: luzequan
    Created: 2018-01-21 16:27:49
*/
type HttpTransfer struct {}

func NewHttpTransfer() Transfer {
	return &HttpTransfer{}
}

// 封装fasthttp服务
func (ft *HttpTransfer) ExecuteMethod(req Request) Response {

	switch req.GetMethod() {
	case "Post":
		fmt.Println("post")
	case "POSTFILE":

		fileName := req.GetPostData()
		err := postFile(fileName, req.GetUrl())
		if err != nil {

		}
	}

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: "000000",
	}
}

func (ft *HttpTransfer) Close() {

}

func postFile(fileName string, targetUrl string) error {

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	//打开文件句柄操作
	filePath := path.Join("D:/output/SOURCE", fileName)
	fh, err := os.Open(filePath)
	defer fh.Close()
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest("POST", targetUrl, bytes.NewReader(bodyBuf.Bytes()))
	req.Header.Set("Content-Type", contentType)
	defaultClient := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(5 * time.Hour)
				c, err := net.DialTimeout(netw, addr, 5*time.Hour)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	// 执行
	fresp, err := defaultClient.Do(req)
	if err != nil {
		return err
	}

	fmt.Println("fresp :", fresp)

	//resp, err := http.Post(targetUrl, contentType, bodyBuf)
	//defer resp.Body.Close()
	//if err != nil {
	//	return err
	//}
	//resp_body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(resp.Status)
	//fmt.Println(string(resp_body))

	return nil
}