package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"

	"github.com/vsabirov/fintech/a/handlers"
	"github.com/vsabirov/fintech/a/servicectx"
)

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

	kafkaWriter := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaHost + ":" + kafkaPort),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	server := echo.New()

	server.Use(servicectx.ServiceContextExtender(servicectx.ServiceContextConfig{
		KafkaWriter:   kafkaWriter,
		ServiceLogger: log,
	}))

	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,

		LogValuesFunc: func(ctx echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))

	server.Use(middleware.Recover())
	server.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))

	server.POST("/accounts/:account-id/transfer", handlers.TransferHandler)

	port := optionalEnv(log, "SERVICE_A_PORT", "2500")
	server.Logger.Fatal(server.Start(":" + port))
}
