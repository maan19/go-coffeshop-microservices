package handlers

import (
	"net/http"

	"github.com/maan19/go-coffeshop-microservices/product-api/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//
//		200: productResponse
//	 422: errorValidation
//	 501: errorResponse
func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	p.l.Debug("Handle POST Products")
	prod := r.Context().Value(KeyProduct).(data.Product)

	p.l.Debug("Inserting product: %#v\n", prod)
	p.productsDB.AddProduct(prod)
	//return success? - that's done automatically
}
