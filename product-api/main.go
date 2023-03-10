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

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/maan19/go-coffeshop-microservices/currency/protos/currency/pb"
	"github.com/maan19/go-coffeshop-microservices/product-api/data"
	"github.com/maan19/go-coffeshop-microservices/product-api/handlers"
	"google.golang.org/grpc"
)

var bindAddress = flag.String("BIND_ADDRESS", ":9090", "Bind address for the server")

func main() {
	flag.Parse()
	l := hclog.Default()
	v := data.NewValidation()

	//create currency client
	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cc := pb.NewCurrencyClient(conn)

	//new database instance
	db := data.NewProductsDB(cc, l)

	//create produc-api handler
	ph := handlers.NewProducts(l, v, db)

	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/products", ph.ListAll).Queries("currency", "{[A-Z]{3}}")
	getR.HandleFunc("/products", ph.ListAll)

	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle).Queries("currency", "{[A-Z]{3}}")
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/products", ph.Update)
	putR.Use(ph.MiddlewareValidateProduct)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/products", ph.Create)
	postR.Use(ph.MiddlewareValidateProduct)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	s := http.Server{
		Addr:         *bindAddress,
		Handler:      ch(sm),
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Info("Starting server on port 9090")
		if err := s.ListenAndServe(); err != nil {
			l.Error("Error starting server %s\n", err)
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
