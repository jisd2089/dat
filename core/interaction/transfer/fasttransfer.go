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
	case "POST":
		execPost(req, dataResponse)
	case "POSTFILE":
		execPostFileStream(req, dataResponse)
	case "FILESTREAM":
		execPostFileStream(req, dataResponse)
	}

	return dataResponse
}

func execPost(req Request, dataResponse *DataResponse) {

	timeout := 30 * 1000

	freq := fasthttp.AcquireRequest()
	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(freq)
	defer fasthttp.ReleaseResponse(fresp)

	freq.SetRequestURI(req.GetUrl())
	freq.Header.SetMethod("POST")
	freq.SetBody(req.GetParameters())

	pipelineClient := getPipelineClient(req.GetUrl())
	for {
		if err := pipelineClient.DoTimeout(freq, fresp, time.Duration(timeout)*time.Millisecond); err != nil {
			if err == fasthttp.ErrPipelineOverflow {
				//time.Sleep(1 * time.Millisecond)
				continue
			}
			dataResponse.SetStatusCode(400)
		}
		dataResponse.SetHeader(&fresp.Header)
		dataResponse.SetBody(fresp.Body())
		dataResponse.SetStatusCode(fresp.StatusCode())
		dataResponse.ReturnCode = "000000"
		break
	}
}

func execPostFile(req Request, dataResponse *DataResponse) error {
	fileName := req.GetPostData()
	targetUrl := req.GetUrl()

	timeOut := time.Duration(50) * time.Minute

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	//打开文件句柄操作
	//filePath := path.Join("/home/ddsdev/data/test/sup/send", fileName)
	filePath := path.Join("D:/dds_send/tmp", fileName)

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
		return err
	}

	dataResponse.SetStatusCode(200)
	dataResponse.ReturnCode = "000000"
	return nil
}

func execPostFileStream(req Request, dataResponse *DataResponse) error {
	fileName := req.GetPostData()
	targetUrl := req.GetUrl()

	timeOut := time.Duration(50) * time.Minute

	//打开文件句柄操作
	//filePath := path.Join("/home/ddsdev/data/test/sup/send", fileName)
	filePath := path.Join("D:/dds_send/tmp", fileName)

	fh, err := os.Open(filePath)
	defer fh.Close()
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	freq := fasthttp.AcquireRequest()
	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(freq)
	defer fasthttp.ReleaseResponse(fresp)

	freq.Header.SetMethod("POST")
	freq.SetRequestURI(targetUrl)

	freq.SetBodyStreamWriter(func(w *bufio.Writer) {
		buf := bufio.NewReader(fh)
		//rows := 0
		//lineCnt := ""

		bufCnt := make([]byte, 104857600)
		//
		//buf.Read(bufCnt)
		//fmt.Fprintf(w, string(bufCnt))
		//// Do not forget flushing streamed data to the client.
		//if err := w.Flush(); err != nil {
		//	return
		//}

		for {
			//line, err := buf.ReadString('\n')
			nr, err := buf.Read(bufCnt)
			//line, _, err := buf.ReadLine()
			//lineStr := string(line) + "\n"
			if err == io.EOF || err != nil {
				//fmt.Println("file end ###############################")
				break
			}
			if nr > 0 {

				//rows ++
				//fmt.Println(rows)
				//
				//lineCnt += lineStr
				//
				//if rows % 1000 == 0 {
				fmt.Fprintf(w, string(bufCnt))
				//lineCnt = ""
				// Do not forget flushing streamed data to the client.
				if err := w.Flush(); err != nil {
					return
				}
				//time.Sleep(10 * time.Millisecond)
				//}
			}
		}

		//for i := 0; i < 10; i++ {
		//	fmt.Fprintf(w, "this is a message number %d", i)
		//
		//	// Do not forget flushing streamed data to the client.
		//	if err := w.Flush(); err != nil {
		//		return
		//	}
		//	time.Sleep(time.Second)
		//}
	})

	err = fasthttp.DoTimeout(freq, fresp, timeOut)
	if err != nil {
		return err
	}

	dataResponse.SetStatusCode(200)
	dataResponse.ReturnCode = "000000"
	return nil
}

func (ft *FastTransfer) Close() {

}

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
