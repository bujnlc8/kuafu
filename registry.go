package kuafu

import "github.com/linghaihui/kuafu/util"

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

func (registry *Registry) doCommon(method string, path string, handler HandlerFunc) {
	router := &Router{
		method:  method,
		path:    path,
		handler: handler,
	}
	s := router.String()
	routerStorage := registry.server.Routers[method]
	if routerStorage == nil {
		routerStorage = &RouterStorage{
			routers: make(map[string]*Router)}
		registry.server.Routers[method] = routerStorage
	} else {
		if _, ok := registry.server.Routers[method].routers[s]; ok {
			panic(util.FormatString("redeclare %s", s))
		}
	}
	routerStorage.routers[s] = router
}

func (registry *Registry) GET(path string, handler HandlerFunc) {
	registry.doCommon(GET, path, handler)
}

func (registry *Registry) POST(path string, handler HandlerFunc) {
	registry.doCommon(POST, path, handler)
}

func (registry *Registry) DELETE(path string, handler HandlerFunc) {
	registry.doCommon(DELETE, path, handler)
}

func (registry *Registry) PUT(path string, handler HandlerFunc) {
	registry.doCommon(PUT, path, handler)
}

type Router struct {
	method  string
	path    string // 支持/a/<:b>/c这种形式
	handler HandlerFunc
	group   *Group
}

type Group struct {
	name   string
	prefix string
	server *Server
}

func (group *Group) doCommon(method string, path string, handler HandlerFunc) {
	router := &Router{
		method:  method,
		path:    group.prefix + path,
		handler: handler,
		group:   group,
	}
	s := router.String()
	routerStorage := group.server.Routers[method]
	if routerStorage == nil {
		routerStorage = &RouterStorage{
			routers: make(map[string]*Router)}
		group.server.Routers[method] = routerStorage
	} else {
		if _, ok := group.server.Routers[method].routers[s]; ok {
			panic(util.FormatString("redeclare router %s", s))
		}
	}
	routerStorage.routers[s] = router
}

func (group *Group) GET(path string, handler HandlerFunc) {
	group.doCommon(GET, path, handler)
}

func (group *Group) POST(path string, handler HandlerFunc) {
	group.doCommon(POST, path, handler)
}

func (group *Group) DELETE(path string, handler HandlerFunc) {
	group.doCommon(DELETE, path, handler)
}

func (group *Group) PUT(path string, handler HandlerFunc) {
	group.doCommon(PUT, path, handler)
}

func (router *Router) String() string {
	return util.FormatString("%s||%s", router.method, router.path)
}
