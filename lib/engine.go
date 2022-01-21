package lib

import (
	"net/http"
)

type Engine struct {
	router *Router
}

// primitive way of mapping paths
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRoute("POST", pattern, handler)
}

func (e *Engine) SetNotFound(handler HandlerFunc) {
	e.router.notFound = handler
}

func (e *Engine) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	e.router.handleContext(newContext(rw, r))
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func New() *Engine {
	return &Engine{
		router: newRouter(func(ctx *Context) {
			ctx.HTML(404, "<html> <head> <title>Error</title> </head> <body> <h1>"+ctx.Path+": 404 Not Found</h1> </body> </html>\r\n")
		}),
	}
}
