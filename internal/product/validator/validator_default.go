package validator

import (
	"fmt"
	"regexp"
)

// ValidatorProductDefault is a struct that implements the ValidatorProduct interface
type ValidatorProductDefault struct {
	// RegexCodeValue is a regex that represents the code value
	regexCodeValue *regexp.Regexp
}

// NewValidatorProductDefault is a method that creates a new ValidatorProductDefault
func NewValidatorProductDefault(regexCodeValue string) *ValidatorProductDefault {
	// default values
	regexCodeValueDefault := `^[A-Z]{3}-[0-9]{3}$`
	if regexCodeValue != "" {
		regexCodeValueDefault = regexCodeValue
	}

	return &ValidatorProductDefault{
		regexCodeValue: regexp.MustCompile(regexCodeValueDefault),
	}
}

// Validate is a method that validates a product
func (v *ValidatorProductDefault) Validate(p *ProductAttributesValidator) (err error) {
	// validate
	// -> required fields
	if p.Name == "" {
		err = fmt.Errorf("%w: name is empty", ErrValidatorProductFieldRequired)
		return
	}
	if p.CodeValue == "" {
		err = fmt.Errorf("%w: code value is empty", ErrValidatorProductFieldRequired)
		return
	}

	// -> quality fields
	if p.Quantity < 0 {
		err = fmt.Errorf("%w: quantity can't be negative", ErrValidatorProductFieldInvalid)
		return
	}
	if !v.regexCodeValue.MatchString(p.CodeValue) {
		err = fmt.Errorf("%w: code value format is invalid", ErrValidatorProductFieldInvalid)
		return
	}
	if p.Expiration.Before(p.Expiration) {
		err = fmt.Errorf("%w: expiration date can't be before created date", ErrValidatorProductFieldInvalid)
		return
	}
	if p.Price < 0 {
		err = fmt.Errorf("%w: price can't be negative", ErrValidatorProductFieldInvalid)
		return
	}

	return
}
