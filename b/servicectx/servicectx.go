package servicectx

import (
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"

	"github.com/vsabirov/fintech/b/services/dal"
)

type ServiceContext struct {
	Logger      *logrus.Logger
	KafkaReader *kafka.Reader
	Database    dal.DataSource
}
