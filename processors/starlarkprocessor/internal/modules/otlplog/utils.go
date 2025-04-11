// package otlplog package defines the internal starlark otlp type for log events
package otlplog

import (
	"fmt"

	jsonlib "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
)

func FromBytes(thread *starlark.Thread, b []byte) (starlark.Value, error) {
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

	return &OTLPLog{data: dict}, nil
}
