package modules

import (
	"bytes"
	"fmt"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// BuildLoggerModule expects a map[string]func(string) in whhich the key
// of the map is the name of the log level and the value is the log function
// itself.
func BuildLoggerModule(loggers map[string]func(*starlark.Thread, string)) *starlarkstruct.Module {
	module := &starlarkstruct.Module{
		Name:    "log",
		Members: starlark.StringDict{},
	}

	for k, v := range loggers {
		module.Members[k] = starlark.NewBuiltin(fmt.Sprintf("log.%s", k), LogFn(v))
	}

	return module
}

func LogFn(logLevelFn func(t *starlark.Thread, msg string)) ModuleAlias {
	return func(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		buf := bytes.NewBuffer([]byte(""))
		for _, v := range args {
			buf.WriteString(v.String() + " ")
		}

		msg := strings.TrimSpace(buf.String())

		logLevelFn(t, msg)
		return starlark.None, nil
	}
}
