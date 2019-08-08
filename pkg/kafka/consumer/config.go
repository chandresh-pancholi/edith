package consumer

import (
	"github.com/chandresh-pancholi/edith/pkg/elasticsearch"
	es "github.com/chandresh-pancholi/edith/pkg/elasticsearch/config"
	"github.com/chandresh-pancholi/edith/pkg/rest"
	"go.uber.org/zap"


	"github.com/spf13/viper"

	"github.com/Shopify/sarama"
)

// GroupHandler represents a Sarama consumer group consumer
type GroupHandler struct {
	esClient *es.ESClient
	logger   *zap.Logger
}

// NewKafkaConsumer creates a new kafka consumer
func NewKafkaConsumer(groupID string, esClient *es.ESClient, logger *zap.Logger) (sarama.ConsumerGroup, *GroupHandler, error) {
	kafkaConfig := sarama.NewConfig()

	kafkaConfig.Version = sarama.V2_1_0_0
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kafkaConfig.Consumer.Return.Errors = true
	kafkaConfig.ClientID = groupID //Client Id
	kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumerGroup, err := sarama.NewConsumerGroup(viper.GetStringSlice("kafka.brokers"), groupID, kafkaConfig)
	if err != nil {
		return nil, nil, err
	}

	handler := GroupHandler{
		esClient: esClient,
		logger:   logger,
	}

	return consumerGroup, &handler, nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *GroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.logger.Info("kafka message claimed", zap.String("value", string(message.Value)), zap.Int64("offset", message.Offset), zap.String("key", string(message.Key)))

		esResponse, err := elasticsearch.Search(c.esClient, viper.GetString("elasticsearch.client.index"), "topic", message.Topic)
		if err != nil {
			c.logger.Error("es search query failed", zap.String("index", viper.GetString("elasticsearch.client.index")), zap.String("topic", message.Topic), zap.Error(err))
		}

		consumerList, err := elasticsearch.ConsumerSearchResult(esResponse)
		if err != nil {
			c.logger.Error("es response to consumer mapping failed", zap.String("index", viper.GetString("elasticsearch.client.index")), zap.String("topic", message.Topic), zap.Any("es_response", esResponse), zap.Error(err))
		}
		rest.Invoke(message, consumerList, c.logger)

		sess.MarkMessage(message, "")
	}
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *GroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	c.logger.Info("setup kafka consumer group", zap.String("MemberID", s.MemberID()), zap.Int32("GenerationID", s.GenerationID()), zap.Any("Claims", s.Claims()))
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *GroupHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	c.logger.Info("cleanup consumer group", zap.String("MemberID", s.MemberID()), zap.Int32("GenerationID", s.GenerationID()), zap.Any("Claims", s.Claims()))
	return nil
}
