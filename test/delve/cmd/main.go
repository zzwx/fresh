package main

import (
	"fmt"
	"os"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func exampleRequestHandler(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Hello World asd")
}

func main() {

	router := router.New()
	router.GET("/", exampleRequestHandler)

	fmt.Println("Server Listening at 0.0.0.0:8101")

	var shutdownCh = make(chan os.Signal, 1)

	var server *fasthttp.Server

	go func(server *fasthttp.Server) {
		server = &fasthttp.Server{
			Handler: router.Handler,
		}
	}(server)

	<-shutdownCh
	server.Shutdown()
}
