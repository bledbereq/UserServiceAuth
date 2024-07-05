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

func (m *MockHandlerUsecase) UpdateUserByID(id uint, user *storage.USERS) error {
	args := m.Called(id, user)
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

	mockUsecase := router.usecase.(*MockHandlerUsecase)
	mockUsecase.On("RegisterUser", mock.Anything).Return(errors.New("DuplicateEmail"))

	err := router.handleRegister(ctx)
	assert.NoError(err)

	assert.Equal(http.StatusBadRequest, rec.Code)

	expectedResponse := `{"error":"DuplicateEmail"}`
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

func TestHandleUpdateUserByID_ValidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	mockUsecase := new(MockHandlerUsecase)
	router := NewHttpRouter(e, mockUsecase, validator.New())

	reqBody := `{"login": "newlogin", "username": "John", "surname": "Doe", "email": "john.doe@example.com", "password": "newPwd123"}`
	req := httptest.NewRequest(http.MethodPut, "/update/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	mockUsecase.On("UpdateUserByID", uint(1), mock.Anything).Return(nil)

	err := router.handleUpdateUserByID(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusOK, rec.Code)

	expectedResponse := `"User updated successfully"`
	assert.JSONEq(expectedResponse, rec.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestHandleUpdateUserByID_UserNotFound(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	mockUsecase := new(MockHandlerUsecase)
	router := NewHttpRouter(e, mockUsecase, validator.New())

	reqBody := `{"login": "newlogin", "username": "John", "surname": "Doe", "email": "john.doe@example.com", "password": "newPwd123"}`
	req := httptest.NewRequest(http.MethodPut, "/update/999", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("999")

	mockUsecase.On("UpdateUserByID", uint(999), mock.Anything).Return(errors.New("user with this id not exists"))

	err := router.handleUpdateUserByID(ctx)

	assert.NoError(err)
	assert.Equal(http.StatusNotFound, rec.Code)

	expectedResponse := `{"error":"user with this id not exists"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	mockUsecase.AssertExpectations(t)
}
func TestHandleUpdateUserByID_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	router := NewHttpRouter(e, &MockHandlerUsecase{}, validator.New())

	reqBody := `{"login": "updated_login", "username": "Updated", "surname": "User", "email": "updated@example.com", "password": "updatedPassword"}`
	req := httptest.NewRequest(http.MethodPut, "/update/invalid_id", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := router.handleUpdateUserByID(ctx)

	assert.Equal(http.StatusBadRequest, rec.Code)

	expectedResponse := `{"error":"Invalid user ID"}`
	assert.JSONEq(expectedResponse, rec.Body.String())

	assert.NoError(err)
}
