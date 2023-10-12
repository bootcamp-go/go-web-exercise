package storage

import (
	"app/internal/product"
	"encoding/json"
	"os"
	"time"
)

type StorageProductJSON struct {
	// filePath is the path of the file
	filePath string
	// layoutDate is the layout of the expiration date
	layoutDate string
}

// NewStorageProductJSON is a method that creates a new storage product json
func NewStorageProductJSON(filePath string, layoutDate string) *StorageProductJSON {
	// default config
	defaulLayoutDate := time.DateOnly
	if layoutDate != "" {
		defaulLayoutDate = layoutDate
	}

	return &StorageProductJSON{filePath: filePath, layoutDate: defaulLayoutDate}
}

// ProductAttributesJSON is a struct that contains the information of a product
type ProductAttributesJSON struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

// ReadAll is a method that returns all products
func (s *StorageProductJSON) ReadAll() (p []product.Product, err error) {
	// open file
	f, err := os.Open(s.filePath)
	if err != nil {
		return
	}
	defer f.Close()

	// decode
	pr := make(map[int]ProductAttributesJSON)
	err = json.NewDecoder(f).Decode(&pr)
	if err != nil {
		return
	}

	// serialization
	for k, v := range pr {
		var exp time.Time
		exp, err = time.Parse(s.layoutDate, v.Expiration)
		if err != nil {
			err = ErrStorageProductTimeLayout
			return
		}
		p = append(p, *product.NewProduct(k, v.Name, v.Quantity, v.CodeValue, v.IsPublished, exp, v.Price))
	}

	return
}

// WriteAll is a method that writes all products
func (s *StorageProductJSON) WriteAll(p []product.Product) (err error) {
	// open file
	f, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}

	// serialization
	pr := make(map[int]ProductAttributesJSON)
	for _, v := range p {
		pr[v.Id()] = ProductAttributesJSON{v.Name(), v.Quantity(), v.CodeValue(), v.IsPublished(), v.Expiration().Format(s.layoutDate), v.Price()}
	}
	
	// encode
	err = json.NewEncoder(f).Encode(pr)
	if err != nil {
		return
	}

	return
}
