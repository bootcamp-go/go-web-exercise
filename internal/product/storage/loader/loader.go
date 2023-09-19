package loader

import (
	"app/internal/product/storage"
	"errors"
)

var (
	// ErrLoaderProductInternal is an error that returns when an internal error occurs
	ErrLoaderProductInternal = errors.New("loader: internal error")
)

// ProductsDB is a struct that contains the information of products db
type ProductsDB struct {
	Db     map[int]storage.ProductAttributesMap
	LastId int
}

// Loader is an interface that contains a method that loads the products into map
type Loader interface {
	// Load is a method that loads the products into memory
	Load() (p *ProductsDB, err error)
}