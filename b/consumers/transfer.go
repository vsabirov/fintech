package consumers

import (
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/vsabirov/fintech/b/servicectx"
	"github.com/vsabirov/fintech/b/services"
)

func Transfer(message kafka.Message, sctx servicectx.ServiceContext) error {
	var request services.TransferRequest
	err := json.Unmarshal(message.Value, &request)
	if err != nil {
		return err
	}

	return services.Transfer(request, sctx)
}
