package services

import (
	"context"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/vsabirov/fintech/b/servicectx"
	"github.com/vsabirov/fintech/b/services/dal"
)

type TransferRequest struct {
	ID     string  `json:"id" xml:"id"`
	Amount float64 `json:"amount" xml:"amount"`

	Receiver string `json:"receiver" xml:"receiver"`
	Sender   string `json:"sender" xml:"sender"`
}

func Transfer(request TransferRequest, sctx servicectx.ServiceContext) error {
	err := sctx.Database.Transfer(context.Background(), dal.TransferParams{
		ID:     request.ID,
		Amount: request.Amount,

		Receiver: request.Receiver,
		Sender:   request.Sender,
	})

	go func() {
		time.Sleep(time.Second * 30)

		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)

		if random.Intn(1) == 0 {
			sctx.Logger.WithFields(logrus.Fields{
				"request": request,
			}).Warn("Transfer was marked invalid, trying to restore account funds.")

			err := sctx.Database.Refund(context.Background(), dal.RefundParams{
				TransferID: request.ID,
				Amount:     request.Amount,

				Receiver: request.Receiver,
				Sender:   request.Sender,
			})

			if err != nil {
				sctx.Logger.WithFields(logrus.Fields{
					"request": request,
					"error":   err,
				}).Warn("Fund restoration failed.")

				return
			}

			sctx.Logger.WithFields(logrus.Fields{
				"request": request,
			}).Info("Funds restored successfully.")
		}
	}()

	return err
}
