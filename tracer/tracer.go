package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type TracerConfig struct {
	// please provide JaegerCollectorURL if jaeger is used
	UseJaeger          bool
	JaegerCollectorURL string
	Environment        string
	ServiceName        string
}

var tracerCfg *TracerConfig

func Init(cfg *TracerConfig) (err error) {

	setGlobalTracerCfg(cfg)
	if tracerCfg.UseJaeger {
		err = initJaeger(tracerCfg.ServiceName, tracerCfg.Environment, tracerCfg.JaegerCollectorURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func setGlobalTracerCfg(cfg *TracerConfig) {
	tracerCfg = cfg
}

func initJaeger(service string, environment, url string) error {
	_, err := tracerProviderJaeger(service, environment, url)
	if err != nil {
		return err
	}

	return nil
}

// tracerProviderJaeger returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProviderJaeger(service, environment, url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
		)),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	return tp, nil
}

func StartSpan(ctx context.Context, name string, attributes map[string]string) (context.Context, trace.Span) {

	ctx, span := otel.GetTracerProvider().Tracer(name).Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(semconv.NetTransportTCP),
	)
	for k, v := range attributes {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return ctx, span
}
