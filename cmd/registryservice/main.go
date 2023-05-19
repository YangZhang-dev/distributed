package main

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var server http.Server

	server.Addr = registry.ServerPort

	go func() {
		log.Println(server.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("registry service started. Press any key to stop. ")
		var s string
		fmt.Scanln(&s)
		server.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("shutting down registry service.")
}
