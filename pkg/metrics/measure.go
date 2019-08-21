package metrics

import "go.opencensus.io/stats"

// metrics
var (
	// KafkaSuccessfulConsumed counts successfully consumed kafka events
	KafkaSuccessfulConsumed = stats.Int64("kafka/consumer/success", "Successful event consumed by kafka consumer", stats.UnitDimensionless)

	KafkaConsumedFailed = stats.Int64("kafka/consumer/failed", "event consumed by kafka consumer failed ", stats.UnitDimensionless)

	// KafkaSuccessfulProduced counts the number of pins the local peer is tracking.
	KafkaSuccessfulProduced = stats.Int64("kafka/producer/success", "Successful event produced to kafka", stats.UnitDimensionless)

	KafkaProducedFailed = stats.Int64("kafka/producer/failed", "failed to produced to kafka", stats.UnitDimensionless)
)


