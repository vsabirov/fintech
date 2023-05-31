package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type getBalanceResponse struct {
	Balance float64 `json:"balance" xml:"balance"`
}

func GetBalanceHandler(ctx echo.Context) error {
	response := &getBalanceResponse{
		Balance: 1337.5,
	}

	return ctx.JSON(http.StatusOK, response)
}
