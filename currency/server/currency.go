package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/maan19/go-coffeshop-microservices/currency/protos/currency/pb"
)

type Currency struct {
	log hclog.Logger
}

func NewCurrency(log hclog.Logger) *Currency {
	return &Currency{
		log: log,
	}
}

func (c *Currency) GetRate(ctx context.Context, rr *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	return &pb.RateResponse{Rate: 0.5}, nil
}
