package functions

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Проверка длины и сложности пароля с помощью регулярного выражения

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
