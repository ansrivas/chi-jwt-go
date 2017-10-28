package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	keyPath string
)

func init() {
	flag.StringVar(&keyPath, "keyPath", "", "path to a key files")
}

func main() {
	flag.Parse()

	if keyPath == "" {
		log.Fatalln("Please pass the path to keys")
	}

	addr := ":8080"
	s := &http.Server{
		Addr:           addr,
		Handler:        NewRouter(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport.
	go func() {
		fmt.Printf("server listening on %s\n", addr)
		errc <- s.ListenAndServe()
	}()
	// Run!
	log.Fatalln("exit", <-errc)
}
