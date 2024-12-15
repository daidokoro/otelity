package modules

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.uber.org/zap"
)

type ModuleAlias = func(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

func EmitFn[T any](logger *zap.Logger, queue chan T, unmarshaler func([]byte) (T, error)) ModuleAlias {
	return func(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var v starlark.Value
		if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 1, &v); err != nil {
			logger.Error("failed to unpack data", zap.Error(err))
			return nil, err
		}

		td, err := unmarshaler([]byte(v.String()))
		if err != nil {
			logger.Error("failed to unmarshal", zap.Error(err), zap.String("data", v.String()))
			return nil, fmt.Errorf("[emit] failed to unmarshal %s: %q", v.String(), err)
		}

		queue <- td
		return starlark.None, nil
	}
}
