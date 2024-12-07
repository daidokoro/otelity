# starlarktransform

<!-- status autogenerated section -->
| Status        |           |
| ------------- |-----------|
| Stability     | [alpha]: traces, metrics, logs   |


[beta]: https://github.com/open-telemetry/opentelemetry-collector#beta
[sumo]: https://github.com/SumoLogic/sumologic-otel-collector
<!-- end autogenerated section -->


The starlarktransform processor modifies telemetry based on configuration using Starlark code.

Starlark is a scripting language used for configuration that is designed to be similar to Python. It is designed to be fast, deterministic, and easily embedded into other software projects.

The processor leverages Starlark to modify telemetry data while using familiar, pythonic syntax. Modifying telemetry data is as a simple as modifying a `Dict`.

## Why?

While there are a number of transform processors, most notably, the main OTTL [transform processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/transformprocessor), this processor aims to grant users more flexibility by allowing them to manipulate telemetry data using a familiar syntax.

Python is a popular, well known language, even among non-developers. By allowing Starlark code to be used as an option to transform telemetry data, users can leverage their existing knowledge of Python.


## Config

| Parameter | Desc |
| ------- | ---- |
| code | Add in-line starlark code directly to your config |
| script | Allows you to set the script source. Supports __File path__ or __HTTP URL__ |

To configure starlarktransform, you can add your code using the `code` option in the config.


```yaml
processors:
  starlarktransform:
    code: |
      def transform(event):
        event = json.decode(event)
        <your starlark code>
        return event
```


Alternatively, in cases where you prefer to not have starlark code present or visible in your Open Telemetry configuration, you can use the `script` parameter to pass your script from a file or http url.

```yaml
processors:
  starlarktransform:
	script: /path/to/script.star

# or
processors:
  starlarktransform:
	script: https://some.url.com/script.star

```



You must define a function called `transform` that accepts a single argument, `event`. This function is called by the processor and is passed the telemetry event. The function **must return** the modified, json decoded event.


## How it works

