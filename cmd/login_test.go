package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	// Создать новый экземпляр сервера Echo
	e := echo.New()

	// Установить маршрут для регистрации пользователя
	e.POST("/login", RegisterHandler)

	// Создать тестовый запрос
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{
        "login": "user_login",
        "password": "pAssw_ord123"
    }`))
	req.Header.Set("Content-Type", "application/json")

	// Отправить тестовый запрос на серверs
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)
	responseBody := strings.Trim(w.Body.String(), "\"\n")
	assert.Equal("1jwt1", responseBody)
}
