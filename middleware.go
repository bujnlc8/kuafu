package kuafu

import (
	"log"
)

// print request after response sent
func PrintRequest(ctx *Context)  {
	ctx.Next()
	log.Println(ctx.request.Method, ctx.request.URL.Path, ctx.httpCode)
	if ctx.server.Debug {
		log.Println("response body:", string(ctx.responseBytes))
	}
}


// handle 404
func Handler404(ctx *Context)  {
	if ctx.httpCode == 404{
		ctx.Response404()
	}else {
		ctx.Next()
	}
}


// add kuafu mark in http header
func KuafuMark(ctx *Context)  {
	ctx.response.Header().Add("X-Server-Framework", FormatString("Kuafu/%s", Version))
	ctx.Next()
}