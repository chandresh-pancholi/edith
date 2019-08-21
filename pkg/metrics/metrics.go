package metrics

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var EMPTY  = ""
func counterMetrics(viewName, name string, labels map[string]string) {

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

func counter(viewName, name string)  {
	//https://github.com/ipfs/ipfs-cluster/blob/d8c20adc4e4779cad10bc8fd3eddfdb694868a9f/observations/metrics.go
	mLinesIn := stats.Int64(name, EMPTY, stats.UnitDimensionless)

	countView := &view.View{
		Name: viewName,
		Measure: mLinesIn,
		Description: EMPTY,
		Aggregation: view.Count(),
	}
	if err := view.Register(countView); err != nil {
		//p.logger.Error("Failed to register views:", zap.Error(err))
	}
}
