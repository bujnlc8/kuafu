package main

import "github.com/linghaihui/kuafu"

type Hello struct {
	Msg string `json:"s"`
	Code int `json:"code"`
}

func SayHello(ctx *kuafu.Context)  {
	ctx.Json(200, Hello{Code:0, Msg:"Hello World"})
}

func main()  {
	server := kuafu.NewServer()
	registry := server.NewRegistry()
	registry.GET("/hello", SayHello)
	server.Run("127.0.0.1:9999")
}