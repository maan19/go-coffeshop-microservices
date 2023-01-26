package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestNewRates(t *testing.T) {
	tr, err := NewEchangeRates(hclog.Default())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Rates are :%#v", tr.rates)
}
