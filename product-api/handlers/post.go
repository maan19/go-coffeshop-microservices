package handlers

import (
	"net/http"

	"github.com/maan19/product-api/product-api/data"
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
	p.l.Println("Handle POST Products")
	prod := r.Context().Value(KeyProduct).(data.Product)

	p.l.Printf("[DEBUG] Inserting product: %#v\n", prod)
	data.AddProduct(prod)
	//return success? - that's done automatically
}
