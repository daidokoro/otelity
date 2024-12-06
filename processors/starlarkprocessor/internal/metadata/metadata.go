package metadata

import (
	"go.opentelemetry.io/collector/component"
)

const (
	typ              = "starlark"
	LogsStability    = component.StabilityLevelAlpha
	MetricsStability = component.StabilityLevelAlpha
	TracesStability  = component.StabilityLevelAlpha
)

var Type, _ = component.NewType(typ)
