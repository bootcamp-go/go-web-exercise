package auth

// NewAuthTokenMock creates a new auth token mock
func NewAuthTokenMock() *AuthTokenMock {
	return &AuthTokenMock{}
}

// AuthTokenMock is a struct that mocks the auth token
type AuthTokenMock struct {
	// FuncAuth is a function to proxy the auth method
	FuncAuth func(token string) (err error)

	// Spy
	Spy struct {
		// AuthCalls is a counter that counts the number of calls to the auth method
		AuthCalls int
	}
}

// Auth is a method that mocks the auth method
func (a *AuthTokenMock) Auth(token string) (err error) {
	// increment the number of calls
	a.Spy.AuthCalls++

	// proxy the method
	err = a.FuncAuth(token)
	return
}