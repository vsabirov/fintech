package servicectx

import (
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type ServiceContext struct {
	echo.Context

	KafkaWriter   *kafka.Writer
	ServiceLogger *logrus.Logger
}

type ServiceContextConfig struct {
	KafkaWriter   *kafka.Writer
	ServiceLogger *logrus.Logger
}

func ServiceContextExtender(config ServiceContextConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sc := &ServiceContext{
				c,

				config.KafkaWriter,
				config.ServiceLogger,
			}

			return next(sc)
		}
	}
}
