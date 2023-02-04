package server

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/maan19/go-coffeshop-microservices/currency/data"
	"github.com/maan19/go-coffeshop-microservices/currency/protos/currency/pb"
)

type Currency struct {
	rates *data.EchangeRates
	log   hclog.Logger
}

func NewCurrency(r *data.EchangeRates, l hclog.Logger) *Currency {
	go func() {
		ru := r.MonitorRates(5 * time.Second)
		for range ru {
			l.Info("got updated rates")
		}
	}()
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

func (c *Currency) SubscribeRates(src pb.Currency_SubscribeRatesServer) error {
	//handle client messages
	go func() {
		for {
			rr, err := src.Recv()
			if err == io.EOF {
				c.log.Info("Client closed connection")
				break
			}
			if err != nil { //Connection forcibly closed.
				c.log.Error("Receive error", "err", err)
				break
			}
			c.log.Info("Client request received", "request", rr)
		}
	}()

	//handler server responses
	for {
		err := src.Send(&pb.RateResponse{Rate: 12.1})
		if err != nil {
			c.log.Error("Send error", "err", err)
			return err
		}
		time.Sleep(5 * time.Second)
	}
}
