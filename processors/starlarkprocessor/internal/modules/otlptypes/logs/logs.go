// package logs package defines the internal starlark otlp type for log events
package logs

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var (
	stringValue = starlark.String("stringValue")
	boolValue   = starlark.String("boolValye")
	doubleValue = starlark.String("doubleValue")
	intValue    = starlark.String("intValue")
)

// Resource represents an OTLP resource
type Resource struct {
	data *starlark.Dict
}

type Attribute struct {
	data *starlark.List
}

func (a *Attribute) Freeze() {
	a.data.Freeze()
}

func (a *Attribute) Truth() starlark.Bool {
	return a.data.Len() > 0
}

func (a *Attribute) Hash() (uint32, error) {
	return a.data.Hash()
}

func (a *Attribute) String() string {
	return a.data.String()
}

func (a *Attribute) Type() string {
	return "otlp.Attribute"
}

func (a *Attribute) AttrNames() []string {
	return []string{"get", "set"}
}

func (a *Attribute) Attr(name string) (starlark.Value, error) {
	switch name {
	case "get":
		return starlark.NewBuiltin("get",
			func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("get: expected 1 argument, got %d", len(args))
				}
				key := args[0]
				iter := a.data.Iterate()
				defer iter.Done()
				var item starlark.Value
				for iter.Next(&item) {
					if dict, ok := item.(*starlark.Dict); ok {
						if kval, found, err := dict.Get(starlark.String("key")); err == nil && found {
							if kval.String() != key.String() {
								continue
							}
							val, _, _ := dict.Get(starlark.String("value"))
							return val, nil
						}
					}
				}
				return starlark.None, nil
			}), nil
	case "set":
		return starlark.NewBuiltin("set",
			func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("set: expected 2 argument, got %d", len(args))
				}

				var typ starlark.Value = stringValue
				if len(kwargs) > 0 {
					typ = kwargs[0][1]
				}

				key, value := args[0], args[1]

				// Check if key already exists and update it
				iter := a.data.Iterate()
				defer iter.Done()
				var item starlark.Value
				for iter.Next(&item) {
					if dict, ok := item.(*starlark.Dict); ok {
						if kval, found, _ := dict.Get(starlark.String("key")); found && kval.String() == key.String() {
							valueDict := starlark.NewDict(1)
							valueDict.SetKey(typ, value)
							dict.SetKey(starlark.String("value"), valueDict)
							return starlark.None, nil
						}
					}
				}

				// If key doesn't exist, append new attribute
				attrDict := starlark.NewDict(2)
				attrDict.SetKey(starlark.String("key"), key)
				valueDict := starlark.NewDict(1)
				valueDict.SetKey(typ, value)
				attrDict.SetKey(starlark.String("value"), valueDict)
				a.data.Append(attrDict)

				return starlark.None, nil
			}), nil
	}
	return starlark.None, nil
}

func (r *Resource) Attr(name string) (starlark.Value, error) {
	switch name {
	case "attributes":
		attrs, found, err := r.data.Get(starlark.String("attributes"))
		if err != nil {
			return nil, err
		}
		if !found {
			attrs = starlark.NewList(nil)
			r.data.SetKey(starlark.String("attributes"), attrs)
		}
		// return WrapList(attrs.(*starlark.List)), nil
		return &Attribute{data: attrs.(*starlark.List)}, nil
	}
	return nil, nil
}

func (r *Resource) AttrNames() []string {
	return []string{"attributes"}
}

// LogRecord represents an OTLP log record
type LogRecord struct {
	data *starlark.Dict
}

func (lr *LogRecord) Attr(name string) (starlark.Value, error) {
	switch name {
	case "observedTimeUnixNano":
		time, _, _ := lr.data.Get(starlark.String("observedTimeUnixNano"))
		return time, nil
	case "body":
		body, _, _ := lr.data.Get(starlark.String("body"))
		return body, nil
	case "attributes":
		attrs, found, err := lr.data.Get(starlark.String("attributes"))
		if err != nil {
			return nil, err
		}
		if !found {
			attrs = starlark.NewList(nil)
			lr.data.SetKey(starlark.String("attributes"), attrs)
		}
		return &Attribute{data: attrs.(*starlark.List)}, nil
	case "traceId":
		traceID, _, _ := lr.data.Get(starlark.String("traceId"))
		return traceID, nil
	case "spanId":
		spanID, _, _ := lr.data.Get(starlark.String("spanId"))
		return spanID, nil
	}
	return nil, nil
}

