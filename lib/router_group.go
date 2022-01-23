package lib

import (
	"net/http"
	"path"
	"strconv"
)

type RouterGroup struct {
	prefix        string
	middlewares   []HandlerFunc
	parantGrp     *RouterGroup
	engine        *Engine
	staticCounter int
}

func (grp *RouterGroup) CreateSubGroup(prefix string) *RouterGroup {
	result := &RouterGroup{
		prefix:    grp.prefix + prefix,
		engine:    grp.engine,
		parantGrp: grp,
	}
	grp.engine.allGroups = append(grp.engine.allGroups, result)
	return result
}

func (grp *RouterGroup) AppendMiddilewares(middlewares ...HandlerFunc) {
	grp.middlewares = append(grp.middlewares, middlewares...)
}

// Route binding now moves to here, so that groups control access

func (grp *RouterGroup) AddRoute(method, pattern string, handler HandlerFunc) {
	pattern = grp.prefix + pattern
	grp.engine.router.addRoute(method, pattern, handler)
}

func (grp *RouterGroup) GET(pattern string, handler HandlerFunc) {
	grp.AddRoute("GET", pattern, handler)
}

func (grp *RouterGroup) POST(pattern string, handler HandlerFunc) {
	grp.AddRoute("POST", pattern, handler)
}

/*
	relativePath:	url part related to current group
	rootPath: 		root path of static resources directory
*/
func (grp *RouterGroup) Static(relativePath string, rootPath string, strictMode bool) {
	param := "/*filepath" + strconv.Itoa(grp.staticCounter)
	pattern := path.Join(relativePath, param)
	absPathPrefix := path.Join(grp.prefix, relativePath)

	fileSystem := http.Dir(rootPath) // open
	// File Server: a handler that serves as a normal static web engine
	fileServer := http.StripPrefix(absPathPrefix, http.FileServer(fileSystem))
	grp.AddRoute("GET", pattern, func(c *Context) {
		file := c.Params[param]
		if _, err := fileSystem.Open(file); err != nil {
			// 404 not found
			c.SetStatusCode(http.StatusNotFound)
			return
		}
		// delegate the context to the file server
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
	grp.staticCounter++
}
