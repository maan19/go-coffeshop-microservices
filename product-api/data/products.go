package data

import (
	"fmt"
)

var ErrProductNotFound = fmt.Errorf("Product not found")

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"sku"`
}

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
