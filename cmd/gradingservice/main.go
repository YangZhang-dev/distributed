package main

import (
	"context"
	"distributed/grades"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

const (
	host        = "127.0.0.1"
	port        = ":6000"
	serviceName = registry.GradingService
)

func main() {

	r := registry.Registration{
		ServiceName: serviceName,
		ServiceURL:  fmt.Sprintf("%v%v", host, port),
	}
	ctx, err := service.Start(
		context.Background(),
		r,
		grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
