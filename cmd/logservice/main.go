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
	host        = "127.0.0.1"
	port        = ":8888"
	serviceName = registry.LogService
	logFileName = "distributed.log"
)

func main() {
	log.Run(logFileName)

	r := registry.Registration{
		ServiceName: serviceName,
		ServiceURL:  fmt.Sprintf("%v%v", host, port),
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