The processor uses the [Starlark-Go](https://github.com/google/starlark-go) interpreter, this allows you to run this processor without having to install a Starlark language interpreter on the host.

## Features

The starlarktransform processor gives you access to the full telemetry event payload. You are able to modify this payload using the Starklark code in any way you want. This allows you do various things such as:

- Filtering
- Adding/Removing attributes
- Modifying attributes
- Modifying telemetry data directly
- Telemetry injection based on existing values
- And more

## Libs, Functions and Functionality

While similar in syntax to Python, Starlack does not have all the functionality associated with Python. This processor does not have access to Python standard libraries and the implementation found in this processor is limited further to only the following libraries and functions:

- **json**

> The JSON library allows you to encode and decode JSON strings. The use of this library is mandatory as the telemetry data is passed to the processor as a JSON string. You must decode the JSON string to a Dict before you can modify it. **You must also return a JSON decoded Dict to the processor.**

```python
# encode dict string to json string
x = json.encode({"foo": ["bar", "baz"]})
print(x)
# output: {"foo":["bar","baz"]}
```

```python
# decode json string to dict
x = json.decode('{"foo": ["bar", "baz"]}')
```

You can read more on the JSON library [here](https://qri.io/docs/reference/starlark-packages/encoding/json)

- **print**
> You are able to use the print function to check outputs of your code. The output of the print function is sent to the Open Telemetry runtime log. Values printed by the Print function only show when running Open Telemetry in Debug mode.

```python
def transform(event):
	print("hello world")
	return json.decode(event)
```

The print statement above would result in the following output in the Open Telemetry runtime log. Again, this output is only visible when running Open Telemetry in Debug mode.
```log
2023-09-23T16:50:17.328+0200	debug	traces/processor.go:25	hello world	{"kind": "processor", "name": "starlarktransform/traces", "pipeline": "traces", "thread": "trace.processor", "source": "starlark/code"}
```


- **re** (regex)
> Support for Regular Expressions coming soon


Note that you can define your own functions within your Starlark code, however, there must be at least one function named `transform` that accepts a single argument `event` and returns a JSON decoded Dict, this function can call all your other functions as needed.


## Examples

This section contains examples of the event payloads that are sent to the starlarktransform processor from each telemetry type. These examples can help you understand the structure of the telemetry events and how to modify them.

##### Log Event Payload Example:

```json
{
	"resourceLogs": [{
		"resource": {
			"attributes": [{
				"key": "log.file.name",
				"value": {
					"stringValue": "test.log"
				}
			}]
		},
		"scopeLogs": [{
			"scope": {},
			"logRecords": [{
				"observedTimeUnixNano": "1694127596456358000",
				"body": {
					"stringValue": "2023-09-06T01:09:24.045+0200 INFO starting app.",
					"attributes": [{

						"key": "app",
						"value": {
							"stringValue": "dev"
						}
					}],
					"traceId": "",
					"spanId": ""
				}
			}]
		}]
	}]
}
```

View the the log.proto type definition [here](https://github.com/open-telemetry/opentelemetry-proto/blob/main/opentelemetry/proto/logs/v1/logs.proto)

##### Metric Event Payload Example:

```json
{
	"resourceMetrics": [{
		"resource": {},
		"scopeMetrics": [{
			"scope": {
				"name": "otelcol/hostmetricsreceiver/memory",
				"version": "0.84.0"
			},
			"metrics": [{
				"name": "system.memory.usage",
				"description": "Bytes of memory in use.",
				"unit": "By",
				"sum": {
					"dataPoints": [{
							"attributes": [{
								"key": "state",
								"value": {
									"stringValue": "used"
								}
							}],
							"startTimeUnixNano": "1694171569000000000",
							"timeUnixNano": "1694189699786689531",
							"asInt": "1874247680"
						},
						{
							"attributes": [{
								"key": "state",
								"value": {
									"stringValue": "free"
								}
							}],
							"startTimeUnixNano": "1694171569000000000",
							"timeUnixNano": "1694189699786689531",
							"asInt": "29214199808"
						}
					],
					"aggregationTemporality": 2
				}
			}]
		}],
		"schemaUrl": "https://opentelemetry.io/schemas/1.9.0"
	}]
} 
```
View the metric.proto type definition [here](https://github.com/open-telemetry/opentelemetry-proto/blob/main/opentelemetry/proto/metrics/v1/metrics.proto)

##### Trace Event Payload Example:

```json
{
	"resourceSpans": [{
		"resource": {
			"attributes": [{
					"key": "telemetry.sdk.language",
					"value": {
						"stringValue": "python"
					}
				},
				{
					"key": "telemetry.sdk.name",
					"value": {
						"stringValue": "opentelemetry"
					}
				},
				{
					"key": "telemetry.sdk.version",
					"value": {
						"stringValue": "1.19.0"
					}
				},
				{
					"key": "telemetry.auto.version",
					"value": {
						"stringValue": "0.40b0"
					}
				},
				{
					"key": "service.name",
					"value": {
						"stringValue": "unknown_service"
					}
				}
			]
		},
		"scopeSpans": [{
			"scope": {
				"name": "opentelemetry.instrumentation.flask",
				"version": "0.40b0"
			},
			"spans": [{
				"traceId": "9cb5bf738137b2248dc7b20445ec2e1c",
				"spanId": "88079ad5c94b5b13",
				"parentSpanId": "",
				"name": "/roll",
				"kind": 2,
				"startTimeUnixNano": "1694388218052842000",
				"endTimeUnixNano": "1694388218053415000",
				"attributes": [{
						"key": "http.method",
						"value": {
							"stringValue": "GET"
						}
					},
					{
						"key": "http.server_name",
						"value": {
							"stringValue": "0.0.0.0"
						}
					},
					{
						"key": "http.scheme",
						"value": {
							"stringValue": "http"
						}
					},
					{
						"key": "net.host.port",
						"value": {
							"intValue": "5001"
						}
					},
					{
						"key": "http.host",
						"value": {
							"stringValue": "localhost:5001"
						}
					},
					{
						"key": "http.target",
						"value": {
							"stringValue": "/roll"
						}
					},
					{
						"key": "net.peer.ip",
						"value": {
							"stringValue": "127.0.0.1"
						}
					},
					{
						"key": "http.user_agent",
						"value": {
							"stringValue": "curl/7.87.0"
						}
					},
					{
						"key": "net.peer.port",
						"value": {
							"intValue": "52365"
						}
					},
					{
						"key": "http.flavor",
						"value": {
							"stringValue": "1.1"
						}
					},
					{
						"key": "http.route",
						"value": {
							"stringValue": "/roll"
						}
					},
					{
						"key": "http.status_code",
						"value": {
							"intValue": "200"
						}
					}
				],
				"status": {}
			}]
		}]
	}]
}
```

View the trace.proto type definition [here](https://github.com/open-telemetry/opentelemetry-proto/blob/main/opentelemetry/proto/trace/v1/trace.proto).


## Full Configuration Example

For following configuration example demonstrates the starlarktransform processor telemetry events for logs, metrics and traces.

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: "0.0.0.0:4318"
      grpc:
        endpoint: "0.0.0.0:4317"

  filelog:
    start_at: beginning
    include_file_name: true
    include: 
      - $LOGFILE

    operators:
      - type: move
        from: attributes["log.file.name"]
        to: resource["log.file.name"]

      - type: add
        field: attributes.app
        value: dev

processors:

  # - change resource attribute log.file.name to source.log
  # - add resource attribute cluster: dev
  # - filter out any logs that contain the word password
  # - add an attribute to each log: language: golang
  starlarktransform/logs:
    code: |
      def transform(event):
        event = json.decode(event)
        # edit resource attributes
        for data in event['resourceLogs']:
          for attr in data['resource']['attributes']:
            attr['value']['stringValue'] = 'source.log'

        # filter/delete logs
        for data in event['resourceLogs']:
          for slog in  data['scopeLogs']:
            slog['logRecords'] = [ lr for lr in slog['logRecords'] if 'internal' not in lr['body']['stringValue']]
            
            # add an attribute to each log
            for lr in slog['logRecords']:
              lr['attributes'].append({
                'key': 'language',
                'value': {
                  'stringValue': 'golang'
                }})
                
        return event
  # - print event received to otel runtime log
  # - if there are no resources, add a resource attribute source starlarktransform
  # - prefix each metric name with starlarktransform
  starlarktransform/metrics:
    code: |
      def transform(event):
        print("received event", event)
        event = json.decode(event)
        for md in event['resourceMetrics']:
          # if resources are empty
          if not md['resource']:
            md['resource'] = {
              'attributes': [
                {
                  "key": "source",
                  "value": {
                    "stringValue": "starlarktransform"
                  }
                }
              ]
            }

          # prefix each metric name with starlarktransform
          for sm in md['scopeMetrics']:
            for m in sm['metrics']:
              m['name'] = 'starlarktransform.' + m['name']

        return event

  # - add resource attribute source starlarktransform
  # - filter out any spans with http.target /roll attribute
  starlarktransform/traces:
    code: |
      def transform(event):
        event = json.decode(event)
        for td in event['resourceSpans']:
          # add resource attribute
          td['resource']['attributes'].append({
            'key': 'source',
            'value': {
              'stringValue': 'starlarktransform'
            }
          })

          # filter spans with http.target /roll attribute
          has_roll = lambda attrs: [a for a in attrs if a['key'] == 'http.target' and a['value']['stringValue'] == '/cats']
          for sd in td['scopeSpans']:
            sd['spans'] = [
              s for s in sd['spans']
              if not has_roll(s['attributes'])
            ]
        return event
exporters:
  logging:
    verbosity: detailed


service:
  pipelines:
    logs:
      receivers:
      - filelog
      processors:
      - starlarktransform/logs
      exporters:
      - logging

    metrics:
      receivers:
      - otlp
      processors:
      - starlarktransform/metrics
      exporters:
      - logging

    traces:
      receivers:
      - otlp
      processors:
      - starlarktransform/traces
      exporters:
      - logging
```


### Warnings

The starlarktransform processor allows you to modify all aspects of your telemetry data. This can result in invalid or bad data being propogated if you are not careful. It is your responsibility to inspect the data and ensure it is valid.

