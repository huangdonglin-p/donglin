package framework

import "net/http"

// Core is the core structure of the framework, similar to Gin's Engine struct
type Core struct {
}

// NewCore returns *Core
func NewCore() *Core {
	return &Core{}
}

// ServeHTTP is used to make Core implement the Handler interface of the http package
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// TODO
}
