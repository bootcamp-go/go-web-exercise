package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// JSON decodes json from request body to ptr
var (
	ErrRequestJSONInvalid = errors.New("request json invalid")
)

func JSON(r *http.Request, ptr any) (err error) {
	// get body
	err = json.NewDecoder(r.Body).Decode(ptr)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrRequestJSONInvalid, err)
		return
	}

	return
}