package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/vsabirov/fintech/a/services"
)

func GetBalanceHandler(ctx echo.Context) error {
	response := services.GetBalance(ctx.Param("account-id"))

	return ctx.JSON(http.StatusOK, response)
}
