package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	// ErrValidatorRequired is an error that returns when a required field is not found
	ErrValidatorRequired = errors.New("validator: required field not found")
)

func RequiredJSON(r io.Reader, keys ...string) (err error) {
	// decode json into map[string]any
	m := make(map[string]any)
	err = json.NewDecoder(r).Decode(&m)
	if err != nil {
		return
	}

	// check if the keys are in the map
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			err = fmt.Errorf("%w - key %s", ErrValidatorRequired, k)
			return
		}
	}

	return
}