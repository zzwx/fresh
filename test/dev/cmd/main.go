package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	var server *fasthttp.Server

	go func() {
		server = &fasthttp.Server{
			Handler: router.Handler,
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Server shutdown")
	server.Shutdown()
	fmt.Println("Server Exiting")
}
