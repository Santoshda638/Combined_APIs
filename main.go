package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// set http handllers routes
	router := newRouter()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Run server in a seperate goroutine so that it doesn't block.
	go func() {
		fmt.Println("Listening at port", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// For gracefull shutdown
	// default wait before shutdown
	wait := time.Second * 10
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// wait until we receive kill signal.
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
