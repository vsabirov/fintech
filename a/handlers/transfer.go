package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func TransferHandler(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
