package data

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductMissingNameReturnsErr(t *testing.T) {
	product := Product{
		Price: 1.23,
		SKU:   "abc-abc-abc",
	}

	v := NewValidation()
	err := v.Validate(product)
	assert.Len(t, err, 1)
}

func TestProductMissingPriceReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: -1,
		SKU:   "abc-abc-abc",
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

func TestProductInvalidSKUReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.22,
		SKU:   "abc",
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

func TestValidProductDoesNOTReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.22,
		SKU:   "abc-efg-hji",
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 0)
}

func TestProductsToJSON(t *testing.T) {
	ps := []*Product{
		&Product{
			Name: "abc",
		},
	}

	b := bytes.NewBufferString("")
	err := ToJSON(ps, b)
	assert.NoError(t, err)
}
