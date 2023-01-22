// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/maan19/product-api/product-api/data"
)

type Products struct {
	l *log.Logger
	v *data.Validation
}

func NewProducts(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}

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
