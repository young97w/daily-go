package main

import (
	"fmt"
	v2 "geektime/web/v2"
	"net/http"
)

func main() {
	s := v2.NewHTTPServer()
	s.addRoute(http.MethodGet, "/user", hello)
	s.Start("localhost:8082")
}

func hello(ctx *v2.Context) {
	fmt.Println("收到请求了")
	fmt.Fprintf(ctx.Resp, "你好user")
}
