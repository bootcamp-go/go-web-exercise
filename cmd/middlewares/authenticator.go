package middlewares

import (
	"app/internal/authenticator"
	"app/pkg/web/response"
	"net/http"
)

type MiddlewareAuthenticator struct {
	// auth is an AuthenticatorToken interface to authenticate via token
	auth authenticator.AuthenticatorToken
}

// NewMiddlewareAuthenticator is a constructor
func NewMiddlewareAuthenticator(auth authenticator.AuthenticatorToken) *MiddlewareAuthenticator {
	return &MiddlewareAuthenticator{
		auth: auth,
	}
}

// Auth is a method that returns a middleware that authenticates via token
func (m *MiddlewareAuthenticator) Auth(hd http.Handler) (newHd http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before
		// -> auth
		token := r.Header.Get("Authorization")
		if err := m.auth.Auth(token); err != nil {
			code := http.StatusUnauthorized
			body := map[string]any{"message": "unauthorized"}

			response.JSON(w, code, body)
			return
		}

		// call
		hd.ServeHTTP(w, r)

		// after
		// ...
	})
}
