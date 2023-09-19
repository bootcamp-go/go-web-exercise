package authenticator

import "errors"

var (
	// ErrAuthenticatorTokenInternal is an error that returns when an internal error occurs
	ErrAuthenticatorTokenInternal = errors.New("authenticator: internal error")

	// ErrAuthenticatorTokenInvalid is an error that returns when a token is invalid
	ErrAuthenticatorTokenInvalid = errors.New("authenticator: token invalid")

	// ErrAuthenticatorTokenNotFound is an error that returns when a token is not found
	ErrAuthenticatorTokenNotFound = errors.New("authenticator: token not found")

	// ErrAuthenticatorTokenExpired is an error that returns when a token is expired
	ErrAuthenticatorTokenExpired = errors.New("authenticator: token expired")
)

// Authenticator is an interface that contains the methods that a authenticator must implement
type AuthenticatorToken interface {
	// Auth is a method that authenticates
	Auth(token string) (err error)
}