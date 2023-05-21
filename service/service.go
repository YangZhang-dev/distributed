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
	// 注册每个服务的handlers
	registerHandlersFunc()

	ctx = startService(ctx, reg)

	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

// 启动服务程序
func startService(ctx context.Context, reg registry.Registration) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var server http.Server
	// 这里只取端口，监听所有地址
	port := strings.Split(reg.ServiceURL, ":")[2]
	server.Addr = ":" + port

	// 一个go routine用于启动服务
	go func() {
		log.Println(server.ListenAndServe())
		err := registry.ShutdownService(reg)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()
	// 另一个再控制台等待用户输入字符后停止服务
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
