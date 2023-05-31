package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

var TransferTopic string = "transfer"

type TransferRequest struct {
	ID     string  `json:"id" xml:"id" validate:"required"`
	Amount float64 `json:"amount" xml:"amount" validate:"required"`

	Receiver string `json:"receiver" xml:"receiver" validate:"required"`
	Sender   string `json:"sender" xml:"sender"`
}

func Transfer(request TransferRequest, kafkaWriter *kafka.Writer) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return kafkaWriter.WriteMessages(
		context.Background(),

		kafka.Message{
			Topic: TransferTopic,
			Time:  time.Now().UTC(),

			Key:   []byte(request.ID),
			Value: []byte(payload),
		},
	)
}
