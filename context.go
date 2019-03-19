package kuafu

import (
	"encoding/json"
	"net/http"
)

const (
	JsonMime = "application/json; charset=utf-8"
	TextMime = "text/plain; charset=utf-8"
)

type Context struct {
	server        *Server
	request       *http.Request
	response      http.ResponseWriter
	session       map[string]string
	handlerChain  []HandlerFunc
	index         int
	responseBytes []byte
	httpCode      int
}

// return json
func (ctx *Context) Json(code int, data interface{}) {
	if jsonBytes, err := json.Marshal(data); err != nil {
		panic(err)
	} else {
		header := ctx.response.Header()
		header["Content-Type"] = []string{JsonMime}
		ctx.response.WriteHeader(code)
		ctx.httpCode = code
		if _, err := ctx.response.Write(jsonBytes); err != nil {
			panic(err)
		}
		ctx.responseBytes = append(ctx.responseBytes, jsonBytes...)
	}
}

// return 404
func (ctx *Context) Response404() {
	ctx.response.WriteHeader(404)
}

func (ctx *Context) Next() {
	ctx.index++
	for total := len(ctx.handlerChain); ctx.index <= total; ctx.index++ {
		ctx.handlerChain[ctx.index-1](ctx)
	}
}
