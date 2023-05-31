package main

import (
	"context"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/vsabirov/fintech/b/consumers"
	"github.com/vsabirov/fintech/b/servicectx"
)

const ConsumerGroupID = "service-b"

func optionalEnv(log *logrus.Logger, key string, zv string) string {
	value, valueSpecified := os.LookupEnv(key)
	if !valueSpecified {
		value = zv

		log.WithFields(logrus.Fields{
			"key":   key,
			"value": value,
		}).Info("Optional environment variable not specified, using default value.")
	}

	return value
}

func main() {
	log := logrus.New()

	kafkaPort := optionalEnv(log, "KAFKA_PORT", "9094")
	kafkaHost := optionalEnv(log, "KAFKA_HOST", "localhost")

	sctx := servicectx.ServiceContext{
		Logger: log,
		KafkaReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{kafkaHost + ":" + kafkaPort},
			GroupID:     ConsumerGroupID,
			GroupTopics: []string{consumers.TransferTopic},
		}),
	}

	consumers.ProcessMessages(context.Background(), sctx)
}
