package service

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Start 服务启动程序
func Start(ctx context.Context, reg registry.Registration,
	registerHandlersFunc func()) (context.Context, error) {

	registerHandlersFunc()
	// 启动服务程序
	ctx = startService(ctx, reg)
	// 在注册中心进行注册
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}
func startService(ctx context.Context, reg registry.Registration) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var server http.Server
	port := strings.Split(reg.ServiceURL, ":")[2]
	server.Addr = ":" + port

	go func() {
		log.Println(server.ListenAndServe())
		err := registry.ShutdownService(reg)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop. \n", reg.ServiceName)
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
