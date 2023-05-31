package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/vsabirov/fintech/a/servicectx"
	"github.com/vsabirov/fintech/a/services"
)

func TransferHandler(ctx echo.Context) error {
	var body services.TransferRequest
	if err := ctx.Bind(&body); err != nil {
		ctx.Logger().Error(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	body.Sender = ctx.Param("account-id")

	sctx := ctx.(*servicectx.ServiceContext)

	if err := services.Transfer(body, sctx.KafkaWriter); err != nil {
		sctx.ServiceLogger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to call transfer service.")

		ctx.Logger().Error(err)

		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusCreated)
}
