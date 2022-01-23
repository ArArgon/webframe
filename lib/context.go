package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONObject map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	StatusCode int
	Params     map[string]string

	// Middleware support
	Middlewares []HandlerFunc
	MidIndex    int
}
type HandlerFunc func(*Context)

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Add(key, value)
}

func (c *Context) SetStatusCode(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) String(code int, format string, val ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatusCode(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, val...)))
}

func (c *Context) JSON(code int, obj map[string]interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatusCode(code)

	encoder := json.NewEncoder(c.Writer)

	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.SetStatusCode(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetStatusCode(code)
	c.SetHeader("Content-Type", "text/html")
	c.Writer.Write([]byte(html))
}

func (ctx *Context) Next() {

}

func (ctx *Context) Fail(code int, message string) {
	ctx.JSON(code, JSONObject{"errMessage": message})
}

func newContext(
	writer http.ResponseWriter,
	req *http.Request) *Context {
	return &Context{
		Writer:  writer,
		Request: req,
		Path:    req.URL.Path,
		Method:  req.Method,
		Params:  make(map[string]string),
	}
}
