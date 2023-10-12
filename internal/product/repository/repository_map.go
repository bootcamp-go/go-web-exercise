package repository

import (
	"app/internal/product"
	"app/platform/patcher"
	"fmt"
	"time"
)

// ProductAttributesMap is a struct that contains the information of a product
type ProductAttributesMap struct {
	Name        string
	Quantity    int
	CodeValue   string
	IsPublished bool
	Expiration  time.Time
	Price       float64
}

// RepositoryProductMap is a struct that contains the information of a storage product
type RepositoryProductMap struct {
	// db is the database of the repository
	db	   map[int]ProductAttributesMap
	// lastId is the last id of the repository
	lastId int
	// layoutDate is the layout date of the repository
	layoutDate string
}

// NewRepositoryProductMap is a method that creates a new storage product
func NewRepositoryProductMap(db map[int]ProductAttributesMap, lastId int, layoutDate string) *RepositoryProductMap {
	// default config
	defaultLayoutDate := time.DateOnly
	if layoutDate == "" {
		layoutDate = defaultLayoutDate
	}
	return &RepositoryProductMap{db: db, lastId: lastId, layoutDate: defaultLayoutDate}
}

// Get is a method that returns all products
func (s *RepositoryProductMap) Get() (p []product.Product, err error) {
	p = make([]product.Product, 0, len(s.db))
	for k, v := range s.db {
		// serialization
		p = append(p, *product.NewProduct(k, v.Name, v.Quantity, v.CodeValue, v.IsPublished, v.Expiration, v.Price))
	}

	return p, nil
}

// GetByID is a method that returns a product by id
func (s *RepositoryProductMap) GetByID(id int) (p product.Product, err error) {
	pr, ok := s.db[id]
	if !ok {
		err = fmt.Errorf("%w: %d", ErrRepositoryProductNotFound, id)
		return
	}

	// serialization
	p = *product.NewProduct(id, pr.Name, pr.Quantity, pr.CodeValue, pr.IsPublished, pr.Expiration, pr.Price)
	
	return
}

// Search is a method that returns filtered products
// valid if at least one field is set
func (s *RepositoryProductMap) Search(query Query) (p []product.Product, err error) {
	valid := query.Id > 0 || query.Name != ""
	
	// filter
	for k, v := range s.db {
		// check if query is valid
		if valid {
			// check if id is valid
			if query.Id > 0 && k != query.Id {
				continue
			}
			// check if name is valid
			if query.Name != "" && v.Name != query.Name {
				continue
			}
		}

		// serialization
		p = append(p, *product.NewProduct(k, v.Name, v.Quantity, v.CodeValue, v.IsPublished, v.Expiration, v.Price))
	}

	return
}

// Create is a method that creates a product
func (s *RepositoryProductMap) Create(p *product.Product) (err error) {
	// deserialization
	pr := ProductAttributesMap{p.Name(), p.Quantity(), p.CodeValue(), p.IsPublished(), p.Expiration(), p.Price()}

	// save
	s.lastId++
	s.db[s.lastId] = pr
	
	// last id
	p.SetId(s.lastId)

	return nil
}

// ProductAttributesPatch is a struct that contains the information of a product
type ProductAttributesPatch struct {
	Name		string		`patcher:"name"`
	Quantity	int			`patcher:"quantity"`
	CodeValue	string		`patcher:"code_value"`
	IsPublished	bool		`patcher:"is_published"`
	Expiration	string		`patcher:"expiration"`
	Price		float64		`patcher:"price"`
}
// Update is a method that updates a product
// - this method handles dynamic patching
func (s *RepositoryProductMap) Update(id int, patch map[string]any) (p product.Product, err error) {
	// search
	pr, ok := s.db[id]
	if !ok {
		err = fmt.Errorf("%w - %d", ErrRepositoryProductNotFound, id)
		return
	}

	// deserialization
	patchAttributes := ProductAttributesPatch{
		Name:        pr.Name,
		Quantity:    pr.Quantity,
		CodeValue:   pr.CodeValue,
		IsPublished: pr.IsPublished,
		Expiration:  pr.Expiration.Format(s.layoutDate),
		Price:       pr.Price,
	}
	
	// patch
	err = patcher.Patch(&patchAttributes, patch)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrRepositoryProductInvalid, err.Error())
		return
	}

	// serialization
	fieldExpiration, err := time.Parse(s.layoutDate, patchAttributes.Expiration)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrRepositoryProductInvalid, err.Error())
		return
	}
	pr = ProductAttributesMap{
		Name:        patchAttributes.Name,
		Quantity:    patchAttributes.Quantity,
		CodeValue:   patchAttributes.CodeValue,
		IsPublished: patchAttributes.IsPublished,
		Expiration:  fieldExpiration,
		Price:       patchAttributes.Price,
	}

	// save
	s.db[id] = pr

	// serialization
	p = *product.NewProduct(id, pr.Name, pr.Quantity, pr.CodeValue, pr.IsPublished, pr.Expiration, pr.Price)
	
	return
}



// UpdateOrCreate is a method that updates or creates a product
func (s *RepositoryProductMap) UpdateOrCreate(p *product.Product) (err error) {
	// deserialization
	pr := ProductAttributesMap{p.Name(), p.Quantity(), p.CodeValue(), p.IsPublished(), p.Expiration(), p.Price()}

	// update or create
	_, ok := s.db[p.Id()]
	switch ok {
	case true:
		// update
		s.db[p.Id()] = pr
	default:
		// save
		s.lastId++
		s.db[s.lastId] = pr
	}
	
	return
}

// Delete is a method that deletes a product by id
func (s *RepositoryProductMap) Delete(id int) (err error) {
	// search
	_, ok := s.db[id]
	if !ok {
		err = fmt.Errorf("%w: %d", ErrRepositoryProductNotFound, id)
		return
	}

	// delete
	delete(s.db, id)
	
	return
}