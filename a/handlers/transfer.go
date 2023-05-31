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
	err := ctx.Bind(&body)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	sctx := ctx.(*servicectx.ServiceContext)

	err = services.Transfer(body, sctx.KafkaWriter)
	if err != nil {
		sctx.ServiceLogger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to call transfer service.")

		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusCreated)
}
