package main

import (
	"context"
	"distributed/log"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("distributed.log")
	host, port := "127.0.0.1", "8888"
	ctx, err := service.Start(
		context.Background(),
		"logService",
		host,
		port,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("shutting down log service.")
}