func (lr *LogRecord) AttrNames() []string {
	return []string{"observed_time", "body", "attributes", "trace_id", "span_id"}
}

// ScopeLogs represents OTLP scope logs
type ScopeLogs struct {
	data *starlark.Dict
}

func (sl *ScopeLogs) Attr(name string) (starlark.Value, error) {
	switch name {
	case "scope":
		scope, _, _ := sl.data.Get(starlark.String("scope"))
		return scope, nil
	case "logRecords":
		records, found, err := sl.data.Get(starlark.String("logRecords"))
		if err != nil {
			return nil, err
		}
		if !found {
			records = starlark.NewList(nil)
			sl.data.SetKey(starlark.String("logRecords"), records)
		}

		// Wrap each record in LogRecord type
		list := records.(*starlark.List)
		wrappedList := make([]starlark.Value, 0, list.Len())
		iter := list.Iterate()
		defer iter.Done()
		var item starlark.Value
		for iter.Next(&item) {
			if dict, ok := item.(*starlark.Dict); ok {
				wrappedList = append(wrappedList, &LogRecord{data: dict})
			}
		}
		return WrapList(starlark.NewList(wrappedList)), nil
	}
	return nil, nil
}

func (sl *ScopeLogs) AttrNames() []string {
	return []string{"scope", "logRecords"}
}

// ResourceLogs represents OTLP resource logs
type ResourceLogs struct {
	data *starlark.Dict
}

func (rl *ResourceLogs) Attr(name string) (starlark.Value, error) {
	switch name {
	case "resource":
		res, found, err := rl.data.Get(starlark.String("resource"))
		if err != nil {
			return nil, err
		}
		if !found {
			res = starlark.NewDict(1)
			rl.data.SetKey(starlark.String("resource"), res)
		}
		return &Resource{data: res.(*starlark.Dict)}, nil
	case "scopeLogs":
		scopeLogs, found, err := rl.data.Get(starlark.String("scopeLogs"))
		if err != nil {
			return nil, err
		}
		if !found {
			scopeLogs = starlark.NewList(nil)
			rl.data.SetKey(starlark.String("scopeLogs"), scopeLogs)
		}
		return WrapList(scopeLogs.(*starlark.List)), nil
	}
	return nil, nil
}

func (rl *ResourceLogs) AttrNames() []string {
	return []string{"resource", "scopeLogs"}
}

// OTLPLog represents an OTLP log record with helper methods
type OTLPLog struct {
	data *starlark.Dict
}

// Build - implement modules.OTLPModuleTelemetryType interface
func (l *OTLPLog) Build(dict *starlark.Dict) starlark.Value {
	l.data = dict
	return l
}

// String implements starlark.Value
func (l *OTLPLog) String() string {
	return l.data.String()
}

// Type implements starlark.Value
func (l *OTLPLog) Type() string {
	return "otlp.Log"
}

// Freeze implements starlark.Value
func (l *OTLPLog) Freeze() {
	l.data.Freeze()
}

// Truth implements starlark.Value
func (l *OTLPLog) Truth() starlark.Bool {
	return l.data.Len() > 0
}

// Hash implements starlark.Value
func (l *OTLPLog) Hash() (uint32, error) {
	return l.data.Hash()
}

