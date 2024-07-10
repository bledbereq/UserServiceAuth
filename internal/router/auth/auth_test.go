package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	dto "UserServiceAuth/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type MockHandlerUsecase struct {
	mock.Mock
}

func (m *MockHandlerUsecase) RegisterUser(user *dto.USERS) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockHandlerUsecase) AuthenticateUser(login, password string) (string, error) {
	args := m.Called(login, password)
	return args.String(0), args.Error(1)
}

func (m *MockHandlerUsecase) UpdateUserByLogin(login, token string, user *dto.USERS) error {
	args := m.Called(login, token, user)
	return args.Error(0)
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
	// Моковый токен
	tokenString := "mock_jwt_token"
	mockUsecase.On("AuthenticateUser", "user_login", "pAssw_ord123").Return(tokenString, nil)

	err := router.handleLogin(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusOK, rec.Code)

	expectedResponse := `{"message":"User authenticated successfully","token":"mock_jwt_token"}`
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

func TestHandleRegister_DuplicateLogin(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"username": "John", "surname": "Doe", "email": "newuser@example.com", "login": "johndoe", "password": "securePwd123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	mockUsecase := router.usecase.(*MockHandlerUsecase)
	mockUsecase.On("RegisterUser", mock.Anything).Return(errors.New("DuplicateLogin"))

	err := router.handleRegister(ctx)
	assert.NoError(err)

	assert.Equal(http.StatusBadRequest, rec.Code)

	expectedResponse := `{"error":"DuplicateLogin"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestHandleUpdateUserByLogin_ValidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	mockUsecase := new(MockHandlerUsecase)
	router := NewHttpRouter(e, mockUsecase, validator.New())

	reqBody := `{"login": "newlogin", "username": "John", "surname": "Doe", "email": "john.doe@example.com", "password": "newPwd123"}`
	req := httptest.NewRequest(http.MethodPut, "/update/user_login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer mock_jwt_token")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("login")
	ctx.SetParamValues("user_login")

	mockUsecase.On("UpdateUserByLogin", "user_login", "Bearer mock_jwt_token", mock.Anything).Return(nil)

	err := router.handleUpdateUserByLogin(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusOK, rec.Code)

	expectedResponse := `"User updated successfully"`
	assert.JSONEq(expectedResponse, rec.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestHandleUpdateUserByLogin_UserNotFound(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	mockUsecase := new(MockHandlerUsecase)
	router := NewHttpRouter(e, mockUsecase, validator.New())

	reqBody := `{"login": "newlogin", "username": "John", "surname": "Doe", "email": "john.doe@example.com", "password": "newPwd123"}`
	req := httptest.NewRequest(http.MethodPut, "/update/nonexistent_user", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer mock_jwt_token")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("login")
	ctx.SetParamValues("nonexistent_user")

	mockUsecase.On("UpdateUserByLogin", "nonexistent_user", "Bearer mock_jwt_token", mock.Anything).Return(errors.New("user with this login not exists"))

	err := router.handleUpdateUserByLogin(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusNotFound, rec.Code)

	expectedResponse := `{"error":"user with this login not exists"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestHandleUpdateUserByLogin_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"login": "updated_login", "username": "Updated", "surname": "User", "email": "updated@example.com", "password": "updatedPassword"}`
	req := httptest.NewRequest(http.MethodPut, "/update/invalid_login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("login")
	ctx.SetParamValues("invalid_login")

	err := router.handleUpdateUserByLogin(ctx)

	assert.Equal(http.StatusUnauthorized, rec.Code)

	expectedResponse := `{"error":"Missing token"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	assert.NoError(err)
}
