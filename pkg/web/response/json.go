package response

import (
	"encoding/json"
	"net/http"
)

// Text writes text response
func Text(w http.ResponseWriter, code int, body string) {
	// set header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// set status code
	w.WriteHeader(code)

	// write body
	w.Write([]byte(body))
}

// JSON writes json response
func JSON(w http.ResponseWriter, code int, body any) {
	// check body
	if body == nil {
		w.WriteHeader(code)
		return
	}
	
	// marshal body
	bytes, err := json.Marshal(body)
	if err != nil {
		// default error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// set status code
	w.WriteHeader(code)

	// write body
	w.Write(bytes)
}