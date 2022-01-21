package lib

import (
	"log"
	"strings"
)

type Router struct {
	trees    map[string]TrieTree
	router   map[string]HandlerFunc
	notFound HandlerFunc
}

func newRouter(notFoundHandler HandlerFunc) *Router {
	return &Router{
		router:   make(map[string]HandlerFunc),
		trees:    make(map[string]TrieTree),
		notFound: notFoundHandler,
	}
}

func parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	var result []string

	for _, part := range parts {
		if part == "" {
			continue
		}
		result = append(result, part)
		if part[0] == '*' {
			break
		}
	}
	return result
}

func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("[Router] Routing %4s, %s", method, pattern)
	if _, ok := r.trees[method]; !ok {
		r.trees[method] = *newTrieTree()
	}

	parts := parsePattern(pattern)
	tree := r.trees[method]
	key := method + "-" + pattern

	tree.addPath(pattern, parts)
	r.router[key] = handler
}

func procWildcard(pattern string, path []string) map[string]string {
	patternParts := parsePattern(pattern)
	res := make(map[string]string)

loop:
	for idx, part := range patternParts {
		switch part[0] {
		case ':':
			// matchOne
			res[part[1:]] = path[idx]
		case '*':
			// matchRest
			res[part[1:]] = strings.Join(path[idx:], "/")
			break loop
		}
	}
	return res
}

func (r *Router) handleContext(ctx *Context) {
	tree, hasSuchMethod := r.trees[ctx.Method]
	parts := parsePattern(ctx.Path)
	if !hasSuchMethod {
		r.trees[ctx.Method] = *newTrieTree()
	}

	if node, ok := tree.matchPath(parts); ok {
		key := ctx.Method + "-" + node.pattern
		if handler, ok := r.router[key]; ok {
			log.Printf("[Router] %s Matched> %s ==> %s", ctx.Method, ctx.Path, node.pattern)
			ctx.Params = procWildcard(node.pattern, parts)
			handler(ctx)
			return
		}
	}
	// Not Found
	log.Printf("[Router] %s Unmatched> %s", ctx.Method, ctx.Path)
	r.notFound(ctx)
}
