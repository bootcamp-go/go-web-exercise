package auth

import "errors"

var (
	// ErrAuthTokenInternal is an error that returns when an internal error occurs
	ErrAuthTokenInternal = errors.New("authenticator: internal error")

	// ErrAuthTokenInvalid is an error that returns when a token is invalid
	ErrAuthTokenInvalid = errors.New("authenticator: token invalid")

	// ErrAuthTokenNotFound is an error that returns when a token is not found
	ErrAuthTokenNotFound = errors.New("authenticator: token not found")

	// ErrAuthTokenExpired is an error that returns when a token is expired
	ErrAuthTokenExpired = errors.New("authenticator: token expired")
)

// AuthToken is an interface that contains the methods that a authenticator must implement
type AuthToken interface {
	// Auth is a method that authenticates
	Auth(token string) (err error)
}