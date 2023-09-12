package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
)

// ________________________________________________________________________________
// storage.go
// Product is a struct that contains the information of a product
type Product struct {
	Name        string
	Quantity    int
	CodeValue   string
	IsPublished bool
	Expiration  time.Time
	Price       float64
}

// ________________________________________________________________________________
// storage.go
// - validate
var (
	ErrValidateRequiredField = errors.New("validate: required field")
	ErrValidateQualityField  = errors.New("validate: quality field")
)
func Validate(p *Product) (err error) {
	// required fields
	if p.Name == "" {
		err = fmt.Errorf("%w: name", ErrValidateRequiredField)
		return
	}
	if p.CodeValue == "" {
		err = fmt.Errorf("%w: code_value", ErrValidateRequiredField)
		return
	}
	if p.Expiration.IsZero() {
		err = fmt.Errorf("%w: expiration", ErrValidateRequiredField)
		return
	}

	// quality fields
	if p.Quantity < 0 {
		err = fmt.Errorf("%w: quantity", ErrValidateQualityField)
		return
	}
	rx := regexp.MustCompile(`^[A-Z]{3}-[0-9]{3}$`)
	if !rx.MatchString(p.CodeValue) {
		err = fmt.Errorf("%w: code_value", ErrValidateQualityField)
		return
	}
	if p.Expiration.Before(time.Now()) {
		err = fmt.Errorf("%w: expiration", ErrValidateQualityField)
		return
	}
	if p.Price < 0 {
		err = fmt.Errorf("%w: price", ErrValidateQualityField)
		return
	}

	return
}

// ________________________________________________________________________________
// LoaderProducts load products from a json file
type ProductAttributesJSON struct {
	Name        string    `json:"name"`
	Quantity    int       `json:"quantity"`
	CodeValue   string    `json:"code_value"`
	IsPublished bool      `json:"is_published"`
	Expiration  time.Time `json:"expiration"`
	Price       float64   `json:"price"`
}

func LoaderProducts(filePath string) (p map[int]*Product, err error) {
	// open file
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	// read file: decode json
	var products map[int]*ProductAttributesJSON
	err = json.NewDecoder(f).Decode(&products)
	if err != nil {
		return
	}

	// serialize
	for k, v := range products {
		p[k] = &Product{
			Name:        v.Name,
			Quantity:    v.Quantity,
			CodeValue:   v.CodeValue,
			IsPublished: v.IsPublished,
			Expiration:  v.Expiration,
			Price:       v.Price,
		}
	}

	return
}
