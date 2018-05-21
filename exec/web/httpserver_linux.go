package web

import (
	"fmt"
	"github.com/valyala/fasthttp/reuseport"
	"github.com/valyala/fasthttp"

	. "drcs/settings"
)

/**
    Author: luzequan
    Created: 2018-01-02 19:39:59
*/
type HttpServer struct{}

func (s *HttpServer) Run() {
	router := NewHttpRouter()
	router.Register()

	var (
		host string
		port int
	)
	for {
		common := GetCommonSettings()
		if common != nil {

			host = common.Node.Host
			port = common.Node.Port
			break
		}
	}

	uri := fmt.Sprintf("%s:%d", host, port)
	ln, err := reuseport.Listen("tcp4", uri)
	if err != nil {
	}

	server := &fasthttp.Server{
		Handler: router.Router.Handler,
		//TODO:body文件大小设置成20G存在安全隐患，在不需要文件时关闭该设置！
		MaxRequestBodySize: 20 * 1024 * 1024 * 1024, //set maxbody size = 20G
	}

	if err = server.Serve(ln); err != nil {
	}
}
