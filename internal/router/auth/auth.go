package auth

import (
	"net/http"
	"strconv"

	dto "UserServiceAuth/storage"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type IHandlerUsecase interface {
	RegisterUser(user *dto.USERS) error
	AuthenticateUser(login, password string) (*dto.USERS, error)
	UpdateUserByID(id uint, user *dto.USERS) error
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
	e.PUT("/update/:id", router.handleUpdateUserByID)

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
		switch req.URL.Path {
		case "/login":
			body = new(dto.LoginRequest)
		case "/register":
			body = new(dto.RegisterRequest)
		case "/update/:id":
			body = new(dto.UpdateRequest)
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
	req := ctx.Get("validatedBody").(*dto.LoginRequest)

	user, err := h.usecase.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	_ = user
	token := "jwt_token"

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Пользователь успешно аутентифицирован",
		"user":    token,
	})
}

func (h *HttpRouter) handleRegister(ctx echo.Context) error {
	req := ctx.Get("validatedBody").(*dto.RegisterRequest)

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

func (h *HttpRouter) handleUpdateUserByID(ctx echo.Context) error {
	req := ctx.Get("validatedBody").(*dto.UpdateRequest)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный ID пользователя")
	}

	updatedUser := &dto.USERS{
		USERNAME: req.Username,
		SURNAME:  req.Surname,
		EMAIL:    req.Email,
		PASSWORD: req.Password,
	}

	if err := h.usecase.UpdateUserByID(uint(id), updatedUser); err != nil {
		if err.Error() == "user with this id not exists" {
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, "Пользователь успешно обновлен")
}
