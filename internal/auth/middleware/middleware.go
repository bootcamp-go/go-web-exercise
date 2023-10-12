package middleware

import (
	"app/internal/auth"
	"app/platform/web/response"
	"net/http"
)

// NewAuthenticator creates an Authenticator to handle authentication via middleware
func NewAuthenticator(au auth.AuthToken) *Authenticator {
	return &Authenticator{
		au: au,
	}
}

// Authenticator handles authentication.
type Authenticator struct {
	// au is the authenticator service.
	au auth.AuthToken
}

// NewAuthenticator creates a middleware to authenticate requests.
func (a *Authenticator) Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before
		// get token
		token := r.Header.Get("Token")
		
		// validate token
		if err := a.au.Auth(token); err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// call
		handler.ServeHTTP(w, r)

		// after
		// ...
	})
}