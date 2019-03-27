package kuafu

import (
	"encoding/json"
	"github.com/linghaihui/kuafu/util"
	"net/http"
)

const (
	JsonMime = "application/json; charset=utf-8"
	TextMime = "text/plain; charset=utf-8"
)

type Context struct {
	Server        *Server
	Request       *http.Request
	Response      http.ResponseWriter
	Session       map[string]string
	HandlerChain  []HandlerFunc
	index         int
	responseBytes []byte
	HttpCode      int
	Params        map[string]string
}

func (ctx *Context) GetParam(name string, args ...interface{}) string {
	if v, ok := ctx.Params[name]; ok {
		return v
	}
	if len(args) > 0 {
		return util.FormatString("%v", args[0])
	}
	return ""
}

// return json
func (ctx *Context) Json(code int, data interface{}) {
	if jsonBytes, err := json.Marshal(data); err != nil {
		panic(err)
	} else {
		header := ctx.Response.Header()
		header["Content-Type"] = []string{JsonMime}
		ctx.Response.WriteHeader(code)
		ctx.HttpCode = code
		if _, err := ctx.Response.Write(jsonBytes); err != nil {
			panic(err)
		}
		ctx.responseBytes = append(ctx.responseBytes, jsonBytes...)
	}
}

// return 404
func (ctx *Context) Response404() {
	ctx.Response.WriteHeader(404)
}

func (ctx *Context) Next() {
	ctx.index++
	for total := len(ctx.HandlerChain); ctx.index <= total; ctx.index++ {
		ctx.HandlerChain[ctx.index-1](ctx)
	}
}
