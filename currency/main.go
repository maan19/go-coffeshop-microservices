package main

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/maan19/go-coffeshop-microservices/currency/data"
	"github.com/maan19/go-coffeshop-microservices/currency/protos/currency/pb"
	"github.com/maan19/go-coffeshop-microservices/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	rates, err := data.NewEchangeRates(log)
	if err != nil {
		log.Error("Unable to get rates", err)
	}
	gs := grpc.NewServer()
	cs := server.NewCurrency(rates, log)

	pb.RegisterCurrencyServer(gs, cs)

	//Disbale in production
	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("failed to listen: %v", err)
		os.Exit(1)
	}

	err = gs.Serve(l)
	if err != nil {
		log.Error("failed to serve: %v", err)
		os.Exit(1)
	}
}
