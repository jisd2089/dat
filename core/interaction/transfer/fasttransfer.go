package transfer

import (
	"github.com/valyala/fasthttp"
	"sync"
	neturl "net/url"
	. "drcs/core/interaction/response"
	"time"
	"path"
	"os"
	"fmt"
	"io"
	"bufio"
	"mime/multipart"
	"bytes"
)

/**
    Author: luzequan
    Created: 2017-12-28 15:41:49
*/

type FastTransfer struct{}

func NewFastTransfer() Transfer {
	return &FastTransfer{}
}

// 封装fasthttp服务
func (ft *FastTransfer) ExecuteMethod(req Request) Response {
	//fmt.Println("execute fasthttp")
	//fmt.Println("fasthttp param:", string(req.GetParameters()))

	dataResponse := &DataResponse{}

	switch req.GetMethod() {
	case "POSTBODY":
		execPostByBody(req, dataResponse)
	case "POSTARGS":
		execPostByArgs(req, dataResponse)
	case "POSTFILE":
		execPostFile(req, dataResponse)
	case "FILESTREAM":
		execPostFileStream(req, dataResponse)
	}

	return dataResponse
}

func execPostByArgs(req Request, dataResponse *DataResponse) {

	freq := fasthttp.AcquireRequest()
	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(freq)
	defer fasthttp.ReleaseResponse(fresp)

	//freq.Header.SetContentType("application/json;charset=UTF-8")
	//freq.Header.SetMethod("POST")
	freq.SetRequestURI(req.GetUrl())

	freq.Header = *req.GetHeaderArgs()

	for k, v := range req.GetPostArgs() {
		freq.PostArgs().Set(k, v)
	}

	err := fasthttp.DoTimeout(freq, fresp, time.Duration(300) * time.Second)
	if err != nil {
		dataResponse.SetStatusCode(200)
		dataResponse.ReturnCode = "000009"
		return
	}
	//fmt.Println(string(fresp.Body()))

	dataResponse.SetStatusCode(200)
	dataResponse.ReturnCode = "000000"
	dataResponse.Body = fresp.Body()
}

func execPostByBody(req Request, dataResponse *DataResponse) {

	freq := fasthttp.AcquireRequest()
	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(freq)
	defer fasthttp.ReleaseResponse(fresp)

	freq.Header.SetContentType("application/json")
	freq.Header.SetMethod("POST")
	freq.SetRequestURI(req.GetUrl())
	freq.SetBody(req.GetParameters())

	err := fasthttp.DoTimeout(freq, fresp, time.Duration(300) * time.Second)
	if err != nil {
		dataResponse.SetStatusCode(200)
		dataResponse.ReturnCode = "000009"
		return
	}
	//fmt.Println(string(fresp.Body()))

	dataResponse.SetStatusCode(200)
	dataResponse.ReturnCode = "000000"
	dataResponse.Body = fresp.Body()

	//pipelineClient := getPipelineClient(req.GetUrl())
	//for {
	//	if err := pipelineClient.DoTimeout(freq, fresp, time.Duration(timeout)*time.Millisecond); err != nil {
	//		if err == fasthttp.ErrPipelineOverflow {
	//			//time.Sleep(1 * time.Millisecond)
	//			continue
	//		}
	//		dataResponse.SetStatusCode(400)
	//	}
	//	dataResponse.SetHeader(&fresp.Header)
	//	dataResponse.SetBody(fresp.Body())
	//	dataResponse.SetStatusCode(fresp.StatusCode())
	//	dataResponse.ReturnCode = "000000"
	//	break
	//}
}

func execPostFile(req Request, dataResponse *DataResponse) error {
	filePath := req.GetPostData()
	fileName := path.Base(filePath)
	targetUrl := req.GetUrl()

	timeOut := time.Duration(50) * time.Minute

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println("error writing to buffer")
		dataResponse.SetStatusCode(500)
		dataResponse.ReturnCode = "000000"
		return err
	}

	//打开文件句柄操作
	fh, err := os.Open(filePath)
	defer fh.Close()
	if err != nil {
		fmt.Println("error opening file")
		dataResponse.SetStatusCode(500)
		dataResponse.ReturnCode = "000000"
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		dataResponse.SetStatusCode(500)
		dataResponse.ReturnCode = "000000"
		return err
	}

	// 文件传输附带参数
	for _, p := range req.GetCommandParams() {
		bodyWriter.WriteField("dataFiles", p)
	}

	// 不定参数
	for _, k := range req.ParamKeys() {
		bodyWriter.WriteField(k, req.Param(k))
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	freq := fasthttp.AcquireRequest()
	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(freq)
	defer fasthttp.ReleaseResponse(fresp)

	freq.Header.SetContentType(contentType)
	freq.Header.SetMethod("POST")
	freq.SetRequestURI(targetUrl)
	freq.SetBody(bodyBuf.Bytes())

	err = fasthttp.DoTimeout(freq, fresp, timeOut)
	if err != nil {
		dataResponse.SetStatusCode(500)
		dataResponse.ReturnCode = "000000"
		return err
	}

	dataResponse.SetStatusCode(200)
	dataResponse.ReturnCode = "000000"
	return nil
}

