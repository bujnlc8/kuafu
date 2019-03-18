package kuafu

import (
	"encoding/json"
	"net/http"
)

const (
	JsonMime = "application/json; charset=utf8"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	session  map[string]string
}

func NewContext(request *http.Request, response http.ResponseWriter) *Context {
	ctx := Context{
		request:  request,
		response: response,
		session:  make(map[string]string),
	}
	return &ctx
}

func (ctx *Context) Json(code int, data interface{}) {
	if jsonBytes, err := json.Marshal(data); err != nil {
		panic(err)
	} else {
		header := ctx.response.Header()
		header["Content-Type"] = []string{JsonMime}
		ctx.response.WriteHeader(code)
		if _, err := ctx.response.Write(jsonBytes); err != nil {
			panic(err)
		}
	}
}
