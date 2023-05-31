package services

import (
	"context"

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
	return sctx.Database.Transfer(context.Background(), dal.TransferParams{
		ID:     request.ID,
		Amount: request.Amount,

		Receiver: request.Receiver,
		Sender:   request.Sender,
	})
}
