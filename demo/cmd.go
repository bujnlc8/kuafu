package main

import (
	"github.com/linghaihui/kuafu"
	"github.com/linghaihui/kuafu/logger"
)

var log = logger.NewSizeLogger("http.log", 0, logger.LevelDebug)

func main() {
	server := kuafu.NewServer()
	registry := server.NewRegistry()
	registry.PUT("/add", AddBook)
	registry.GET("/get/<:id>", GetBook)
	registry.DELETE("/delete/<:id>", DeleteBook)
	registry.GET("/redirect", func(ctx *kuafu.Context) {
		ctx.RedirectResponse("https://www.baidu.com", 302)
	})
	server.SetDebugMode()
	server.Run("127.0.0.1:9999")
}
