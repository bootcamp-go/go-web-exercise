package repository

import "app/internal/product"

// NewRepositoryProductMock creates a new mock of the RepositoryProduct interface
func NewRepositoryProductMock() *RepositoryProductMock {
	return &RepositoryProductMock{}
}

// RepositoryProductMock is a mock of the RepositoryProduct interface
type RepositoryProductMock struct {
	// FuncGet is a function to proxy the Get method
	FuncGet            func() (p []product.Product, err error)
	// FuncGetByID is a function to proxy the GetByID method
	FuncGetByID        func(id int) (p product.Product, err error)
	// FuncSearch is a function to proxy the Search method
	FuncSearch         func(query Query) (p []product.Product, err error)
	// FuncCreate is a function to proxy the Create method
	FuncCreate         func(p *product.Product) (err error)
	// FuncUpdateOrCreate is a function to proxy the UpdateOrCreate method
	FuncUpdateOrCreate func(p *product.Product) (err error)
	// FuncUpdate is a function to proxy the Update method
	FuncUpdate         func(id int, patch map[string]any) (p product.Product, err error)
	// FuncDelete is a function to proxy the Delete method
	FuncDelete         func(id int) (err error)

	// Spy
	Spy struct {
		// GetCalls is a counter of the Get method
		GetCalls            int
		// GetByIDCalls is a counter of the GetByID method
		GetByIDCalls        int
		// SearchCalls is a counter of the Search method
		SearchCalls         int
		// CreateCalls is a counter of the Create method
		CreateCalls         int
		// UpdateOrCreateCalls is a counter of the UpdateOrCreate method
		UpdateOrCreateCalls int
		// UpdateCalls is a counter of the Update method
		UpdateCalls         int
		// DeleteCalls is a counter of the Delete method
		DeleteCalls         int
	}
}

// Get is a method that returns all products
func (s *RepositoryProductMock) Get() (p []product.Product, err error) {
	// increment calls counter
	s.Spy.GetCalls++

	// call mock function
	p, err = s.FuncGet()
	return
}

// GetByID is a method that returns a product by id
func (s *RepositoryProductMock) GetByID(id int) (p product.Product, err error) {
	// increment calls counter
	s.Spy.GetByIDCalls++

	// call mock function
	p, err = s.FuncGetByID(id)
	return
}

// Search is a method that returns a product by query
func (s *RepositoryProductMock) Search(query Query) (p []product.Product, err error) {
	// increment calls counter
	s.Spy.SearchCalls++

	// call mock function
	p, err = s.FuncSearch(query)
	return
}

// Create is a method that creates a product
func (s *RepositoryProductMock) Create(p *product.Product) (err error) {
	// increment calls counter
	s.Spy.CreateCalls++

	// call mock function
	err = s.FuncCreate(p)
	return
}

// UpdateOrCreate is a method that updates or creates a product
func (s *RepositoryProductMock) UpdateOrCreate(p *product.Product) (err error) {
	// increment calls counter
	s.Spy.UpdateOrCreateCalls++

	// call mock function
	err = s.FuncUpdateOrCreate(p)
	return
}

// Update is a method that updates a product
func (s *RepositoryProductMock) Update(id int, patch map[string]any) (p product.Product, err error) {
	// increment calls counter
	s.Spy.UpdateCalls++

	// call mock function
	p, err = s.FuncUpdate(id, patch)
	return
}

// Delete is a method that deletes a product by id
func (s *RepositoryProductMock) Delete(id int) (err error) {
	// increment calls counter
	s.Spy.DeleteCalls++

	// call mock function
	err = s.FuncDelete(id)
	return
}