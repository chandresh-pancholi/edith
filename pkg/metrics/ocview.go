package metrics

import (
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
)

type OCView struct {
	log *zap.Logger
}


// OCViews, which is just the aggregation of the metrics
func (oc *OCView) OCViews() []*view.View {
	kafkaConsumerSuccessView := &view.View{
		Name: "edith/kafka/consumer/success",
		Measure:     KafkaSuccessfulConsumed,
		TagKeys: []tag.Key{oc.Key("topic"), oc.Key("groupId")},
		Aggregation: view.Count(),
	}

	kafkaConsumerFailureView := &view.View{
		Name: "edith/kafka/consumer/failure",
		Measure:     KafkaConsumedFailed,
		TagKeys: []tag.Key{oc.Key("topic"), oc.Key("groupId")},
		Aggregation: view.Count(),
	}

	kafkaProducedSuccess := &view.View{
		Name: "edith/kafka/produced/success",
		Measure: KafkaSuccessfulProduced,
		TagKeys: []tag.Key{oc.Key("topic")},
		Aggregation: view.Count(),
	}

	kafkaProducedFailed := &view.View{
		Name: "edith/kafka/produced/success",
		Measure: KafkaProducedFailed,
		TagKeys: []tag.Key{oc.Key("topic")},
		Aggregation: view.Count(),
	}

	return []*view.View{
		kafkaConsumerSuccessView,
		kafkaConsumerFailureView,
		kafkaProducedSuccess,
		kafkaProducedFailed,
	}
}

func (oc *OCView) Key(name string) tag.Key {
	key, err := tag.NewKey(name)
	if err != nil {
		oc.log.Error("oc tag new Key failed", zap.String("Key", name), zap.Error(err))
	}
	return key
}