// Attr implements starlark.HasAttrs
func (l *OTLPLog) Attr(name string) (starlark.Value, error) {
	switch name {
	case "resourceLogs":
		resourceLogs, found, err := l.data.Get(starlark.String("resourceLogs"))
		if err != nil {
			return nil, err
		}
		if !found {
			resourceLogs = starlark.NewList(nil)
			l.data.SetKey(starlark.String("resourceLogs"), resourceLogs)
		}
		list := resourceLogs.(*starlark.List)
		wrappedList := make([]starlark.Value, 0, list.Len())
		iter := list.Iterate()
		defer iter.Done()
		var item starlark.Value
		for iter.Next(&item) {
			if dict, ok := item.(*starlark.Dict); ok {
				wrappedList = append(wrappedList, &ResourceLogs{data: dict})
			}
		}
		return WrapList(starlark.NewList(wrappedList)), nil
	case "scopeLogs":
		scopeLogs, found, err := l.data.Get(starlark.String("scopeLogs"))
		if err != nil {
			return nil, err
		}
		if !found {
			scopeLogs = starlark.NewList(nil)
			l.data.SetKey(starlark.String("scopeLogs"), scopeLogs)
		}
		return WrapList(scopeLogs.(*starlark.List)), nil
	case "attributes":
		attrs, found, err := l.data.Get(starlark.String("attributes"))
		if err != nil {
			return nil, err
		}
		if !found {
			attrs = starlark.NewList(nil)
			l.data.SetKey(starlark.String("attributes"), attrs)
		}
		return WrapList(attrs.(*starlark.List)), nil
	}
	return nil, nil
}

// AttrNames implements starlark.HasAttrs
func (l *OTLPLog) AttrNames() []string {
	return []string{"resourceLogs", "scopeLogs", "attributes"}
}

// // ToJSON converts the log back to JSON
// func (l *OTLPLog) ToJSON() ([]byte, error) {
// 	// Convert starlark.Dict to regular map
// 	m := make(map[string]any)
// 	for _, k := range l.data.Keys() {
// 		v, _, _ := l.data.Get(k)
// 		m[k.String()] = v.String()
// 	}
// 	return json.Marshal(m)
// }

// StarlarkList wraps starlark.List to add methods
type StarlarkList struct {
	*starlark.List
}

func (l *StarlarkList) Attr(name string) (starlark.Value, error) {
	switch name {
	case "filter":
		return starlark.NewBuiltin("filter", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("filter: expected 1 argument, got %d", len(args))
			}

			lambda, ok := args[0].(starlark.Callable)
			if !ok {
				return nil, fmt.Errorf("filter: expected function, got %s", args[0].Type())
			}

			result := starlark.NewList(nil)
			iter := l.List.Iterate()
			defer iter.Done()
			var item starlark.Value
			for iter.Next(&item) {
				lambdaResult, err := starlark.Call(thread, lambda, starlark.Tuple{item}, nil)
				if err != nil {
					return nil, err
				}

				condition, ok := lambdaResult.(starlark.Bool)
				if !ok {
					return nil, fmt.Errorf("filter: function must return bool, got %s", lambdaResult.Type())
				}

				if condition {
					result.Append(item)
				}
			}
			return WrapList(result), nil
		}), nil

	case "apply":
		return starlark.NewBuiltin("apply", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("apply: expected 1 argument, got %d", len(args))
			}

			lambda, ok := args[0].(starlark.Callable)
			if !ok {
				return nil, fmt.Errorf("apply: expected function, got %s", args[0].Type())
			}

			iter := l.List.Iterate()
			defer iter.Done()
			var item starlark.Value
			for iter.Next(&item) {
				_, err := starlark.Call(thread, lambda, starlark.Tuple{item}, nil)
				if err != nil {
					return nil, err
				}
			}
			return l, nil
		}), nil

	case "range":
		return starlark.NewBuiltin("range", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("range: expected 1 argument, got %d", len(args))
			}

			lambda, ok := args[0].(starlark.Callable)
			if !ok {
				return nil, fmt.Errorf("range: expected function, got %s", args[0].Type())
			}

			result := starlark.NewList(nil)
			iter := l.List.Iterate()
			defer iter.Done()
			var item starlark.Value
			for iter.Next(&item) {
				v, err := starlark.Call(thread, lambda, starlark.Tuple{item}, nil)
				if err != nil {
					return nil, err
				}
				result.Append(v)
			}
			return WrapList(result), nil
		}), nil
	}
	return nil, nil
}

func (l *StarlarkList) AttrNames() []string {
	return []string{"range", "if", "apply", "get"}
}

// WrapList wraps a starlark.List to add methods
func WrapList(list *starlark.List) *StarlarkList {
	return &StarlarkList{List: list}
}

