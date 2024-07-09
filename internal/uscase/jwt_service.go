package uscase

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
	privateKey *rsa.PrivateKey
}

func NewTokenService(privateKey *rsa.PrivateKey) *TokenService {
	return &TokenService{privateKey: privateKey}
}

func (s *TokenService) GenerateToken(userID uint, username, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
