package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context,
	serviceName, host, port string, registerHandlersFunc func()) (context.Context, error) {
	registerHandlersFunc()
	ctx = startService(ctx, serviceName, host, port)
	return ctx, nil
}
func startService(ctx context.Context, serviceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var server http.Server
	server.Addr = host + ":" + port

	go func() {
		log.Println(server.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop. ", serviceName)
		var s string
		fmt.Scanln(&s)
		server.Shutdown(ctx)
		cancel()
	}()
	return ctx
}
