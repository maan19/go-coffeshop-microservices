package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{Name: "Jag", Price: 1.23, SKU: "ad-qw-as"}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
