package register

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func RegisterHandler(ctx echo.Context) error {
	// Получить данные пользователя из запроса
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

	// Отправить ответ с данными пользователя
	return ctx.JSON(http.StatusCreated, "1jwt1")
}
