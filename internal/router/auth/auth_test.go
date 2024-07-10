package auth

import (
	"encoding/json"
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

	var body LoginRequest
	if err := json.Unmarshal([]byte(reqBody), &body); err != nil {
		t.Fatal(err)
	}
	ctx.Set("validatedBody", &body)

	mockUsecase := router.usecase.(*MockHandlerUsecase)
	tokenString := "mock_jwt_token"
	mockUsecase.On("AuthenticateUser", "user_login", "pAssw_ord123").Return(tokenString, nil)

	err := router.handleLogin(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusOK, rec.Code)

	expectedResponse := `{"message":"Пользователь успешно аутентифицирован","token":"mock_jwt_token"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	mockUsecase.AssertExpectations(t)
}
func TestHandleLogin_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	mockUsecase := new(MockHandlerUsecase)

	// Утверждаем, что при вызове AuthenticateUser с определёнными параметрами будет возвращаться ошибка
	mockUsecase.On("AuthenticateUser", "johndoe", "").Return("", errors.New("invalid login or password"))

	router := NewHttpRouter(e, mockUsecase, validator.New())
	reqBody := `{"login": "johndoe", "password":""}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	var body LoginRequest
	if err := json.Unmarshal([]byte(reqBody), &body); err != nil {
		t.Fatal(err)
	}
	ctx.Set("validatedBody", &body)

	err := router.handleLogin(ctx)
	if assert.Error(err) {
		he, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatalf("expected *echo.HTTPError, got %T", err)
		}
		assert.Equal(http.StatusUnauthorized, he.Code)
	}
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
	var body RegisterRequest
	if err := json.Unmarshal([]byte(reqBody), &body); err != nil {
		t.Fatal(err)
	}
	ctx.Set("validatedBody", &body)

	mockUsecase := router.usecase.(*MockHandlerUsecase)
	mockUsecase.On("RegisterUser", mock.Anything).Return(nil)

	err := router.handleRegister(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusCreated, rec.Code)

	expectedResponse := `"Пользователь успешно зарегистрирован"`
	assert.Equal(strings.TrimSpace(expectedResponse), strings.TrimSpace(rec.Body.String()))

	mockUsecase.AssertExpectations(t)
}
func TestHandleRegister_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	mockUsecase := new(MockHandlerUsecase)

	user := &dto.USERS{
		EMAIL:    "john.doe@example.com",
		LOGIN:    "johndoe",
		SURNAME:  "Doe",
		PASSWORD: "securePwd123",
	}
	mockUsecase.On("RegisterUser", user).Return(errors.New("Validation error"))

	router := NewHttpRouter(e, mockUsecase, validator.New())
	reqBody := `{"username": "", "surname": "Doe", "email": "john.doe@example.com", "login": "johndoe", "password": "securePwd123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	var body RegisterRequest
	if err := json.Unmarshal([]byte(reqBody), &body); err != nil {
		t.Fatal(err)
	}
	ctx.Set("validatedBody", &body)

	err := router.handleRegister(ctx)
	if assert.Error(err) {
		he, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatalf("expected *echo.HTTPError, got %T", err)
		}
		assert.Equal(http.StatusBadRequest, he.Code)
	}
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

	// Set validatedBody in context manually (simulating middleware behavior)
	var body UpdateRequest
	if err := json.Unmarshal([]byte(reqBody), &body); err != nil {
		t.Fatal(err)
	}
	ctx.Set("validatedBody", &body)

	mockUsecase.On("UpdateUserByLogin", "user_login", "Bearer mock_jwt_token", mock.Anything).Return(nil)

	err := router.handleUpdateUserByLogin(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusOK, rec.Code)

	expectedResponse := `"Пользователь успешно обновлен"`
	assert.JSONEq(expectedResponse, rec.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestHandleUpdateUserByLogin_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"username": "Updated", "surname": "User", "email": "updated@example.com", "password": "updatedPassword"}`
	req := httptest.NewRequest(http.MethodPut, "/update/invalid_login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("login")
	ctx.SetParamValues("invalid_login")

	var body UpdateRequest
	err := json.Unmarshal([]byte(reqBody), &body)
	if err != nil {
		t.Fatal(err)
	}
	ctx.Set("validatedBody", &body)

	err = router.handleUpdateUserByLogin(ctx)
	if assert.Error(err) {
		he, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatalf("expected *echo.HTTPError, got %T", err)
		}
		assert.Equal(http.StatusUnauthorized, he.Code)

	}
}
