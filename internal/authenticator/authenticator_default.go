package authenticator

// NewAuthenticatorTokenBasic returns a new AuthenticatorBasic
func NewAuthenticatorTokenBasic(token string) *AuthenticatorBasic {
	return &AuthenticatorBasic{
		Token: token,
	}
}

// AuthenticatorBasic is a struct that contains the basic data of a authenticator
type AuthenticatorBasic struct {
	// Token is a string that contains the token
	Token string
}

// Auth is a method that authenticates
func (a *AuthenticatorBasic) Auth(token string) (err error) {
	if a.Token != token {
		return ErrAuthenticatorTokenInvalid
	}
	return nil
}