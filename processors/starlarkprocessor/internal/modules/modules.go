package modules

import (
	"fmt"

	jsonlib "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	LogsModuleName string = "log"
	JSONModuleName string = "json"
	EmitFnName     string = "emit"
)

type OTLPModuleTelemetryType interface {
	Build(dict *starlark.Dict) starlark.Value
}

// BuildOTLPModule creates the OTLP module
func BuildOTLPModule[T OTLPModuleTelemetryType](otlptype T, telemetryType string) *starlarkstruct.Module {
	return &starlarkstruct.Module{
		Name: "otlp",
		Members: starlark.StringDict{
			telemetryType: starlark.NewBuiltin("otlp.log", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("otlp.log: expected exactly 1 argument, got %d", len(args))
				}

				dict, ok := args[0].(*starlark.Dict)
				if !ok {
					return nil, fmt.Errorf("otlp.log: expected dict argument, got %s", args[0].Type())
				}

				return otlptype.Build(dict), nil
			}),
		},
	}
}

func OTLPTypeFromBytes[T OTLPModuleTelemetryType](otlptype T, thread *starlark.Thread, b []byte) (starlark.Value, error) {
	decoder := jsonlib.Module.Members["decode"]
	val, err := starlark.Call(thread, decoder, starlark.Tuple{starlark.String(b)}, nil)
	if err != nil {
		return nil, err
	}

	// Convert starlark value to dict
	dict, ok := val.(*starlark.Dict)
	if !ok {
		return nil, fmt.Errorf("expected dict, got %s", val.Type())
	}

	return otlptype.Build(dict), nil
}
