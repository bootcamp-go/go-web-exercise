package repository

import (
	"app/internal/product"
	"app/internal/product/storage"
	"app/platform/patcher"
	"time"
)

type RepositoryProductStore struct {
	// store is the store of the repository
	st storage.StorageProduct
	// lastID is the last id of the repository
	lastID int
	// layoutDate is the layout date of the repository
	layoutDate string
}

// NewRepositoryProductStore is a method that creates a new repository product store
func NewRepositoryProductStore(st storage.StorageProduct, lastId int, layoutDate string) *RepositoryProductStore {
	// default config
	defaultLayoutDate := time.DateOnly
	if layoutDate == "" {
		layoutDate = defaultLayoutDate
	}
	return &RepositoryProductStore{
		st:         st,
		lastID:     lastId,
		layoutDate: defaultLayoutDate,
	}
}

// Get is a method that returns all products
func (r *RepositoryProductStore) Get() (p []product.Product, err error) {
	// get all products
	p, err = r.st.ReadAll()
	return
}

// GetByID is a method that returns a product by id
func (r *RepositoryProductStore) GetByID(id int) (p product.Product, err error) {
	// get all products
	ps, err := r.st.ReadAll()
	if err != nil {
		return
	}

	// search product
	var exists bool
	for _, v := range ps {
		if v.Id() == id {
			p = v
			exists = true
			break
		}
	}

	// check if product exists
	if !exists {
		err = ErrRepositoryProductNotFound
		return
	}

	return
}

// Search is a method that returns a product by query
func (r *RepositoryProductStore) Search(query Query) (p []product.Product, err error) {
	// get all products
	ps, err := r.st.ReadAll()
	if err != nil {
		return
	}

	// search product
	for _, v := range ps {
		if query.Id != 0 && v.Id() != query.Id {
			continue
		}
		if query.Name != "" && v.Name() != query.Name {
			continue
		}
		p = append(p, v)
	}

	return
}

// Create is a method that creates a product
func (r *RepositoryProductStore) Create(p *product.Product) (err error) {
	// get all products
	ps, err := r.st.ReadAll()
	if err != nil {
		return
	}

	// validation - consistency
	// ...

	// set id
	r.lastID++
	p.SetId(r.lastID)

	// append product
	ps = append(ps, *p)

	// write all products
	err = r.st.WriteAll(ps)
	if err != nil {
		return
	}

	return
}

// UpdateOrCreate is a method that updates or creates a product
func (r *RepositoryProductStore) UpdateOrCreate(p *product.Product) (err error) {
	// get all products
	ps, err := r.st.ReadAll()
	if err != nil {
		return
	}

	// validation - consistency
	// ...

	// search product
	var exists bool; var ix int
	for k, v := range ps {
		if v.Id() == p.Id() {
			ix = k
			exists = true
			break
		}
	}

	// update or create product
	switch exists {
	case true:
		ps[ix] = *p
	default:
		// set id
		r.lastID++
		p.SetId(r.lastID)
		ps = append(ps, *p)
	}

	// write all products
	err = r.st.WriteAll(ps)
	if err != nil {
		return
	}

	return
}

// Update is a method that updates a product
func (r *RepositoryProductStore) Update(id int, patch map[string]interface{}) (p product.Product, err error) {
	// get all products
	ps, err := r.st.ReadAll()
	if err != nil {
		return
	}

	// search product
	var exists bool; var ix int
	for k, v := range ps {
		if v.Id() == id {
			ix = k
			exists = true
			break
		}
	}

	// check if product exists
	if !exists {
		err = ErrRepositoryProductNotFound
		return
	}

	// set product
	p = ps[ix]

	// deserialize patch
	pr := ProductAttributesPatch{
		Name:        p.Name(),
		Quantity:    p.Quantity(),
		CodeValue:   p.CodeValue(),
		IsPublished: p.IsPublished(),
		Expiration:  p.Expiration().Format(r.layoutDate),
		Price:       p.Price(),
	}
	err = patcher.Patch(&pr, patch)
	if err != nil {
		err = ErrRepositoryProductPatchInvalid
		return
	}

	// serialize product
	exp, err := time.Parse(r.layoutDate, pr.Expiration)
	if err != nil {
		err = ErrRepositoryProductPatchInvalid
		return
	}
	p = *product.NewProduct(
		p.Id(),
		pr.Name,
		pr.Quantity,
		pr.CodeValue,
		pr.IsPublished,
		exp,
		pr.Price,
	)

	// update product
	ps[ix] = p

	// write all products
	err = r.st.WriteAll(ps)
	if err != nil {
		return
	}

	return
}

// Delete is a method that deletes a product by id
func (r *RepositoryProductStore) Delete(id int) (err error) {
	// get all products
	ps, err := r.st.ReadAll()
	if err != nil {
		return
	}

	// search product
	var exists bool; var ix int
	for k, v := range ps {
		if v.Id() == id {
			ix = k
			exists = true
			break
		}
	}

	// check if product exists
	if !exists {
		err = ErrRepositoryProductNotFound
		return
	}

	// delete product
	ps = append(ps[:ix], ps[ix+1:]...)

	// write all products
	err = r.st.WriteAll(ps)
	if err != nil {
		return
	}

	return
}