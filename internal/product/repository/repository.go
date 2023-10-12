package repository

import (
	"app/internal/product"
	"errors"
)

var (
	// ErrRepositoryProductInternal is an error that returns when an internal error occurs
	ErrRepositoryProductInternal = errors.New("repository: internal error")

	// ErrRepositoryProductNotFound is an error that returns when a product is not found
	ErrRepositoryProductNotFound = errors.New("repository: product not found")

	// ErrRepositoryProductInvalid is an error that returns when a product is invalid
	ErrRepositoryProductInvalid = errors.New("repository: product invalid")

	// ErrRepositoryProductPatchInvalid is an error that returns when a product patch is invalid
	ErrRepositoryProductPatchInvalid = errors.New("repository: product patch invalid")
)

type Query struct {
	Id			int
	Name        string
}

// RepositoryProduct is an interface that contains the methods that a repository product must implement
type RepositoryProduct interface {
	// Get is a method that returns all products
	Get() (p []product.Product, err error)

	// GetByID is a method that returns a product by id
	GetByID(id int) (p product.Product, err error)

	// Search is a method that returns a product by query
	Search(query Query) (p []product.Product, err error)

	// Create is a method that creates a product
	Create(p *product.Product) (err error)

	// UpdateOrCreate is a method that updates or creates a product
	UpdateOrCreate(p *product.Product) (err error)

	// Update is a method that updates a product
	Update(id int, patch map[string]any) (p product.Product, err error)

	// Delete is a method that deletes a product by id
	Delete(id int) (err error)
}