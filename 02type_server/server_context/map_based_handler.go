package server_context

import (
	"fmt"
	"net/http"
)

type Routable interface {
	Route(method, pattern string, handleFunc func(ctx *Context))
}

type Handler interface {
	http.Handler
	Routable
}

type HandlerBasedOnMap struct {
	//key: method + url
	handlers map[string]func(ctx *Context)
}

//检测是否实现
var _ Handler = &HandlerBasedOnMap{}

func (h *HandlerBasedOnMap) Route(method, pattern string, handleFunc func(ctx *Context)) {
	//http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
	//	ctx := NewContext(w, r)
	//	handleFunc(ctx)
	//})
	key := h.Key(method, pattern)
	h.handlers[key] = handleFunc
}

func (h *HandlerBasedOnMap) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := h.Key(request.Method, request.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		ctx := NewContext(writer, request)
		handler(ctx)
	} else {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("not any router match"))
	}
}

func (h *HandlerBasedOnMap) Key(method, path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}

func NewHandlerBasedOnMap() Handler {
	return &HandlerBasedOnMap{handlers: map[string]func(ctx *Context){}}
}
