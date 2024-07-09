package auth

import (
	"net/http"

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

	token, err := h.usecase.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "User authenticated successfully", "token": token})
}

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

	user := &dto.USERS{
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

func (h *HttpRouter) handleUpdateUserByLogin(ctx echo.Context) error {
	type UpdateRequest struct {
		Login    string `json:"login" validate:"required"`
		Username string `json:"username" validate:"required"`
		Surname  string `json:"surname" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	req := new(UpdateRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if err := h.validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation error",
			"details": err.Error(),
		})
	}

	login := ctx.Param("login")
	token := ctx.Request().Header.Get("Authorization")

	if token == "" {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
	}

	updatedUser := &dto.USERS{
		LOGIN:    req.Login,
		USERNAME: req.Username,
		SURNAME:  req.Surname,
		EMAIL:    req.Email,
		PASSWORD: req.Password,
	}

	if err := h.usecase.UpdateUserByLogin(login, token, updatedUser); err != nil {
		if err.Error() == "user with this login not exists" {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, "User updated successfully")
}
