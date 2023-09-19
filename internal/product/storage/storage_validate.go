package storage

import (
	"app/internal/product/validator"
	"fmt"
)

// NewStorageProductValidate is a method that creates a new storage product validate
func NewStorageProductValidate(st StorageProduct, vl validator.ValidatorProduct) *StorageProductValidate {
	return &StorageProductValidate{
		st: st,
		vl: vl,
	}
}

// StorageProductValidate is a struct that contains the validates product before storage
type StorageProductValidate struct {
	// st is the storage of products
	st StorageProduct
	// vl is the validator of products
	vl validator.ValidatorProduct
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
func (s *StorageProductValidate) Search(query Query) (p []*Product, err error) {
	p, err = s.st.Search(query)
	return
}

// Create is a method that creates a product with validations
func (s *StorageProductValidate) Create(p *Product) (err error) {
	// validate
	pv := validator.ProductAttributesValidator{
		Name:        p.Name,
		Quantity:    p.Quantity,
		CodeValue:   p.CodeValue,
		IsPublished: p.IsPublished,
		Expiration:  p.Expiration,
		Price:       p.Price,
	}
	err = s.vl.Validate(&pv)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrStorageProductInvalid, err.Error())
		return
	}

	// save
	err = s.st.Create(p)
	return
}

// Update is a method that updates a product with validations
func (s *StorageProductValidate) Update(p *Product) (err error) {
	// validate
	pv := validator.ProductAttributesValidator{
		Name:        p.Name,
		Quantity:    p.Quantity,
		CodeValue:   p.CodeValue,
		IsPublished: p.IsPublished,
		Expiration:  p.Expiration,
		Price:       p.Price,
	}
	err = s.vl.Validate(&pv)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrStorageProductInvalid, err.Error())
		return
	}

	// save
	err = s.st.Update(p)
	return
}