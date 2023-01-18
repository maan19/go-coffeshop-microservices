package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maan19/product-api/product-api/handlers"
)

var bindAddress = flag.String("BIND_ADDRESS", ":9090", "Bind address for the server")

func main() {
	flag.Parse()
	l := log.New(os.Stdout, "products-api", log.LstdFlags)
	ph := handlers.NewProducts(l)
	sm := http.NewServeMux()
	sm.Handle("/", ph)
	s := http.Server{
		Addr:         *bindAddress,
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9090")
		if err := s.ListenAndServe(); err != nil {
			l.Printf("Error starting server %s\n", err)
			os.Exit(1)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	log.Println("Got signal:", sig)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
