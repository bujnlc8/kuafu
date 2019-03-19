package kuafu

import (
	"fmt"
	"net/http"
	"strings"
)

type Server struct {
	Routers     map[string]*Router
	MiddleWares []HandlerFunc
	Debug       bool
}

func (server *Server) Use(m ...HandlerFunc) {
	server.MiddleWares = append(server.MiddleWares, m...)
}

func NewServer() *Server {
	return &Server{
		Routers:     make(map[string]*Router),
		MiddleWares: []HandlerFunc{KuafuMark},
		Debug:       false,
	}
}

func (server *Server) NewRegistry() *Registry {
	return &Registry{server: server}
}

func (server *Server) findRouter(ctx *Context) (error, *Router) {
	router := &Router{Path: ctx.request.URL.Path, Method: ctx.request.Method}
	if v, ok := server.Routers[router.ToString()]; ok {
		return nil, v
	} else {
		return RouterNotFound, nil
	}
}

func (server *Server) NewContext(request *http.Request, response http.ResponseWriter) *Context {
	ctx := Context{
		server:       server,
		request:      request,
		response:     response,
		session:      make(map[string]string),
		handlerChain: nil,
		index:        0,
	}
	return &ctx
}

func (server *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := server.NewContext(req, resp)
	for _, v := range server.MiddleWares {
		ctx.handlerChain = append(ctx.handlerChain, v)
	}
	if err, router := server.findRouter(ctx); err != nil {
		if server.Debug {
			fmt.Println(req.URL.Path, err)
		}
		ctx.httpCode = 404
	} else {
		// merge router handle and middleware
		ctx.handlerChain = append(ctx.handlerChain, router.Handler)
	}
	ctx.Next()
}

func (server *Server) routers() []string {
	rst := []string{}
	for _, v := range server.Routers {
		rst = append(rst, v.ToString())
	}
	return rst
}

// set debug mode so that you can see more log
func (server *Server) SetDebugMode() {
	server.Debug = true
}

func (server *Server) Run(addr string) {
	fmt.Println(fmt.Sprintf("Kuafu is listening at %s", addr))
	fmt.Println("router list:")
	fmt.Println(strings.Join(server.routers(), "\n"))
	if err := http.ListenAndServe(addr, server); err != nil {
		panic(err)
	}
}
