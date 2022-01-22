package lib

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parantGrp   *RouterGroup
	engine      *Engine
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
