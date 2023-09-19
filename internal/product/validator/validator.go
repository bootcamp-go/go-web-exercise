package validator

import (
	"errors"
	"time"
)

var (
	// ErrValidatorProductFieldRequired is an error that returns when a product field is required
	ErrValidatorProductFieldRequired = errors.New("validator: product field required")
	
	// ErrValidatorProductFieldInvalid is an error that returns when a product field is invalid
	ErrValidatorProductFieldInvalid = errors.New("validator: product field invalid")
)

// ProductAttributesValidator is a struct that contains the information of a product
type ProductAttributesValidator struct {
	Name        string
	Quantity    int
	CodeValue   string
	IsPublished bool
	Expiration  time.Time
	Price       float64
}

// ValidatorProduct is an interface that contains the methods that a validator product must implement
type ValidatorProduct interface {
	// Validate is a method that validates a product
	Validate(p *ProductAttributesValidator) (err error)
}