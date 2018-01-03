package transfer

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"sync"
	neturl "net/url"
	. "dat/core/interaction/response"
	"time"
)

/**
    Author: luzequan
    Created: 2017-12-28 15:41:49
*/

type FastTransfer struct {}

func NewFastTransfer() Transfer {
	return &FastTransfer{}
}

// 封装fasthttp服务
func (ft *FastTransfer) ExecuteMethod(req Request) Response {
	fmt.Println("execute fasthttp")
	timeout := 30*1000
	dataResponse := &DataResponse{}

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
		break
	}

	return dataResponse
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
