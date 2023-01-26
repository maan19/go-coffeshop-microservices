package handlers

import (
	"net/http"

	"github.com/maan19/go-coffeshop-microservices/product-api/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//
//		201: noContentResponse
//	 404: errorResponse
//	 422: errorValidation
func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	p.l.Debug("Handle PUT Product")

	prod := r.Context().Value(KeyProduct).(data.Product)
	err := p.productsDB.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
