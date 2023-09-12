package request

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for PathParam function
func TestPathParam(t *testing.T) {
	type input struct { r *http.Request }
	type output struct { value string; err error; errMsg string }
	type testCase struct {
		name string
		input input
		output output
	}

	cases := []testCase{
		// valid cases
		{
			name: "valid path param",
			input: input{r: &http.Request{URL: &url.URL{Path: "/api/v1/products/123"}}},
			output: output{value: "123", err: nil, errMsg: ""},
		},
		{
			name: "valid path, only 2 slashes",
			input: input{r: &http.Request{URL: &url.URL{Path: "/clients/123"}}},
			output: output{value: "123", err: nil, errMsg: ""},
		},

		// invalid cases
		{
			name: "invalid path, ends with slash",
			input: input{r: &http.Request{URL: &url.URL{Path: "/api/v1/products/123/"}}},
			output: output{value: "", err: ErrRequestPathParamInvalid, errMsg: "request path param invalid"},
		},
		{
			name: "invalid path, does not start with slash",
			input: input{r: &http.Request{URL: &url.URL{Path: "api/v1/products/123"}}},
			output: output{value: "", err: ErrRequestPathParamInvalid, errMsg: "request path param invalid"},
		},
		{
			name: "invalid path, multiple slashes in the end",
			input: input{r: &http.Request{URL: &url.URL{Path: "/api/v1/products/123///"}}},
			output: output{value: "", err: ErrRequestPathParamInvalid, errMsg: "request path param invalid"},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			// ...

			// act
			value, err := PathLastParam(c.input.r)

			// assert
			require.Equal(t, c.output.value, value)
			require.Equal(t, c.output.err, err)
			if c.output.err != nil {
				require.EqualError(t, err, c.output.errMsg)
			}
		})
	}
}