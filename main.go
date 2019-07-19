package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func main() {

	handller := NewHandller()
	port := "8080"
	mux := http.NewServeMux()
	mux.HandleFunc("/GetNewJoke", handller.GetNewJokeHandller)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	// Lets build the cache for names and jokes.
	go buildCache()

	// Run server in a seperate goroutine so that it doesn't block.
	go func() {
		fmt.Println("Listening at port :", port)
		log.Fatal(srv.ListenAndServe())
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
