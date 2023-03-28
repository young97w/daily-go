package web

import (
	"log"
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error

	// addRoute 注册一个路由
	// method 是 HTTP 方法
	addRoute(method string, path string, handler HandleFunc)
	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type HTTPServer struct {
	router
	mdls []Middleware
}

type ServerOption func(server *HTTPServer)

func NewHTTPServer(opts ...ServerOption) *HTTPServer {
	s := &HTTPServer{
		router: newRouter(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ServeHTTP HTTPServer 处理请求的入口
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	s.serve(ctx)
}

// Start 启动服务器
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) Post(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n == nil || mi.n.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("Not Found"))
		return
	}
	ctx.MatchedRoute = mi.n.path
	ctx.PathParams = mi.pathParams
	root := mi.n.handler

	//链接中间件
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}
	root = flashResp(root)
	root(ctx)
}

//响应
func flashResp(next HandleFunc) HandleFunc {
	return func(ctx *Context) {
		//执行中间件
		next(ctx)
		//写状态码
		if ctx.RespStatusCode > 0 {
			ctx.Resp.WriteHeader(ctx.RespStatusCode)
		}
		//写响应数据
		//fmt.Println(string(ctx.RespData))
		_, err := ctx.Resp.Write(ctx.RespData)
		if err != nil {
			log.Fatalln("响应失败", err)
		}
	}
}

//Use 注册中间件
func (s *HTTPServer) Use(mdls ...Middleware) {
	if s.mdls == nil {
		s.mdls = mdls
		return
	}
	s.mdls = append(s.mdls, mdls...)
}

func ServerWithMiddleware(mdls ...Middleware) ServerOption {
	return func(s *HTTPServer) {
		if s.mdls == nil {
			s.mdls = mdls
		} else {
			s.mdls = append(s.mdls, mdls...)
		}
	}
}
