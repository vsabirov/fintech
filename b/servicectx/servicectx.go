package servicectx

import (
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type ServiceContext struct {
	Logger      *logrus.Logger
	KafkaReader *kafka.Reader
}
