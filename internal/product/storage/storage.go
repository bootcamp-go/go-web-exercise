package storage

import (
	"app/internal/product"
	"errors"
)

var (
	// ErrStorageProductTimeLayout is an error that returns when the time layout is invalid
	ErrStorageProductTimeLayout = errors.New("storage: time layout invalid")
)

// StorageProduct is an interface that contains the methods that a storage product must implement
type StorageProduct interface {
	// ReadAll is a method that returns all products
	ReadAll() (p []product.Product, err error)

	// WriteAll is a method that writes all products
	WriteAll(p []product.Product) (err error)
}