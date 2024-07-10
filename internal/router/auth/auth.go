package auth

import (
	"net/http"
	"strings"

	dto "UserServiceAuth/storage"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type IHandlerUsecase interface {
	RegisterUser(user *dto.USERS) error
	AuthenticateUser(login, password string) (string, error)
	UpdateUserByLogin(login, token string, user *dto.USERS) error
}

type HttpRouter struct {
	validator *validator.Validate
	usecase   IHandlerUsecase
}

func NewHttpRouter(e *echo.Echo, usecase IHandlerUsecase, validator *validator.Validate) *HttpRouter {
	e.Validator = &CustomValidator{validator}

	router := &HttpRouter{
		validator: validator,
		usecase:   usecase,
	}
	e.Use(router.validateMiddleware)

	e.POST("/login", router.handleLogin)
	e.POST("/register", router.handleRegister)
	e.PUT("/update/:login", router.handleUpdateUserByLogin)

	return router
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error":   "Ошибка валидации",
			"details": err.Error(),
		})
	}
	return nil
}

func (h *HttpRouter) validateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()

		contentType := req.Header.Get(echo.HeaderContentType)
		if contentType != echo.MIMEApplicationJSON {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Content-Type, expected application/json")
		}

		var body interface{}
		path := ctx.Path()
		switch {
		case strings.HasSuffix(path, "/login"):
			body = new(LoginRequest)
		case strings.HasSuffix(path, "/register"):
			body = new(RegisterRequest)
		case strings.HasSuffix(path, "/update/:login"):
			body = new(UpdateRequest)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request path")
		}

		if err := ctx.Bind(body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		if err := h.validator.Struct(body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
				"error":   "Validation error",
				"details": err.Error(),
			})
		}

		ctx.Set("validatedBody", body)
		return next(ctx)
	}
}

func (h *HttpRouter) handleLogin(ctx echo.Context) error {
	req := ctx.Get("validatedBody").(*LoginRequest)
	token, err := h.usecase.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Пользователь успешно аутентифицирован",
		"token":   token,
	})
}

func (h *HttpRouter) handleRegister(ctx echo.Context) error {
	req := ctx.Get("validatedBody").(*RegisterRequest)

	user := &dto.USERS{
		USERNAME: req.Username,
		SURNAME:  req.Surname,
		EMAIL:    req.Email,
		LOGIN:    req.Login,
		PASSWORD: req.Password,
	}

	if err := h.usecase.RegisterUser(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, "Пользователь успешно зарегистрирован")
}

func (h *HttpRouter) handleUpdateUserByLogin(ctx echo.Context) error {
	req := ctx.Get("validatedBody").(*UpdateRequest)
	login := ctx.Param("login")
	token := ctx.Request().Header.Get("Authorization")

	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Отсутствует токен авторизации")
	}

	updatedUser := &dto.USERS{
		USERNAME: req.Username,
		SURNAME:  req.Surname,
		EMAIL:    req.Email,
		PASSWORD: req.Password,
	}

	err := h.usecase.UpdateUserByLogin(login, token, updatedUser)
	if err != nil {
		if err.Error() == "user with this login not exists" {
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, "Пользователь успешно обновлен")
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateRequest struct {
	Username string `json:"username" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
