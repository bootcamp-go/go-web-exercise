package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
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

// PathLastParam returns the value of the last path parameter
var (
	ErrRequestPathParamInvalid = errors.New("request path param invalid")
)
func PathLastParam(r *http.Request) (value string, err error) {
	// get url path (example: /api/v1/products/123) (query not included such as ?page=1)
	path := r.URL.Path

	// check path matches regexp pattern
	rx := regexp.MustCompile(`^/(.*/)*([0-9a-zA-Z]+)$`)
	if !rx.MatchString(path) {
		err = ErrRequestPathParamInvalid
		return
	}

	// split path by slash
	sl := strings.Split(path, "/")
	value = sl[len(sl)-1]
	return
}