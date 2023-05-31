package consumers

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/vsabirov/fintech/b/servicectx"
)

const (
	TransferTopic string = "transfer"
)

func ProcessMessages(ctx context.Context, sctx servicectx.ServiceContext) {
	sctx.Logger.Info("Listening for messages.")

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		message, err := sctx.KafkaReader.ReadMessage(ctx)
		if err != nil {
			sctx.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to fetch message.")

			continue
		}

		sctx.Logger.WithFields(logrus.Fields{
			"message": message,
		}).Info("Processing new message.")

		switch message.Topic {
		case TransferTopic:
			go func() {
				err = Transfer(message)

				if err != nil {
					sctx.Logger.WithFields(logrus.Fields{
						"message": message,
						"error":   err,
					}).Error("Failed to process transfer request.")
				}
			}()
		}
	}
}
