package token

type Manager interface {
	ParseAccessToken(token string) (*AccessClaims, error)
}
