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
	p.l.Println("[DEBUG] get all records")

	prods := data.GetProducts()

	err := data.ToJSON(prods, rw)
	if err != nil {
		p.l.Println(err)
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a product from the database
// Responses:
//	 200: productResponse
//	 404: errorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)
	p.l.Println("[DEBUG] get single record for id ", id)

	prod, err := data.GetProductByID(id)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Println("[ERROR] product not found")
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return

	default:
		p.l.Println("[ERROR] fetching product", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	err = data.ToJSON(prod, rw)
	if err != nil {
		p.l.Println(err)
	}
}
