package logs

import (
	"context"
	"fmt"

	"github.com/daidokoro/otelity/processors/starlarkprocessor/internal/modules"
	"github.com/qri-io/starlib/re"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	jsonlib "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"go.uber.org/zap"
)

func NewProcessor(ctx context.Context, logger *zap.Logger,
	code, entry string, consumer consumer.Logs) *Processor {
	return &Processor{
		logger: logger,
		code:   code,
		entry:  entry,
		queue:  make(chan plog.Logs),
		next:   consumer,
		thread: &starlark.Thread{
			Name: "log/processor",
			Print: func(thread *starlark.Thread, msg string) {
				logger.Debug(msg, zap.String("thread", thread.Name), zap.String("source", "starlark/code"))
			},
		},
	}
}

type Processor struct {
	plog.JSONMarshaler
	plog.JSONUnmarshaler
	logger         *zap.Logger
	code           string
	entry          string
	queue          chan plog.Logs
	thread         *starlark.Thread
	transformFn    starlark.Value
	convertEventFn starlark.Value
	next           consumer.Logs
}

func (p *Processor) Start(ctx context.Context, _ component.Host) error {
	modules, err := p.loadModules()
	if err != nil {
		return fmt.Errorf("failed to load starlark modules; %q", err)
	}

	globals, err := starlark.ExecFileOptions(&syntax.FileOptions{}, p.thread, "", p.code, modules)
	if err != nil {
		return err
	}

	// Retrieve a module global.
	var ok bool
	if p.transformFn, ok = globals[p.entry]; !ok {
		return fmt.Errorf("starlark: no '%s' function defined in script for entrypoint", p.entry)
	}

	if p.convertEventFn, ok = jsonlib.Module.Members["decode"]; !ok {
		return fmt.Errorf("starlark: no 'json.decode' function defined in env")
	}

	go func() {
		for ld := range p.queue {
			p.logger.Debug("emitting telemetry data", zap.Int("#logs", ld.LogRecordCount()))
			if err := p.next.ConsumeLogs(ctx, ld); err != nil {
				p.logger.Error("failed to emit telemetry data", zap.Error(err))
			}
		}
	}()

	return nil
}

func (p *Processor) Shutdown(context.Context) error { return nil }

func (p *Processor) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	b, err := p.MarshalLogs(ld)
	if err != nil {
		return err
	}

	// convert to starlark value
	event, err := starlark.Call(p.thread, p.convertEventFn, starlark.Tuple{starlark.String(string(b))}, nil)
	if err != nil {
		return fmt.Errorf("error converting telemetry event to starlark: %w", err)
	}

	// Call the function.
	result, err := starlark.Call(p.thread, p.transformFn, starlark.Tuple{event}, nil)
	if err != nil {
		return fmt.Errorf("error calling entrypoint function: %w", err)
	}

	if result.String() == "None" {
		p.logger.Error("entrypoint function returned an empty value, passing record with no changes", zap.String("result", result.String()))
		return p.next.ConsumeLogs(ctx, ld)
	}

	if ld, err = p.UnmarshalLogs([]byte(result.String())); err != nil {
		return fmt.Errorf("error unmarshalling logs data from starlark: %w", err)
	}

	// if there are no logs, return
	if ld.LogRecordCount() == 0 {
		return nil
	}

	return p.next.ConsumeLogs(ctx, ld)
}

func (p *Processor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (p *Processor) loadModules() (starlark.StringDict, error) {
	// define loggers for starklark logger moduels
	loggers := modules.BuildLoggerModule(map[string]func(*starlark.Thread, string){
		"info": func(t *starlark.Thread, msg string) {
			p.logger.Info(msg, zap.String("thread", t.Name), zap.String("source", "starlark/code"))
		},
		"warn": func(t *starlark.Thread, msg string) {
			p.logger.Warn(msg, zap.String("thread", t.Name), zap.String("source", "starlark/code"))
		},
		"error": func(t *starlark.Thread, msg string) {
			p.logger.Error(msg, zap.String("thread", t.Name), zap.String("source", "starlark/code"))
		},
	})

	modules := starlark.StringDict{
		"json": jsonlib.Module,
		"log":  loggers,
		"emit": starlark.NewBuiltin("emit", modules.EmitFn(p.logger, p.queue, func(b []byte) (plog.Logs, error) {
			return p.UnmarshalLogs(b)
		})),
	}

	// add regex module
	regexMod, err := re.LoadModule()
	if err != nil {
		return nil, err
	}

	for k, v := range regexMod {
		modules[k] = v
	}

	return modules, nil
}
