// package otlplog package defines the internal starlark otlp type for log events
package otlptypes

import (
	"fmt"

	jsonlib "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
)

// FromDictFn - a function that takes a starlark.Dict and returns a metric/log/trace starlark
// otlp type
type FromDictFn func(*starlark.Dict) starlark.Value

func FromBytes(thread *starlark.Thread, b []byte, fromDictFn FromDictFn) (starlark.Value, error) {
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

	return fromDictFn(dict), nil
}
