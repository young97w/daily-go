package main

import (
	"fmt"
	v2 "geektime/web/v2"
)

func main() {
	s := v2.NewHTTPServer()
	s.GET("/user", hello)
	s.Start("localhost:8082")
}

func hello(ctx *v2.Context) {
	fmt.Println("收到请求了")
	fmt.Fprintf(ctx.Resp, "你好user")
}
