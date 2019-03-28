package kuafu

import (
	"encoding/json"
	"github.com/linghaihui/kuafu/util"
	"io/ioutil"
	"net/http"
	"reflect"
)

const (
	JsonMime = "application/json; charset=utf-8"
	HtmlMime = "text/html; charset=utf-8"
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
	params        map[string]string
	Json          interface{}
}

// get param from query or variable in path
func (ctx *Context) GetParam(name string, args ...interface{}) string {
	if v, ok := ctx.params[name]; ok {
		return v
	}
	if len(args) > 0 {
		return util.FormatString("%v", args[0])
	}
	return ""
}

// bind request data to json, only the first time will work
func (ctx *Context) BindJSON(ptr interface{}) {
	if ctx.Json == nil {
		buff, _ := ioutil.ReadAll(ctx.Request.Body)
		t := reflect.TypeOf(ptr)
		if t.Kind() != reflect.Ptr {
			panic("obj arg must be a ptr")
		}
		if err := json.Unmarshal(buff, ptr); err != nil {
			panic("cannot bind data")
		}
	}
}

// return json
func (ctx *Context) JsonResponse(code int, data interface{}) {
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
	ctx.HttpCode = 404
	ctx.Response.WriteHeader(404)
}

// return specific code
func (ctx *Context) ResponseAnyCode(code int, resp ...interface{}) {
	ctx.Response.WriteHeader(code)
	ctx.HttpCode = code
	if len(resp) > 0 {
		buff := []byte(util.FormatString("%v", resp[0]))
		if _, err := ctx.Response.Write(buff); err != nil {
			panic(err)
		} else {
			ctx.responseBytes = append(ctx.responseBytes, buff...)
		}
	}
}

// return redirect response, such as 302, 301
func (ctx *Context) RedirectResponse(location string, code int) {
	header := ctx.Response.Header()
	header["Location"] = []string{location}
	header["Content-Type"] = []string{HtmlMime}
	redirectHtml := util.FormatString(`
        '<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">\n'
        '<title>Redirecting...</title>\n'
        '<h1>Redirecting...</h1>\n'
        '<p>You should be redirected automatically to target URL: '
        '<a href="%s">%s</a>.  If not click the link.'`, location, "redirect to")
	ctx.HttpCode = code
	ctx.Response.WriteHeader(code)
	if _, err := ctx.Response.Write([]byte(redirectHtml)); err != nil {
		panic(err)
	}
	ctx.responseBytes = append(ctx.responseBytes, []byte(redirectHtml)...)
}

func (ctx *Context) Next() {
	ctx.index++
	for total := len(ctx.HandlerChain); ctx.index <= total; ctx.index++ {
		ctx.HandlerChain[ctx.index-1](ctx)
	}
}
