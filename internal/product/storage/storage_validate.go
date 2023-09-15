package storage

import (
	"fmt"
	"regexp"
)

// NewStorageProductValidate returns a new StorageProductValidate
type ConfigStorageProductValidate struct {
	// st is the storage of products
	St StorageProduct
	RegexCodeValue string
}

func NewStorageProductValidate(cfg ConfigStorageProductValidate) *StorageProductValidate {
	// default values
	if cfg.RegexCodeValue == "" {
		cfg.RegexCodeValue = `^[A-Z]{3}-[0-9]{3}$`
	}

	return &StorageProductValidate{
		st: cfg.St,
		regexCodeValue: regexp.MustCompile(cfg.RegexCodeValue),
	}
}

// StorageProductValidate is a struct that contains the validates product before storage
type StorageProductValidate struct {
	// st is the storage of products
	st StorageProduct

	// regexCodeValue is the regex pattern for code value
	regexCodeValue *regexp.Regexp
}


// Get is a method that returns all products
func (s *StorageProductValidate) Get() (p []*Product, err error) {
	p, err = s.st.Get()
	return
}

// GetByID is a method that returns a product by id
func (s *StorageProductValidate) GetByID(id int) (p *Product, err error) {
	p, err = s.st.GetByID(id)
	return
}

// Search is a method that returns a product by query
func (s *StorageProductValidate) Search(query *Query) (p []*Product, err error) {
	p, err = s.st.Search(query)
	return
}

// Create is a method that creates a product with validations
func (s *StorageProductValidate) Create(p *Product) (err error) {
	// validate
	// -> required fields
	if p.Name == "" {
		err = fmt.Errorf("%w: name is empty", ErrStorageProductInvalid)
		return
	}
	if p.CodeValue == "" {
		err = fmt.Errorf("%w: code value is empty", ErrStorageProductInvalid)
		return
	}
	// -> quality fields
	if p.Quantity < 0 {
		err = fmt.Errorf("%w: quantity can't be negative", ErrStorageProductInvalid)
		return
	}
	if !s.regexCodeValue.MatchString(p.CodeValue) {
		err = fmt.Errorf("%w: code value format is invalid", ErrStorageProductInvalid)
		return
	}
	if p.Expiration.Before(p.Expiration) {
		err = fmt.Errorf("%w: expiration date can't be before created date", ErrStorageProductInvalid)
		return
	}
	if p.Price < 0 {
		err = fmt.Errorf("%w: price can't be negative", ErrStorageProductInvalid)
		return
	}

	// save
	err = s.st.Create(p)
	return
}
