package storage

import (
	"errors"
	"time"
)

var (
	// ErrStorageProductInternal is an error that returns when an internal error occurs
	ErrStorageProductInternal = errors.New("storage: internal error")

	// ErrStorageProductNotFound is an error that returns when a product is not found
	ErrStorageProductNotFound = errors.New("storage: product not found")

	// ErrStorageProductInvalid is an error that returns when a product is invalid
	ErrStorageProductInvalid = errors.New("storage: product invalid")
)

// Product is a struct that contains the information of a product
type Product struct {
	Id			int
	Name        string
	Quantity    int
	CodeValue   string
	IsPublished bool
	Expiration  time.Time
	Price       float64
}

type Query struct {
	Id			int
	Name        string
}

// StorageProduct is an interface that contains the methods that a storage product must implement
type StorageProduct interface {
	// Get is a method that returns all products
	Get() (p []*Product, err error)

	// GetByID is a method that returns a product by id
	GetByID(id int) (p *Product, err error)

	// Search is a method that returns a product by query
	Search(query Query) (p []*Product, err error)

	// Create is a method that creates a product
	Create(p *Product) (err error)
}