package jwt_helper

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"os"
)

// Config конфиг
type Config struct {
	PublicKeyFile string `yaml:"public_key_file"`
}

// Claims стандартные claims
type Claims interface {
	jwt.Claims
}

// JwtHelper jwt хелпер
type JwtHelper struct {
	publicKey *rsa.PublicKey
}

// NewHelper новый jwt хелпер
func NewHelper(cfg Config) (*JwtHelper, error) {
	buf, err := os.ReadFile(cfg.PublicKeyFile)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(buf)
	if err != nil {
		return nil, err
	}

	return &JwtHelper{
		publicKey: publicKey,
	}, nil
}

// Parse  validate, and return a token.
func (h *JwtHelper) Parse(token string, claims Claims) error {
	_, err := jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) {
		return h.publicKey, nil
	})
	if err != nil {
		return err
	}
	return claims.Valid()
}