func execPostFileStream(req Request, dataResponse *DataResponse) error {
	filePath := req.GetPostData()
	targetUrl := req.GetUrl()

	timeOut := time.Duration(50) * time.Minute

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//打开文件句柄操作
	fh, err := os.Open(filePath)
	defer fh.Close()
	if err != nil {
		fmt.Println("error opening file")
		dataResponse.SetStatusCode(500)
		dataResponse.ReturnCode = "000000"
		return err
	}

	// 文件传输附带参数
	//for _, p := range req.GetCommandParams() {
	//	bodyWriter.WriteField("dataFiles", p)
	//}

	// 不定参数
	//for _, k := range req.ParamKeys() {
	//	bodyWriter.WriteField(k, req.Param(k))
	//
	//}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	freq := fasthttp.AcquireRequest()
	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(freq)
	defer fasthttp.ReleaseResponse(fresp)

	freq.Header.SetContentType(contentType)
	freq.Header.SetMethod("POST")
	freq.SetRequestURI(targetUrl)
	//freq.SetBody(bodyBuf.Bytes())

	freq.Header.Set("dataFile", path.Base(filePath))

	// 不定参数
	for _, k := range req.ParamKeys() {
		freq.Header.Set(k, req.Param(k))
	}

	writeFlag := 1 // 0:buf 1:line

	freq.SetBodyStreamWriter(func(w *bufio.Writer) {
		buf := bufio.NewReader(fh)

		switch writeFlag {
		case 0:
			bufCnt := make([]byte, 1024)
			for {
				nr, err := buf.Read(bufCnt)
				if err == io.EOF || err != nil {
					//fmt.Println("file end ###############################")
					break
				}
				if nr > 0 {
					fmt.Fprintf(w, string(bufCnt))
					// Do not forget flushing streamed data to the client.
					if err := w.Flush(); err != nil {
						return
					}
				}
			}
		case 1:
			lineCnt := 100
			cntBuf := &bytes.Buffer{}
			c := lineCnt
			for {
				c--
				line, _, err := buf.ReadLine()

				if err != nil {
					if err == io.EOF && cntBuf.Len() == 0 {
						break
					} else if err == io.EOF && cntBuf.Len() > 0 {
						fmt.Fprintf(w, cntBuf.String()[:len(cntBuf.String())-1])
						// Do not forget flushing streamed data to the client.
						if err := w.Flush(); err != nil {
							return
						}
						break
					} else {
						return
					}
				}

				cntBuf.Write(line)
				cntBuf.WriteByte('\n')

				if c == 0 {

					fmt.Fprintf(w, cntBuf.String())
					// Do not forget flushing streamed data to the client.
					if err := w.Flush(); err != nil {
						return
					}

					c = lineCnt
					cntBuf.Reset()
				}
			}
		}
	})

	err = fasthttp.DoTimeout(freq, fresp, timeOut)
	if err != nil {
		dataResponse.SetStatusCode(500)
		dataResponse.ReturnCode = "000000"
		return err
	}

	dataResponse.SetStatusCode(200)
	dataResponse.ReturnCode = "000000"
	return nil
}

func (ft *FastTransfer) Close() {}

var (
	_clientMutex sync.RWMutex
	_clientMap   = make(map[string]*fasthttp.PipelineClient)
)

const (
	ClientMaxConns            = 500
	ClientMaxPendingRequests  = 1024
	ClientMaxIdleConnDuration = 10 * time.Second
)

func getClientFromMap(key string) *fasthttp.PipelineClient {
	_clientMutex.RLock()
	defer _clientMutex.RUnlock()
	return _clientMap[key]
}

func setClientToMap(key string, client *fasthttp.PipelineClient) {
	_clientMutex.Lock()
	defer _clientMutex.Unlock()
	if _clientMap[key] != nil {
		_clientMap[key] = client
	}
}

func newClientAndSetToMap(url string) *fasthttp.PipelineClient {
	_clientMutex.Lock()
	defer _clientMutex.Unlock()

	client := _clientMap[url]
	if client == nil {
		host := getHostFroURL(url)

		client = &fasthttp.PipelineClient{
			Addr:                host,
			MaxConns:            ClientMaxConns,
			MaxPendingRequests:  ClientMaxPendingRequests,
			MaxIdleConnDuration: ClientMaxIdleConnDuration,
		}

		_clientMap[url] = client
		_clientMap[host] = client
	}

	return client
}

func getHostFroURL(url string) string {
	u, _ := neturl.Parse(url)
	return u.Host
}

func getPipelineClient(url string) *fasthttp.PipelineClient {
	client := getClientFromMap(url)
	if client != nil {
		return client
	}

	host := getHostFroURL(url)
	client = getClientFromMap(host)
	if client != nil {
		setClientToMap(url, client)
		return client
	}

	client = newClientAndSetToMap(url)
	return client
}
