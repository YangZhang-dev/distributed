package main

import (
	"context"
	"distributed/log"
	"distributed/portal"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

const (
	host         = "127.0.0.1"
	port         = ":5000"
	serviceName  = registry.PortalService
	updateURL    = "/services"
	heartbeatURL = "/heartbeat"
)

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}
	serviceAddress := fmt.Sprintf("http://%v%v", host, port)

	r := registry.Registration{
		ServiceName: registry.PortalService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
			registry.GradingService,
		},
		ServiceUpdateURL: serviceAddress + updateURL,
		HeartbeatURL:     serviceAddress + heartbeatURL,
	}

	ctx, err := service.Start(context.Background(),
		r,
		portal.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	LogServiceURLs := registry.GetProvider(registry.LogService)
	if len(LogServiceURLs) < 1 {
		fmt.Printf("Logging service not found")
	} else {
		fmt.Printf("Logging service found at URL: %v \n", LogServiceURLs)
		log.SetClientLogger(LogServiceURLs[0], serviceName)
		stlog.Printf("%v connect log service \n", serviceName)
	}
	<-ctx.Done()
	fmt.Println("Shutting down portal.")
}
