package main

import (
	"context"
	"distributed/log"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

const (
	host         = "127.0.0.1"
	port         = ":9988"
	serviceName  = registry.LogService
	updateURL    = "/services"
	heartbeatURL = "/heartbeat"
)

func main() {
	log.Run()
	serviceAddress := fmt.Sprintf("http://%v%v", host, port)
	r := registry.Registration{
		ServiceName:      serviceName,
		ServiceURL:       serviceAddress,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceAddress + updateURL,
		HeartbeatURL:     serviceAddress + heartbeatURL,
	}
	ctx, err := service.Start(
		context.Background(),
		r,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("shutting down log service.")
}
