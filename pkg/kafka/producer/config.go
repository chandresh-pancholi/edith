package producer

import (
	"github.com/chandresh-pancholi/edith/model"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

//EventPublishMetrics is custom metrics for MSG
type EventPublishMetrics struct {
	EvntWrittenSuccess  prometheus.Counter
	EventWrittenFailure prometheus.Counter
}

//Producer metrics
type Producer struct {
	sarama.AsyncProducer
	logger *zap.Logger
}

//NewProducer is to create new producer at server start time
func NewProducer(logger *zap.Logger) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.RequiredAcks(sarama.WaitForLocal)
	saramaConfig.Producer.Compression = sarama.CompressionCodec(sarama.CompressionSnappy)
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	asyncProducer, err := sarama.NewAsyncProducer(viper.GetStringSlice("kafka.brokers"), saramaConfig)

	p := &Producer{
		AsyncProducer: asyncProducer,
		logger:        logger,
	}
	p.Errors()
	p.Successes()

	return p, err
}

//Publish is to emit event to Kafka
func (p *Producer) Publish(event model.Event) {
	p.AsyncProducer.Input() <- &sarama.ProducerMessage{
		Topic:    event.Topic,
		Key:      sarama.StringEncoder(event.Key),
		Value:    sarama.StringEncoder(event.Payload),
		Headers:  event.Headers,
		Metadata: event.Metadata,
	}
}

//Errors keep  the track of failed messages.
func (p *Producer) Errors() {
	go func() {
		for err := range p.AsyncProducer.Errors() {
			keyBytes, errEncode := err.Msg.Key.Encode()
			if errEncode != nil {
				p.logger.Error("key encoding failed for failed message", zap.String("topic", err.Msg.Topic))
			}

			//TODO custom prometheus counter metrics for failure
			//TODO pass this message to Elasticsearch as undelivered message so that daemon can process it further

			//metrics.Counter("kafka_event_ingestion_failed", map[string]string{"status": "failed", "topic": err.Msg.Topic}).Inc()
			p.logger.Error("failed to emit event to kafka", zap.String("topic", err.Msg.Topic), zap.String("key", string(keyBytes)), zap.String("error", err.Error()))
		}
	}()
}

//Successes is to check if message successfully delivered to kafka
func (p *Producer) Successes() {
	go func() {
		for msg := range p.AsyncProducer.Successes() {
			keyBytes, errEncode := msg.Key.Encode()
			if errEncode != nil {
				p.logger.Error("key encoding failed for failed message", zap.String("topic", msg.Topic))
			}

			//https://github.com/census-instrumentation/opencensus-go/blob/master/examples/quickstart/stats.go
			//metrics.Counter("kafka_event_ingestion_success", map[string]string{"status": "success", "topic": edith.Topic}).Inc()
			//TODO custom prometheus counter metrics for success
			p.logger.Info("successfully ingested event to kafka", zap.String("topic", msg.Topic), zap.String("key", string(keyBytes)))
		}
	}()
}
