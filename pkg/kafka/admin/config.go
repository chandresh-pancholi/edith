package admin

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AdminClient struct {
	adminClient *kafka.AdminClient
}

func NewAdminClient(logger *zap.Logger) *AdminClient {
	ac := new(AdminClient)

	logger.Info("Creating kafka admin client")

	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": viper.GetString("kafka.brokers")})
	if err != nil {
		logger.Error("Failed to create Kafka Admin Client", zap.Error(err))
	}

	ac.adminClient = adminClient

	return ac
}

func (a *AdminClient) CreateTopic(topic string, partition int, replicationFactor int, logger *zap.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topicResult, err := a.adminClient.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     partition,
			ReplicationFactor: replicationFactor}},
	)

	kafka.SetAdminOperationTimeout(60 * time.Second)

	if err != nil {
		return err
	}

	for _, result := range topicResult {
		if result.Error.Code() != kafka.ErrNoError {
			logger.Error("Kafka topic creation failed",
				zap.String("topic", topic),
				zap.Int("partition", partition),
				zap.Error(err))
		} else {
			logger.Info("Kafka topic created",
				zap.String("topic", topic),
				zap.Int("partition", partition))
		}
	}

	return nil
}

func (a *AdminClient) Close() {
	a.adminClient.Close()
}
