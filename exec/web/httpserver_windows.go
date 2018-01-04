package web

import (
	"fmt"
	"net"
	"github.com/valyala/fasthttp"
)

/**
    Author: luzequan
    Created: 2018-01-04 14:04:14
*/

type HttpServer struct{}

func (s *HttpServer) Run() {
	router := NewHttpRouter()
	router.Register()

	host := "127.0.0.1"
	port := 8899

	uri := fmt.Sprintf("%s:%d", host, port)
	ln, err := net.Listen("tcp4", uri)
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