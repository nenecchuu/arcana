# tracer

`tracer` is a helper to send the distributed tracing metric based open telemetry for monitoring. Currently we still using as our main tracing.

## Usage Example
- Make sure you have set up your own jaeger collector, exporter and UI. For the starter, you can check this guide https://www.jaegertracing.io/docs/1.26/getting-started/#all-in-one

- Init jaeger first, for example like this
```go
    import "github.com/nenecchuu/arcana/tracer"

    func main(){
        // currently we are using jaeger only, in the future we will add more implementation, such as newrelic
        tracer.Init(&tracer.TracerConfig{
            UseJaeger:          true,
            JaegerCollectorURL: "http://localhost:14268/api/traces",
            Environment:        "production",
            ServiceName:        "trace-demo",
        })
        // some implementation here
    }
```

- Typically after you already init the tracer, we already set up the global tracer using provided config, so you only need to start the span like this 
```go
    ctx, span := tracer.StartSpan(r.Context(), "main-func", nil)
    defer span.End()
```

Or if you want to add some additional attribute on each span, you can add it like this 
```go
    ctx, span := tracer.StartSpan(r.Context(), "main-func", map[string]string{
        "key": "value",
    })
    defer span.End()
```

- For http server implementation, you can check the example on `/jaeger_example_with_chi_server`, just run the main func and test the endpoint there
