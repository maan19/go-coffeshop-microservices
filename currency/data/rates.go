package data

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
)

type EchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

func NewEchangeRates(log hclog.Logger) (*EchangeRates, error) {
	er := &EchangeRates{
		log:   log,
		rates: make(map[string]float64),
	}
	err := er.getRates()
	if err != nil {
		return nil, err
	}
	return er, nil
}

func (e *EchangeRates) GetRates(base, dest string) (float64, error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("base rate not found: %s", base)
	}

	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("dest rate not found: %s", dest)
	}
	return dr / br, nil
}

func (e *EchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	ret := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				for k, v := range e.rates {

					change := (rand.Float64() / 10)
					direction := rand.Intn(1)

					if direction == 0 {
						change = 1 - change
					} else {
						change = 1 + change
					}
					e.rates[k] = v * change
				}

				ret <- struct{}{}
			}
		}
	}()

	return ret
}

func (e *EchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 but got code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	err = xml.NewDecoder(resp.Body).Decode(md)
	if err != nil {
		return err
	}

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}

		e.rates[c.Currency] = r
	}
	e.rates["EUR"] = 1
	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
