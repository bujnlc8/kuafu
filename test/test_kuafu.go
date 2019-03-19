package main

import "github.com/linghaihui/kuafu"

type Hello struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func SayHello(ctx *kuafu.Context) {
	ctx.Json(200, Hello{Code: 0, Msg: "Hello World"})
}

func main() {
	server := kuafu.NewServer()
	server.SetDebugMode()
	server.Use(kuafu.PrintRequest, kuafu.Handler404)
	registry := server.NewRegistry()
	registry.GET("/hello", SayHello)
	server.Run("127.0.0.1:9999")
}
