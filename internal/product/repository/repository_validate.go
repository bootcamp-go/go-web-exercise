package repository

import (
	"app/internal/product"
	"app/internal/product/validator"
	"app/platform/patcher"
	"fmt"
)

// NewRepositoryProductValidate is a method that creates a new repository product validate
func NewRepositoryProductValidate(st RepositoryProduct, vl validator.ValidatorProduct) *RepositoryProductValidate {
	return &RepositoryProductValidate{
		st: st,
		vl: vl,
	}
}

// RepositoryProductValidate is a struct that contains the validates product before repository
type RepositoryProductValidate struct {
	// st is the repository of products
	st RepositoryProduct
	// vl is the validator of products
	vl validator.ValidatorProduct
}


// Get is a method that returns all products
func (s *RepositoryProductValidate) Get() (p []product.Product, err error) {
	p, err = s.st.Get()
	return
}

// GetByID is a method that returns a product by id
func (s *RepositoryProductValidate) GetByID(id int) (p product.Product, err error) {
	p, err = s.st.GetByID(id)
	return
}

// Search is a method that returns a product by query
func (s *RepositoryProductValidate) Search(query Query) (p []product.Product, err error) {
	p, err = s.st.Search(query)
	return
}

// Create is a method that creates a product with validations
func (s *RepositoryProductValidate) Create(p *product.Product) (err error) {
	// validate
	pv := validator.ProductAttributesValidator{
		Name:        p.Name(),
		Quantity:    p.Quantity(),
		CodeValue:   p.CodeValue(),
		IsPublished: p.IsPublished(),
		Expiration:  p.Expiration(),
		Price:       p.Price(),
	}
	err = s.vl.Validate(&pv)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrRepositoryProductInvalid, err.Error())
		return
	}

	// save
	err = s.st.Create(p)
	return
}

// Update is a method that updates a product with validations
func (s *RepositoryProductValidate) Update(id int, patch map[string]any) (p product.Product, err error) {
	// validate
	var pv validator.ProductAttributesValidator
	err = patcher.Patch(&pv, patch)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrRepositoryProductInvalid, err.Error())
		return
	}

	err = s.vl.Validate(&pv)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrRepositoryProductInvalid, err.Error())
		return
	}

	// update
	p, err = s.st.Update(id, patch)
	return
}

// UpdateOrCreate is a method that updates or creates a product with validations
func (s *RepositoryProductValidate) UpdateOrCreate(p *product.Product) (err error) {
	// validate
	pv := validator.ProductAttributesValidator{
		Name:        p.Name(),
		Quantity:    p.Quantity(),
		CodeValue:   p.CodeValue(),
		IsPublished: p.IsPublished(),
		Expiration:  p.Expiration(),
		Price:       p.Price(),
	}
	err = s.vl.Validate(&pv)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrRepositoryProductInvalid, err.Error())
		return
	}

	// update or create
	err = s.st.UpdateOrCreate(p)
	return
}

// Delete is a method that deletes a product by id
func (s *RepositoryProductValidate) Delete(id int) (err error) {
	// delete
	err = s.st.Delete(id)
	return
}