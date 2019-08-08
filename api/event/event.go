package event

import (
	"encoding/json"
	"github.com/chandresh-pancholi/edith/model"
	"github.com/chandresh-pancholi/edith/pkg/kafka/producer"
	"net/http"

	"go.uber.org/zap"
)

//Handler is a struct for event emmitter
type Handler struct {
	producer *producer.Producer
	logger   *zap.Logger
}

//NewEventHandler to register event handler
func NewEventHandler(producer *producer.Producer, logger *zap.Logger) *Handler {
	eh := new(Handler)

	eh.producer = producer
	eh.logger = logger

	return eh

}

//Emit is to emit event to Kafka
func (eh *Handler) Emit(w http.ResponseWriter, r *http.Request) {
	var event model.Event

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&event)
	if err != nil {
		eh.logger.Error("Json decoding failed here", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
	}

	eh.logger.Info("ingesting event to kafka", zap.Any("event", event))
	eh.producer.Publish(event)
	if err != nil {
		eh.logger.Error("Kafka publisher failed",
			zap.String("topic", event.Topic),
			zap.String("service_name", event.ServiceName),
			zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
