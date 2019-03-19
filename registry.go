package kuafu

import "fmt"

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

type HandlerFunc func(ctx *Context)

type Registry struct {
	server *Server
}

func (registry *Registry) GET(path string, handler HandlerFunc) {
	router := &Router{Method: GET, Path: path, Handler: handler}
	s := router.ToString()
	if _, ok := registry.server.Routers[s]; ok {
		panic(fmt.Sprintf("redeclare %s", s))
	} else {
		registry.server.Routers[s] = router
	}
}

func (registry *Registry) POST(path string, handler HandlerFunc) {
	router := &Router{Method: POST, Path: path, Handler: handler}
	s := router.ToString()
	if _, ok := registry.server.Routers[s]; ok {
		panic(fmt.Sprintf("redeclare %s", s))
	} else {
		registry.server.Routers[s] = router
	}
}

func (registry *Registry) DELETE(path string, handler HandlerFunc) {
	router := &Router{Method: DELETE, Path: path, Handler: handler}
	s := router.ToString()
	if _, ok := registry.server.Routers[s]; ok {
		panic(fmt.Sprintf("redeclare %s", s))
	} else {
		registry.server.Routers[s] = router
	}
}

func (registry *Registry) PUT(path string, handler HandlerFunc) {
	router := &Router{Method: PUT, Path: path, Handler: handler}
	s := router.ToString()
	if _, ok := registry.server.Routers[s]; ok {
		panic(fmt.Sprintf("redeclare %s", s))
	} else {
		registry.server.Routers[s] = router
	}
}

type Router struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

func (router *Router) ToString() string {
	return fmt.Sprintf("%s_%s", router.Path, router.Method)
}
