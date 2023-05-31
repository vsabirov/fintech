package validators

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type TransferValidator struct {
	Validator *validator.Validate
}

func (tv *TransferValidator) Validate(input interface{}) error {
	if err := tv.Validator.Struct(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
