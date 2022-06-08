package main

import (
	"fmt"
	"net/http"

	"donglin.framework.use/framework"
)

func main() {
	server := &http.Server{
		Handler: framework.NewCore(),
		Addr:    "Localhost:8080",
	}

	fmt.Println("启动服务")
	server.ListenAndServe()
}
