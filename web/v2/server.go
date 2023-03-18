package v2

import "net/http"

//Server 注册路由，
//构建路由树
//生命周期控制（启动关闭）
//web到http（寻找handler）
type Server interface {
	http.Handler
	addRoute(method, path string, handleFunc HandleFunc)
	Start(addr string) error
}

type HandleFunc func(ctx *Context)

type HTTPServer struct {
	router
}

//检测HTTPServer是否实现了 Server

var _ Server = &HTTPServer{}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		newRouter(),
	}
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) Serve(ctx *Context) error {
	//TODO
	return nil
}

func (s *HTTPServer) GET(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) POST(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *HTTPServer) PUT(path string, handler HandleFunc) {
	s.addRoute(http.MethodPut, path, handler)
}
