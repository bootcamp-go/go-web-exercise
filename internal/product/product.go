package product

import "time"

// NewProduct is a function that creates a new product.
func NewProduct(id int, name string, quantity int, codeValue string, isPublished bool, expiration time.Time, price float64) *Product {
	return &Product{id, name, quantity, codeValue, isPublished, expiration, price}
}

// Product is a struct that represents a product.
type Product struct {
	// id is the unique identifier of the product.
	id int
	// name is the name of the product.
	name string
	// quantity is the quantity of the product.
	quantity int
	// codeValue is the code value of the product. Should be unique.
	codeValue string
	// isPublished is the published state of the product.
	isPublished bool
	// expiration is the expiration date of the product.
	expiration time.Time
	// price is the price of the product.
	price float64
}

// Id returns the id of the product.
func (p *Product) Id() int {
	return p.id
}

// Name returns the name of the product.
func (p *Product) Name() string {
	return p.name
}

// Quantity returns the quantity of the product.
func (p *Product) Quantity() int {
	return p.quantity
}

// CodeValue returns the code value of the product.
func (p *Product) CodeValue() string {
	return p.codeValue
}

// IsPublished returns the published state of the product.
func (p *Product) IsPublished() bool {
	return p.isPublished
}

// Expiration returns the expiration date of the product.
func (p *Product) Expiration() time.Time {
	return p.expiration
}

// Price returns the price of the product.
func (p *Product) Price() float64 {
	return p.price
}

// SetId sets the id of the product.
func (p *Product) SetId(id int) {
	(*p).id = id
}