// BuildOTLPModule creates the OTLP module
func BuildOTLPModule() *starlarkstruct.Module {
	return &starlarkstruct.Module{
		Name: "otlp",
		Members: starlark.StringDict{
			"log": starlark.NewBuiltin("otlp.log", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("otlp.log: expected exactly 1 argument, got %d", len(args))
				}

				dict, ok := args[0].(*starlark.Dict)
				if !ok {
					return nil, fmt.Errorf("otlp.log: expected dict argument, got %s", args[0].Type())
				}

				return &OTLPLog{data: dict}, nil
			}),
		},
	}
}

// Add these methods to Resource
func (r *Resource) String() string        { return r.data.String() }
func (r *Resource) Type() string          { return "otlp.Resource" }
func (r *Resource) Freeze()               { r.data.Freeze() }
func (r *Resource) Truth() starlark.Bool  { return r.data.Len() > 0 }
func (r *Resource) Hash() (uint32, error) { return r.data.Hash() }

// Add these methods to ResourceLogs
func (rl *ResourceLogs) String() string        { return rl.data.String() }
func (rl *ResourceLogs) Type() string          { return "otlp.ResourceLogs" }
func (rl *ResourceLogs) Freeze()               { rl.data.Freeze() }
func (rl *ResourceLogs) Truth() starlark.Bool  { return rl.data.Len() > 0 }
func (rl *ResourceLogs) Hash() (uint32, error) { return rl.data.Hash() }

// Add these methods to LogRecord
func (lr *LogRecord) String() string        { return lr.data.String() }
func (lr *LogRecord) Type() string          { return "otlp.LogRecord" }
func (lr *LogRecord) Freeze()               { lr.data.Freeze() }
func (lr *LogRecord) Truth() starlark.Bool  { return lr.data.Len() > 0 }
func (lr *LogRecord) Hash() (uint32, error) { return lr.data.Hash() }

// Add these methods to ScopeLogs
func (sl *ScopeLogs) String() string        { return sl.data.String() }
func (sl *ScopeLogs) Type() string          { return "otlp.ScopeLogs" }
func (sl *ScopeLogs) Freeze()               { sl.data.Freeze() }
func (sl *ScopeLogs) Truth() starlark.Bool  { return sl.data.Len() > 0 }
func (sl *ScopeLogs) Hash() (uint32, error) { return sl.data.Hash() }

// Add Iterate method to StarlarkList to handle custom types
func (l *StarlarkList) Iterate() starlark.Iterator {
	return &listIterator{
		iter:     l.List.Iterate(),
		wrapItem: l.wrapItem,
	}
}

// listIterator wraps the standard iterator to handle custom types
type listIterator struct {
	iter     starlark.Iterator
	wrapItem func(starlark.Value) starlark.Value
}

func (it *listIterator) Done() {
	it.iter.Done()
}

func (it *listIterator) Next(p *starlark.Value) bool {
	if !it.iter.Next(p) {
		return false
	}
	if it.wrapItem != nil {
		*p = it.wrapItem(*p)
	}
	return true
}

// wrapItem wraps dict items in appropriate types
func (l *StarlarkList) wrapItem(item starlark.Value) starlark.Value {
	if dict, ok := item.(*starlark.Dict); ok {
		// Check parent type to determine what wrapper to use
		switch l.Type() {
		case "otlp.ResourceLogs":
			return &ResourceLogs{data: dict}
		case "otlp.ScopeLogs":
			return &ScopeLogs{data: dict}
		case "otlp.LogRecord":
			return &LogRecord{data: dict}
		}
	}
	return item
}

func (l *StarlarkList) Type() string {
	// Return appropriate type based on content
	if l.Len() > 0 {
		var first starlark.Value
		iter := l.List.Iterate()
		iter.Next(&first)
		iter.Done()
		if dict, ok := first.(*starlark.Dict); ok {
			if _, found, _ := dict.Get(starlark.String("resource")); found {
				return "otlp.ResourceLogs"
			}
			if _, found, _ := dict.Get(starlark.String("logRecords")); found {
				return "otlp.ScopeLogs"
			}
			if _, found, _ := dict.Get(starlark.String("body")); found {
				return "otlp.LogRecord"
			}
		}
	}
	return "list"
}
