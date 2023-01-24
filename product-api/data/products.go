package data

import (
	"fmt"
)

// Product defines the structure for an API product
// swagger:model listSingle
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"`

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`
}

var ErrProductNotFound = fmt.Errorf("Product not found")

type Products []*Product

// return all products from database
func GetProducts() Products {
	return productList
}

// return a single product with matching ID
func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}
	return productList[i], nil
}

// Adds a new Product to the database
func AddProduct(p Product) {
	lp := productList[len(productList)-1]
	p.ID = lp.ID + 1
	productList = append(productList, &p)
}

// Update replaces the product in the databse witha given product.
func UpdateProduct(p Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	productList[i] = &p
	return nil
}

// Delete deletes the product from the databse.
func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1:]...)
	return nil
}

func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "short and string coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
