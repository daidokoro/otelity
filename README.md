# Otelity (_like "Utility" but with "Otel"_)

A repo for my Open Telemetry components that don't quite fit with the main contrib distro. Can contain anything from new and experimental to modified versions of existing components.

## Components
- [starklarktransform processor](./processors/starlarktransformprocessor/README.md)


## How to use

In order use a component in this repo, you must build it into your own custom OpenTelemetry distribution using the ocb cli tool described [here](https://opentelemetry.io/docs/collector/custom-collector/#step-1---install-the-builder).

Simply add the comonent to your `build.yaml` file and run `ocb` CLI command to build your custom distribution.

__For eg.__


_build.yaml:_
```yaml
dist:
  name: myot
  output_path: ./myot
  description: otelity binary.
  version: 1.0.0
  # otelcol_version: 0.115.0

extensions:

exporters:
  - gomod: go.opentelemetry.io/collector/exporter/debugexporter v0.115.0
  - gomod: go.opentelemetry.io/collector/exporter/nopexporter v0.115.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/exporter/coralogixexporter v0.115.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.115.0


processors:
  - gomod: go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.115.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.115.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.115.0
  
  # add the otelity starlark processor here
  - gomod: github.com/daidokoro/otelity/processors/starlarkprocessor v0.1.0


receivers:
  - gomod: go.opentelemetry.io/collector/receiver/nopreceiver v0.115.0
  - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver v0.115.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/receiver/datadogreceiver v0.115.0


providers:
  - gomod: go.opentelemetry.io/collector/confmap/provider/envprovider v1.17.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/fileprovider v1.17.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/httpprovider v1.17.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/httpsprovider v1.17.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/yamlprovider v1.17.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/confmap/provider/s3provider v0.115.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/confmap/provider/secretsmanagerprovider v0.115.0
```

