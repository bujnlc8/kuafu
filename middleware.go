package kuafu

import (
	"github.com/linghaihui/kuafu/util"
	"log"
)

// print request after response sent
func PrintRequest(ctx *Context) {
	ctx.Next()
	log.Println(ctx.Request.Method, ctx.Request.URL.Path, ctx.HttpCode)
	if ctx.Server.debug {
		log.Println("response body:", string(ctx.responseBytes))
	}
}

// handle 404
func Handler404(ctx *Context) {
	if ctx.HttpCode == 404 {
		ctx.Response404()
	} else {
		ctx.Next()
	}
}

// add kuafu mark in http header
func KuafuMark(ctx *Context) {
	ctx.Response.Header().Add("X-Server-Framework", util.FormatString("Kuafu/%s", Version))
	ctx.Next()
}
