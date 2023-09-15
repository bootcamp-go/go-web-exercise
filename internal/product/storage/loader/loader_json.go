package loader

import (
	"app/internal/product/storage"
	"encoding/json"
	"fmt"
	"os"
)

// NewLoaderJSON is a method that creates a new loader json
func NewLoaderJSON(filePath string) *LoaderJSON {
	return &LoaderJSON{filePath}
}

// LoaderJSON is a struct that represents an implementation of the Loader interface for json files
type LoaderJSON struct {
	FilePath string
}

// ProductAttributesMap is a struct that contains the information of a product attributes in json format
type ProductAttributesJSON struct {
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
	CodeValue   string `json:"code_value"`
	IsPublished bool   `json:"is_published"`
	Expiration  string `json:"expiration"`
	Price       string `json:"price"`
}

// Load is a method that loads the products into memory
func (l *LoaderJSON) Load() (p *ProductsDB, err error) {
	// open file
	f, err := os.Open(l.FilePath)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrLoaderProductInternal, err)
		return
	}
	defer f.Close()

	// read file
	// -> get db
	db := make(map[int]*storage.ProductAttributesMap)
	err = json.NewDecoder(f).Decode(&db)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrLoaderProductInternal, err)
		return
	}
	// -> get last id
	lastId := len(db) + 1

	// set products db
	p = &ProductsDB{
		Db:     db,
		LastId: lastId,
	}

	return
}