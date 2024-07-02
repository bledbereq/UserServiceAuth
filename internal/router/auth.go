package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Обработчик для логина
type LoginHandler struct {
	// зависимости, такие как usecases или сервисы
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		// Инициализация зависимости здесь
	}
}

func (h *LoginHandler) Handle(ctx echo.Context) error {
	type LoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	req := &LoginRequest{}
	if err := ctx.Bind(req); err != nil {
		return err
	}

	// Логика обработки логина с использованием юзкейсов

	return ctx.JSON(http.StatusCreated, "1jwt1")
}

// Обработчик для регистрации
type RegisterHandler struct {
	// зависимости, такие как usecases или сервисы
}

func NewRegisterHandler() *RegisterHandler {
	return &RegisterHandler{
		// Инициализация зависимости здесь
	}
}

func (h *RegisterHandler) Handle(ctx echo.Context) error {
	type RegisterRequest struct {
		Username string `json:"username"`
		Surname  string `json:"surname"`
		Email    string `json:"email"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	req := &RegisterRequest{}
	if err := ctx.Bind(req); err != nil {
		return err
	}

	// Логика обработки регистрации с использованием юзкейсов

	return ctx.JSON(http.StatusCreated, "1jwt1")
}

// Настройка маршрутов
func SetupRoutes(e *echo.Echo, loginHandler *LoginHandler, registerHandler *RegisterHandler) {
	e.POST("/login", loginHandler.Handle)
	e.POST("/register", registerHandler.Handle)
}
