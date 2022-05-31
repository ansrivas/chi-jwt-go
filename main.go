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
	addr    string
)

func init() {
	flag.StringVar(&keyPath, "keyPath", "", "path to a key files")
	flag.StringVar(&addr, "addr", ":8080", "Server address to bind to, defaults to :8080")
}

func main() {
	flag.Parse()

	if keyPath == "" {
		flag.Usage()
		return
	}

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
