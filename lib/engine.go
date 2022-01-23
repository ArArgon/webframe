package lib

import (
	"net/http"
	"strings"
)

type Engine struct {
	*RouterGroup // Engine itself has all capabilities of a group
	router       *Router
	allGroups    []*RouterGroup
}

func (e *Engine) SetNotFound(handler HandlerFunc) {
	e.router.notFound = handler
}

func (e *Engine) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// append middlewares
	requestCtx := newContext(rw, r)
	for _, group := range e.allGroups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			requestCtx.Handlers = append(requestCtx.Handlers, group.middlewares...)
		}
	}
	e.router.handleContext(requestCtx)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func New() *Engine {
	result := &Engine{
		router: newRouter(func(ctx *Context) {
			ctx.HTML(404, "<html> <head> <title>Error</title> </head> <body> <h1>"+ctx.Path+": 404 Not Found</h1> </body> </html>\r\n")
		}),
	}
	// Engine's group: supreme group
	result.RouterGroup = &RouterGroup{engine: result}
	// The only group
	result.allGroups = []*RouterGroup{result.RouterGroup}
	// enable panic recovery support
	result.AppendMiddilewares(Recovery())
	return result
}
