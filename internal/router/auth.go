package auth

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type HttpRouter struct {
	//  зависимости или сервисы
	validator *validator.Validate
}

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

// Специальный валидатор для Echo, который использует go-playground/validator
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	// Валидируем структуру i с помощью внешнего валидатора
	if err := cv.validator.Struct(i); err != nil {
		// Если есть ошибки валидации, создаем HTTP ошибку с кодом 400 и сообщением об ошибке
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// Обработчик для /login
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
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Логика обработки логина с использованием юзкейсов

	return ctx.JSON(http.StatusCreated, "1jwt1")
}

// Обработчик для /register
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
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Логика обработки регистрации с использованием юзкейсов

	return ctx.JSON(http.StatusCreated, "1jwt1")
}
