package kuafu

import (
	"fmt"
	"net/http"
	"strings"
)

var KuaFu = map[string]string{
	"version": "0.0.1",
}

type Server struct {
	Routers map[string]*Router
}

func NewServer() *Server {
	return &Server{Routers: make(map[string]*Router)}
}

func (server *Server) NewRegistry() *Registry {
	return &Registry{server: server}
}

func (server *Server) findRouter(ctx *Context) (error, *Router) {
	fmt.Println(ctx.request.URL.Path)
	router := &Router{Path: ctx.request.URL.Path, Method: ctx.request.Method}
	if v, ok := server.Routers[router.ToString()]; ok {
		return nil, v
	} else {
		return RouterNotFound, nil
	}
}

func (server *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := NewContext(req, resp)
	if err, router := server.findRouter(ctx); err != nil {
		fmt.Println(fmt.Sprintf("find router error happened %s", err))
	} else {
		router.Do(ctx)
	}
}

func (server *Server) routers() []string {
	rst := []string{}
	for _, v := range server.Routers {
		rst = append(rst, v.ToString())
	}
	return rst
}

func (server *Server) Run(addr string) {
	fmt.Println(fmt.Sprintf("Kuafu is listening at %s", addr))
	fmt.Println("router list:")
	fmt.Println(strings.Join(server.routers(), "\n"))
	if err := http.ListenAndServe(addr, server); err != nil {
		panic(err)
	}
}
