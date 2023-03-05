package server_context

import (
	"encoding/json"
	"io"
	"net/http"
)

type Server interface {
	Route(method, pattern string, handleFunc func(ctx *Context))
	Start(address string) error
}

type sdkHttpServer struct {
	Name    string
	Handler *HandlerBasedOnMap
}

func (s *sdkHttpServer) Route(method, pattern string, handleFunc func(ctx *Context)) {
	//http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
	//	ctx := NewContext(w, r)
	//	handleFunc(ctx)
	//})
	key := s.Handler.Key(method, pattern)
	s.Handler.handlers[key] = handleFunc
}

func (s sdkHttpServer) Start(address string) error {
	http.Handle("/", s.Handler)
	return http.ListenAndServe(address, nil)
}

func NewServer() Server {
	return &sdkHttpServer{}
}

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		W: writer,
		R: request,
	}
}

func (c *Context) ReadJson(obj interface{}) error {
	r := c.R
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) WriteJson(code int, resp interface{}) error {
	c.W.WriteHeader(code)
	respJson, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = c.W.Write(respJson)
	if err != nil {
		return err
	}
	return nil
}

func SignUp(c *Context) {
	req := &signUpReq{}
	err := c.ReadJson(req)
	if err != nil {
		//do something
	}
	res := commonResponse{
		BizCode: 666,
		Msg:     "交易成功，请注意服务质量！",
		Data:    nil,
	}
	err = c.WriteJson(http.StatusOK, res)
	if err != nil {
		//
	}
}

type signUpReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
