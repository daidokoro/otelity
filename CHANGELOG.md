## Otelity

### v0.1.1 / 2024-12-18
- [Fix] `log` module in starlark processor now logs the correct log level per log function


### v0.1.0 / 2024-12-15 (Breaking)
- [Feat] add `emit` funciton to starlark processor
- [Feat] added `log` module to starlark processor
- [Feat] added `entrypoint` option to starlark config, allowing to specify the entry point of the starlark script/code
- [Fix] `json.decode` no longer required to load telemetry events in starlark processor

### v0.0.1 / 2024-12-07
- [Feat] initial release
    - starlark processor
