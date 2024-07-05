package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"UserServiceAuth/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type MockHandlerUsecase struct {
	mock.Mock
}

func (m *MockHandlerUsecase) RegisterUser(user *storage.USERS) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockHandlerUsecase) AuthenticateUser(login, password string) (*storage.USERS, error) {
	args := m.Called(login, password)
	return args.Get(0).(*storage.USERS), args.Error(1)
}

func TestHandleLogin_ValidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"login": "user_login", "password": "pAssw_ord123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	mockUsecase := router.usecase.(*MockHandlerUsecase)
	user := &storage.USERS{USERID: 1}
	mockUsecase.On("AuthenticateUser", "user_login", "pAssw_ord123").Return(user, nil)

	err := router.handleLogin(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusOK, rec.Code)

	expectedResponse := `{"message":"User authenticated successfully","user":"jwt"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestHandleLogin_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"password": "pAssw_ord123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := router.handleLogin(ctx)

	assert.Equal(http.StatusBadRequest, rec.Code)

	expectedResponse := `{"error":"Validation error","details":"Key: 'LoginRequest.Login' Error:Field validation for 'Login' failed on the 'required' tag"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	assert.NoError(err)
}

func TestHandleRegister_ValidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"username": "John", "surname": "Doe", "email": "john.doe@example.com", "login": "johndoe", "password": "securePwd123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	mockUsecase := router.usecase.(*MockHandlerUsecase)
	mockUsecase.On("RegisterUser", mock.Anything).Return(nil)

	err := router.handleRegister(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusCreated, rec.Code)

	expectedResponse := `"User registered successfully"`
	assert.Equal(strings.TrimSpace(expectedResponse), strings.TrimSpace(rec.Body.String()))

	mockUsecase.AssertExpectations(t)
}

func TestHandleRegister_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"surname": "Doe", "email": "john.doe@example.com", "login": "johndoe", "password": "securePwd123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := router.handleRegister(ctx)

	assert.Equal(http.StatusBadRequest, rec.Code)

	expectedResponse := `{"error":"Validation error","details":"Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	assert.NoError(err)
}

func TestHandleRegister_UsecaseError(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"username": "John", "surname": "Doe", "email": "john.doe@example.com", "login": "johndoe", "password": "securePwd123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	mockUsecase := router.usecase.(*MockHandlerUsecase)
	mockUsecase.On("RegisterUser", mock.Anything).Return(errors.New("Database error"))

	err := router.handleRegister(ctx)
	assert.NoError(err)

	assert.Equal(http.StatusBadRequest, rec.Code)

	expectedResponse := `{"error":"Database error"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestHandleRegister_DuplicateEmail(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"username": "John", "surname": "Doe", "email": "john.doe@example.com", "login": "johndoe", "password": "securePwd123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Устанавливаем мок для возврата ошибки DuplicateEmail при регистрации пользователя
	mockUsecase := router.usecase.(*MockHandlerUsecase)
	mockUsecase.On("RegisterUser", mock.Anything).Return(errors.New("DuplicateEmail"))

	// Вызываем обработчик
	err := router.handleRegister(ctx)
	assert.NoError(err)

	assert.Equal(http.StatusBadRequest, rec.Code)

	// Проверяем ожидаемый JSON-ответ с ошибкой дубликата email
	expectedResponse := `{"error":"DuplicateEmail"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	// Проверяем вызовы к моку
	mockUsecase.AssertExpectations(t)
}

func TestHandleRegister_DuplicateLogin(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"username": "John", "surname": "Doe", "email": "newuser@example.com", "login": "johndoe", "password": "securePwd123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Устанавливаем мок для возврата ошибки DuplicateLogin при регистрации пользователя
	mockUsecase := router.usecase.(*MockHandlerUsecase)
	mockUsecase.On("RegisterUser", mock.Anything).Return(errors.New("DuplicateLogin"))

	// Вызываем обработчик
	err := router.handleRegister(ctx)
	assert.NoError(err)

	assert.Equal(http.StatusBadRequest, rec.Code)

	// Проверяем ожидаемый JSON-ответ с ошибкой дубликата логина
	expectedResponse := `{"error":"DuplicateLogin"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	// Проверяем вызовы к моку
	mockUsecase.AssertExpectations(t)
}
