package uscase

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewTokenService(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *TokenService {
	return &TokenService{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (s *TokenService) GenerateToken(username, email, login string, isadmin bool) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"email":    email,
		"login":    login,
		"isadmin":  isadmin,
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *TokenService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Возвращаем публичный ключ для проверки подписи
		return s.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
