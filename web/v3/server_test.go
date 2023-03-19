package v2

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHTTPServer_Start(t *testing.T) {
	s := NewHTTPServer()
	s.addRoute(http.MethodGet, "/user", hello)
	s.Start(":8081")
}

func hello(ctx *Context) {
	fmt.Println("收到请求了")
	fmt.Fprintf(ctx.Resp, "你好user")
}
