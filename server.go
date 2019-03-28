package kuafu

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type RouterStorage struct {
	routers map[string]*Router
}

type Server struct {
	Routers     map[string]*RouterStorage
	MiddleWares []HandlerFunc
	debug       bool
}

func (server *Server) Use(m ...HandlerFunc) {
	server.MiddleWares = append(server.MiddleWares, m...)
}

func NewServer() *Server {
	return &Server{
		Routers:     make(map[string]*RouterStorage),
		MiddleWares: []HandlerFunc{KuafuMark, PrintRequest, Handler404},
		debug:       false,
	}
}

func (server *Server) NewRegistry() *Registry {
	return &Registry{server: server}
}

func (server *Server) NewGroup(name string, prefix string) *Group {
	return &Group{server: server, name: name, prefix: prefix}
}

var reg, _ = regexp.Compile("<:(.*?)>")

// find the handle
func (server *Server) findRouter(ctx *Context) (error, *Router) {
	method := ctx.Request.Method
	path := ctx.Request.URL.Path
	router := &Router{path: path, method: method}
	if v, ok := server.Routers[method].routers[router.String()]; ok {
		return nil, v
	} else {
		routerStorage := server.Routers[method]
		if routerStorage != nil {
			for _, r := range routerStorage.routers {
				if ok, _ := regexp.MatchString(reg.String(), r.path); !ok {
					continue
				}
				params := reg.FindAllStringSubmatch(r.path, -1)
				var paramName []string
				for _, v := range params {
					paramName = append(paramName, v[1])
				}
				regAnother, _ := regexp.Compile("^" + reg.ReplaceAllString(r.path, "(.*?)") + "$")
				pathValue := regAnother.FindStringSubmatch(path)
				if len(pathValue)-1 == len(paramName) {
					for k, v := range paramName {
						ctx.params[v] = pathValue[k+1]
					}
					return nil, r
				}
			}
		}
	}
	return RouterNotFound, nil
}

func (server *Server) NewContext(request *http.Request, response http.ResponseWriter) *Context {
	ctx := Context{
		Server:       server,
		Request:      request,
		Response:     response,
		Session:      make(map[string]string),
		HandlerChain: nil,
		index:        0,
		params:       make(map[string]string),
	}
	//put param in ctx.params
	for k, v := range request.URL.Query() {
		ctx.params[k] = v[0]
	}
	return &ctx
}

func (server *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := server.NewContext(req, resp)
	for _, v := range server.MiddleWares {
		ctx.HandlerChain = append(ctx.HandlerChain, v)
	}
	if err, router := server.findRouter(ctx); err != nil {
		if server.debug {
			fmt.Println(req.URL.Path, err)
		}
		ctx.HttpCode = 404
	} else {
		// merge router handle and middleware
		ctx.HandlerChain = append(ctx.HandlerChain, router.handler)
	}
	ctx.Next()
}

func (server *Server) routers() []string {
	var rst []string
	for _, routerStorage := range server.Routers {
		for _, router := range routerStorage.routers {
			rst = append(rst, router.String())
		}
	}
	return rst
}

// set debug mode so that you can see more log
func (server *Server) SetDebugMode() {
	server.debug = true
}

func (server *Server) Run(addr string) {
	fmt.Println(fmt.Sprintf("Kuafu is listening at %s", addr))
	fmt.Println("router list:")
	fmt.Println(strings.Join(server.routers(), "\n"))
	if err := http.ListenAndServe(addr, server); err != nil {
		panic(err)
	}
}
