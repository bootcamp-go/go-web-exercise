package storage

import (
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

// StorageProductMap is a struct that contains the information of a storage product
type StorageProductMap struct {
	db	   map[int]*ProductAttributesMap
	lastId int
}

// NewStorageProductMap is a method that creates a new storage product
func NewStorageProductMap(db map[int]*ProductAttributesMap, lastId int) *StorageProductMap {
	return &StorageProductMap{db, lastId}
}

// Get is a method that returns all products
func (s *StorageProductMap) Get() (p []*Product, err error) {
	p = make([]*Product, 0, len(s.db))
	for k, v := range s.db {
		// serialization
		p = append(p, &Product{k, v.Name, v.Quantity, v.CodeValue, v.IsPublished, v.Expiration, v.Price})
	}

	return p, nil
}

// GetByID is a method that returns a product by id
func (s *StorageProductMap) GetByID(id int) (p *Product, err error) {
	product, ok := s.db[id]
	if !ok {
		err = fmt.Errorf("%w: %d", ErrStorageProductNotFound, id)
		return
	}

	// serialization
	p = &Product{id, product.Name, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration, product.Price}
	
	return
}

// Search is a method that returns a product by query
func (s *StorageProductMap) Search(query *Query) (p []*Product, err error) {
	p = make([]*Product, 0)
	for k, v := range s.db {
		// filter: check if query is set
		if query != nil {
			// invalid cases
			if query.Id != 0 && query.Id != k {
				continue
			}
			if query.Name != "" && query.Name != v.Name {
				continue
			}
		}

		// add: by default add the product and serialize
		p = append(p, &Product{k, v.Name, v.Quantity, v.CodeValue, v.IsPublished, v.Expiration, v.Price})
	}

	return p, nil
}

// Create is a method that creates a product
func (s *StorageProductMap) Create(p *Product) (err error) {
	// deserialization
	product := &ProductAttributesMap{p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price}

	// save
	s.lastId++
	s.db[s.lastId] = product
	
	return nil
}
