package v1

import (
	"fmt"
	"net"
	"net/http"
)

type HandleFunc func(ctx Context)

//Server 抽象
type Server interface {
	http.Handler
	//启动服务器
	Start(addr string) error
	//注册路由
	AddRoute(method, path string, handler HandleFunc)
}

type HTTPServer struct {
}

var _ Server = &HTTPServer{}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	ctx := &Context{
		Req:  *request,
		Resp: writer,
	}
	s.Serve(ctx)
}

func (s *HTTPServer) Serve(ctx *Context) error {
	//TODO implement me
	fmt.Println("启动服务器！")
	return nil
}

func (s *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return http.Serve(l, s)
}

func (s *HTTPServer) AddRoute(method, path string, handler HandleFunc) {
	//TODO implement me

}
