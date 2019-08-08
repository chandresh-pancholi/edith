package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/metric"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

//CreateMetricsFactory is
func CreateMetricsFactory(namespace string) {
	metric.NewRegistry()
}

func counterMetrics(namespace, name string, labels map[string]string) {
	keyMethod, _ := tag.NewKey("status")
	mLinesIn := stats.Int64("repl/lines_in", "The number of lines read in", stats.UnitDimensionless)
	lineCountView := &view.View{
		Name:        "demo/lines_in/test",
		Measure:     mLinesIn,
		Description: "The number of lines from standard input",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyMethod},
	}
	// Register the views
	if err := view.Register(lineCountView); err != nil {
		//p.logger.Error("Failed to register views:", zap.Error(err))
	}
	ctx, _ := tag.New(context.Background(), tag.Insert(keyMethod, "success"))

	stats.Record(ctx, mLinesIn.M(1))
}

func gaugeMetrics(namespace, name string, labels map[string]string) prometheus.Gauge {
	gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        name,
			ConstLabels: labels,
		})
	return gauge
}

func histogramMetrics(namespace, name string, labels map[string]string) prometheus.Histogram {
	histogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace:   namespace,
			Name:        name,
			ConstLabels: labels,
		})
	return histogram
}

func summaryMetrics(namespace, name string, labels map[string]string) prometheus.Summary {
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:   namespace,
			Name:        name,
			ConstLabels: labels,
		})

	return summary
}
