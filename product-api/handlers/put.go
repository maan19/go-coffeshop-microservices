package handlers

import (
	"net/http"

	"github.com/maan19/product-api/product-api/data"
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
	p.l.Println("Handle PUT Product")

	prod := r.Context().Value(KeyProduct).(data.Product)
	err := data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
