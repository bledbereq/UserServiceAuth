package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleLogin_ValidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e)

	reqBody := `{"login": "user_login", "password": "pAssw_ord123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Вызываем обработчик
	err := router.handleLogin(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusCreated, rec.Code)

	// Проверяем, что токен вернулся корректно
	assert.Equal("1jwt1", strings.Trim(rec.Body.String(), "\"\n"))
}

func TestHandleLogin_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e)

	// Некорректный JSON запроса (отсутствует обязательное поле login)
	reqBody := `{"password": "pAssw_ord123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Вызываем обработчик
	err := router.handleLogin(ctx)

	assert.Equal(http.StatusBadRequest, rec.Code)

	// Проверяем, что в ответе есть ключ error, связанный с валидацией поля login
	expectedResponse := `{"error":"Key: 'LoginRequest.Login' Error:Field validation for 'Login' failed on the 'required' tag"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	assert.NoError(err) // Не ожидаем ошибки валидации, только статус ошибки в ответе
}

func TestHandleRegister_ValidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e)

	reqBody := `{"username": "John", "surname": "Doe", "email": "john.doe@example.com", "login": "johndoe", "password": "securePwd123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Вызываем обработчик
	err := router.handleRegister(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusCreated, rec.Code)

	// Проверяем, что токен вернулся корректно
	assert.Equal("1jwt1", strings.Trim(rec.Body.String(), "\"\n"))
}

func TestHandleRegister_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e)

	// Некорректный JSON запроса (отсутствует обязательное поле username)
	reqBody := `{"surname": "Doe", "email": "john.doe@example.com", "login": "johndoe", "password": "securePwd123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Вызываем обработчик
	err := router.handleRegister(ctx)

	assert.Equal(http.StatusBadRequest, rec.Code)

	// Проверяем, что в ответе есть ключ error, связанный с валидацией поля username
	expectedResponse := `{"error":"Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	assert.NoError(err) // Не ожидаем ошибки валидации, только статус ошибки в ответе
}
