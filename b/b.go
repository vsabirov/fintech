package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/segmentio/kafka-go"

	"github.com/sirupsen/logrus"

	"github.com/vsabirov/fintech/b/consumers"
	"github.com/vsabirov/fintech/b/servicectx"
	"github.com/vsabirov/fintech/b/services/dal"
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

	databaseDSN := optionalEnv(log, "DB_DSN", "user=root password=1fintechpassword1 host=localhost port=5432 dbname=fintech sslmode=disable pool_max_conns=10")
	databaseConfig, err := pgxpool.ParseConfig(databaseDSN)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to parse database connection config.")
	}

	databasePool, err := pgxpool.NewWithConfig(context.Background(), databaseConfig)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to connect to the database.")
	}

	defer databasePool.Close()

	kafkaPort := optionalEnv(log, "KAFKA_PORT", "9094")
	kafkaHost := optionalEnv(log, "KAFKA_HOST", "localhost")

	sctx := servicectx.ServiceContext{
		Logger: log,

		KafkaReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{kafkaHost + ":" + kafkaPort},
			GroupID:     ConsumerGroupID,
			GroupTopics: []string{consumers.TransferTopic},
		}),

		Database: dal.PostgresDataSource{
			Pool: databasePool,
		},
	}

	consumers.ProcessMessages(context.Background(), sctx)
}
