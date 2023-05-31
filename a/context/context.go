package context

import (
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
)

type ServiceContext struct {
	echo.Context

	KafkaWriter *kafka.Writer
}

type ServiceContextConfig struct {
	KafkaWriter *kafka.Writer
}

func ServiceContextExtender(config ServiceContextConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sc := &ServiceContext{
				c,

				config.KafkaWriter,
			}

			return next(sc)
		}
	}
}
