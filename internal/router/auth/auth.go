package auth

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// HttpRouter представляет маршрутизатор HTTP с зависимостями и сервисами
type HttpRouter struct {
	validator *validator.Validate
}

// NewHttpRouter создает новый HttpRouter и настраивает маршруты
func NewHttpRouter(e *echo.Echo) *HttpRouter {
	// Создание экземпляра валидатора
	validator := validator.New()

	// Установка валидатора в Echo
	e.Validator = &CustomValidator{validator}

	// Инициализация роутера
	router := &HttpRouter{
		validator: validator,
	}

	// Настройка маршрутов
	e.POST("/login", router.handleLogin)
	e.POST("/register", router.handleRegister)

	return router
}

// CustomValidator представляет специальный валидатор для Echo, который использует go-playground/validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate валидирует структуру i с помощью внешнего валидатора и возвращает HTTP-ошибку, если валидация не прошла
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

	req := &LoginRequest{}
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат запроса"})
	}

	// Проверка валидности полей запроса
	if err := h.validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Ошибка валидации",
			"details": err.Error(),
		})
	}

	// Логика обработки логина с использованием юзкейсов
	// ...

	return ctx.JSON(http.StatusCreated, "1jwt1")
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

	req := &RegisterRequest{}
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат запроса"})
	}

	// Проверка валидности полей запроса
	if err := h.validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Ошибка валидации",
			"details": err.Error(),
		})
	}

	// Логика обработки регистрации с использованием юзкейсов
	// ...

	return ctx.JSON(http.StatusCreated, "1jwt1")
}
