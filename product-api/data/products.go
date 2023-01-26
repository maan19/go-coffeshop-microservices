package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/maan19/go-coffeshop-microservices/currency/protos/currency/pb"
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
	Price float64 `json:"price" validate:"gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`
}

var ErrProductNotFound = fmt.Errorf("Product not found")

type Products []*Product

type ProductsDB struct {
	currency pb.CurrencyClient
	log      hclog.Logger
}

func NewProductsDB(currency pb.CurrencyClient, logger hclog.Logger) *ProductsDB {
	return &ProductsDB{currency: currency, log: logger}
}

// return all products from database
func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	//get rate
	p.log.Info("Currency", currency)
	rate, err := p.getRate(currency)
	if err != nil {
		return nil, err
	}
	p.log.Info("Rate is", rate)
	pr := Products{}
	for _, p := range productList {
		np := *p
		np.Price = np.Price * rate
		pr = append(pr, &np)
	}

	return pr, nil
}

// return a single product with matching ID
func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	if currency == "" {
		return productList[i], nil
	}

	//get rate
	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("error getting rate", err)
		return nil, err
	}

	np := *productList[i]
	np.Price = np.Price * rate
	return &np, nil
}

// Adds a new Product to the database
func (p *ProductsDB) AddProduct(pr Product) {
	lp := productList[len(productList)-1]
	pr.ID = lp.ID + 1
	productList = append(productList, &pr)
}

// Update replaces the product in the databse witha given product.
func (p *ProductsDB) UpdateProduct(pr Product) error {
	i := findIndexByProductID(pr.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	productList[i] = &pr
	return nil
}

// Delete deletes the product from the databse.
func (p *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1:]...)
	return nil
}

// call currency service to get rate
func (p *ProductsDB) getRate(destination string) (float64, error) {
	//call the currency-client
	rr := &pb.RateRequest{
		Base:        pb.Currencies(pb.Currencies_value["EUR"]),
		Destination: pb.Currencies(pb.Currencies_value[destination]),
	}
	rs, err := p.currency.GetRate(context.Background(), rr)
	p.log.Info("Rate for currency is", destination, rs.Rate)
	return rs.Rate, err
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
