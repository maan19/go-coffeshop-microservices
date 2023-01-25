package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/maan19/go-coffeshop-microservices/product-images/files"
	"github.com/maan19/go-coffeshop-microservices/product-images/handlers"
)

var bindAddress = flag.String("BIND_ADDRESS", ":9091", "server port")
var logLevel = flag.String("LOG_LEVEL", "debug", "log-level for server- [debug|info|trace]")
var basePath = flag.String("BASE_PATH", "./imagestore", "path to store images")

func main() {
	flag.Parse()

	fmt.Printf("bindAddress: %s\n", *bindAddress)
	fmt.Printf("logLevel: %s\n", *logLevel)
	fmt.Printf("basePath: %s\n", *basePath)
	l := hclog.New(&hclog.LoggerOptions{
		Name:  "product-image",
		Level: hclog.LevelFromString(*logLevel),
	})
	//create a standard logger for the server
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// File store instance
	stor, err := files.NewLocal(*basePath)
	if err != nil {
		panic(err)
	}

	//Create handlers and middlewares
	fh := handlers.NewFiles(l, stor)
	gzipm := handlers.GzipHandler{}

	//Mux router
	sm := mux.NewRouter()

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	//post image route
	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.UploadREST)
	ph.HandleFunc("/", fh.UploadMultipart)

	//get image route
	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir("./imagestore"))))

	//use gzipmiddleware on get image router
	gh.Use(gzipm.GzipMiddleware)

	//create a new server
	s := http.Server{
		Addr:         *bindAddress,
		Handler:      ch(sm),
		ErrorLog:     sl,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//stat the server
	go func() {
		l.Info("starting server on", *bindAddress)
		err = s.ListenAndServe()
		if err != nil {
			l.Error("Error starting the server", err)
			os.Exit(1)
		}
	}()

	//wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	//Block until a signal is received
	sig := <-quit
	l.Info("Shutting down the server got signal:", sig)

	//graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	if err != nil {
		l.Error("Error shutting down the server", err)
	}
}
