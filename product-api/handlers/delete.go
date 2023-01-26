package handlers

import (
	"net/http"

	"github.com/maan19/go-coffeshop-microservices/product-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Update a products details
//
// responses:
//
//	201: noContentResponse
//	404: errorResponse
//	501: errorResponse
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	p.l.Debug("deleting record id", id)

	err := p.productsDB.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("error deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	rw.WriteHeader(http.StatusNoContent)

}
