package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"

	"github.com/dynatrace-oss/opentelemetry-metric-go/dynatrace"
)

// Ex5-2 OneAgentAPI endpoint for metrics ingestion
// MintURL represents the OneAgentAPI endpoint for metrics ingestion
const MintURL = "http://localhost:14499/metrics/ingest"

// APIToken authenticates against the metrics ingestion API of Dynatrace
const APIToken = "################################################################################################"

// MintURL represents the DynatraceAPI endpoint for metrics ingestion
// const MintURL = "https://mou612.managed-sprint.dynalabs.io/e/<REPLACE with environmentID>/api/v2/metrics/ingest"

// initMetricsProvider configures the Dynatrace Metrics Exporter
func initMetricsProvider() error {
	opts := dynatrace.Options{}
	//opts.APIToken = APIToken
	opts.URL = MintURL

	var err error
	var exporter *dynatrace.Exporter

	if exporter, err = dynatrace.NewExporter(opts); err != nil {
		return err
	}

	processor := basic.New(simple.NewWithExactDistribution(), exporter)

	pusher := controller.New(
		processor,
		controller.WithPusher(exporter),
	)
	err = pusher.Start(context.Background())

	pusher.Start(context.Background())

	global.SetMeterProvider(pusher.MeterProvider())

	return nil
}

func reportMetric(input int, result int) {
	labelInput := label.Key("input")

	meter := global.GetMeterProvider().Meter("")
	valueRecorder := metric.Must(meter).NewInt64ValueRecorder("otel.oahttp.fib.iterations")
	//valueRecorder := metric.Must(meter).NewInt64ValueRecorder("otel.api.fib.iterations")
	
	valueRecorder.Record(context.Background(), int64(result), labelInput.String(fmt.Sprintf("n=%d", input)))
}
