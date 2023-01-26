package handlers

import (
	"net/http"

	"github.com/maan19/go-coffeshop-microservices/product-api/data"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// Responses:
//	 200: productsResponse

// Returns all products from db
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("get all records")

	rw.Header().Add("Content-Type", "application/json")

	cur := r.URL.Query().Get("currency")

	prods, err := p.productsDB.GetProducts(cur)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	err = data.ToJSON(prods, rw)
	if err != nil {
		p.l.Error("error serializing products", err)
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a product from the database
// Responses:
//	 200: productResponse
//	 404: errorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)
	cur := r.URL.Query().Get("currency")

	p.l.Debug("get single record for id ", id)

	prod, err := p.productsDB.GetProductByID(id, cur)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("product not found")
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return

	default:
		p.l.Error("error fetching product", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		p.l.Error("error serializing product", err)
	}
}
