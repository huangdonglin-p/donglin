package framework

import (
	"log"
	"net/http"
)

// Core is the core structure of the framework, similar to Gin's Engine struct
type Core struct {
	router map[string]ControllerHandler
}

// NewCore returns *Core
func NewCore() *Core {
	return &Core{router: map[string]ControllerHandler{}}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}

// ServeHTTP is used to make Core implement the Handler interface of the http package
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.serveHTTP")
	ctx := NewContext(request, response)

	// 一个简单的路由选择器，这里直接写死为测试路由foo
	router := c.router["foo"]
	if router == nil {
		return
	}
	log.Println("core.router")

	router(ctx)
}
