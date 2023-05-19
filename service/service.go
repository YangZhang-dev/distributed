package service

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
)

// Start 服务启动程序
func Start(ctx context.Context, reg registry.Registration,
	registerHandlersFunc func()) (context.Context, error) {

	registerHandlersFunc()
	ctx = startService(ctx, reg)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}
func startService(ctx context.Context, reg registry.Registration) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var server http.Server
	server.Addr = reg.ServiceURL

	go func() {
		log.Println(server.ListenAndServe())
		err := registry.ShutdownService(reg)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop. ", reg.ServiceName)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(reg)
		if err != nil {
			log.Println(err)
		}
		server.Shutdown(ctx)
		cancel()
	}()
	return ctx
}
