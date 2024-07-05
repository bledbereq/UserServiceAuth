package auth

import (
	"net/http"

	models "UserServiceAuth/storage"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// IHandlerUsecase определяет интерфейс для использования в HTTP роутере
type IHandlerUsecase interface {
	RegisterUser(user *models.USERS) error
	AuthenticateUser(login, password string) (*models.USERS, error)
}

// HttpRouter представляет HTTP маршрутизатор
type HttpRouter struct {
	validator *validator.Validate
	usecase   IHandlerUsecase
}

// NewHttpRouter создает новый HttpRouter
func NewHttpRouter(e *echo.Echo, usecase IHandlerUsecase, validator *validator.Validate) *HttpRouter {
	// Установка валидатора в Echo
	e.Validator = &CustomValidator{validator}

	// Инициализация роутера
	router := &HttpRouter{
		validator: validator,
		usecase:   usecase,
	}

	// Настройка маршрутов
	e.POST("/login", router.handleLogin)
	e.POST("/register", router.handleRegister)

	return router
}

// CustomValidator представляет специальный валидатор для Echo
type CustomValidator struct {
	validator *validator.Validate
}

// Validate валидирует структуру
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error":   "Ошибка валидации",
			"details": err.Error(),
		})
	}
	return nil
}

// handleLogin обрабатывает запросы на /login
func (h *HttpRouter) handleLogin(ctx echo.Context) error {
	type LoginRequest struct {
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	req := new(LoginRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if err := h.validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation error",
			"details": err.Error(),
		})
	}

	user, err := h.usecase.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	_ = user

	return ctx.JSON(http.StatusOK, map[string]string{"message": "User authenticated successfully", "user": "jwt"})
}

// handleRegister обрабатывает запросы на /register
func (h *HttpRouter) handleRegister(ctx echo.Context) error {
	type RegisterRequest struct {
		Username string `json:"username" validate:"required"`
		Surname  string `json:"surname" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	req := new(RegisterRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if err := h.validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation error",
			"details": err.Error(),
		})
	}

	user := &models.USERS{
		USERNAME: req.Username,
		SURNAME:  req.Surname,
		EMAIL:    req.Email,
		LOGIN:    req.Login,
		PASSWORD: req.Password,
	}

	if err := h.usecase.RegisterUser(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, "User registered successfully")
}
