package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/maan19/go-coffeshop-microservices/currency/data"
	"github.com/maan19/go-coffeshop-microservices/currency/protos/currency/pb"
)

type Currency struct {
	rates *data.EchangeRates
	log   hclog.Logger
}

func NewCurrency(r *data.EchangeRates, l hclog.Logger) *Currency {
	return &Currency{
		log:   l,
		rates: r,
	}
}

func (c *Currency) GetRate(ctx context.Context, rr *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
	rate, err := c.rates.GetRates(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &pb.RateResponse{Rate: rate}, nil
}
