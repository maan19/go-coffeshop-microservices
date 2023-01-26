package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/maan19/go-coffeshop-microservices/product-api/data"
)

type Products struct {
	l          hclog.Logger
	v          *data.Validation
	productsDB *data.ProductsDB
}

func NewProducts(l hclog.Logger, v *data.Validation, pdb *data.ProductsDB) *Products {
	return &Products{l, v, pdb}
}

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("invalid Path, path should be /products/[id]")

type contextKey string

const KeyProduct contextKey = contextKey("product")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// get id from request
func getProductID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}
	return id
}
