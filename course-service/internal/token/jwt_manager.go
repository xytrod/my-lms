package token

type JWTManager struct {
	accessSecret []byte
	issuer       string
}

var _ Manager = (*JWTManager)(nil)

func NewJWTManager(accessSecret string, issuer string) Manager {
	return &JWTManager{
		accessSecret: []byte(accessSecret),
		issuer:       issuer,
	}
}
