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
	//新建context
	//转交给serve
	ctx := &Context{
		Req:  *request,
		Resp: writer,
	}
	s.Serve(ctx)
}

func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) Serve(ctx *Context) {
	//先find node
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n == nil || mi.n.handler == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}

	mi.n.handler(ctx)

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
