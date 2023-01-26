package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

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
	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
