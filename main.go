package main

import (
	"fmt"
	"net/http"

	"donglin.framework.use/framework"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    "Localhost:8080",
	}

	fmt.Println("启动服务")
	server.ListenAndServe()
}
