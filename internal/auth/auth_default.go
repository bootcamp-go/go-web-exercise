package auth

// NewAuthTokenBasic returns a new AuthBasic
func NewAuthTokenBasic(token string) *AuthBasic {
	return &AuthBasic{
		Token: token,
	}
}

// AuthBasic is a struct that contains the basic data of a authenticator
type AuthBasic struct {
	// Token is a string that contains the token
	Token string
}

// Auth is a method that authenticates
func (a *AuthBasic) Auth(token string) (err error) {
	if a.Token != token {
		return ErrAuthTokenInvalid
	}
	return nil
}