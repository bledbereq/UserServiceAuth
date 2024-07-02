package login

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func LognHandler(ctx echo.Context) error {
	type LoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	req := &LoginRequest{}
	if err := ctx.Bind(req); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, "1jwt1")
}
