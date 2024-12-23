package starlarktransform

import (
	"context"
	"errors"

	"github.com/daidokoro/otelity/processors/starlarkprocessor/internal/logs"
	"github.com/daidokoro/otelity/processors/starlarkprocessor/internal/metadata"
	"github.com/daidokoro/otelity/processors/starlarkprocessor/internal/metrics"
	"github.com/daidokoro/otelity/processors/starlarkprocessor/internal/traces"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

// NewFactory creates a factory for the routing processor.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		metadata.Type,
		createDefaultConfig,
		processor.WithLogs(createLogsProcessor, metadata.LogsStability),
		processor.WithMetrics(createMetricsProcessor, metadata.MetricsStability),
		processor.WithTraces(createTracesProcessor, metadata.TracesStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		// pass event back to otel without changes
		Code: "def transform(e): return json.decode(e)",
	}
}

func createLogsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (processor.Logs, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("invalid config type")
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	code, err := config.GetCode()
	if err != nil {
		return nil, err
	}

	return logs.NewProcessor(ctx, set.Logger, code, config.EntryPoint, nextConsumer), nil
}

func createMetricsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (processor.Metrics, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("invalid config type")
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	code, err := config.GetCode()
	if err != nil {
		return nil, err
	}

	return metrics.NewProcessor(ctx, set.Logger, code, config.EntryPoint, nextConsumer), nil
}

func createTracesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (processor.Traces, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("invalid config type")
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	code, err := config.GetCode()
	if err != nil {
		return nil, err
	}

	return traces.NewProcessor(ctx, set.Logger, code, config.EntryPoint, nextConsumer), nil
}
