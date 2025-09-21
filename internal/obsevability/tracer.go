package obsevability

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

type TracerProvider interface {
	Tracer(name string) trace.Tracer
}

type DefaultTracerProvider struct{}

func (d *DefaultTracerProvider) Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// InitTracer will init Tracer globally
func InitTracer(ctx context.Context) func(context.Context) error {
	exp, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("otel-collector:4318"),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithURLPath("/v1/traces"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Define atributos do recurso (inclui o nome da aplicação)
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("budget-tracker-api-v2"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Cria o TracerProvider com o recurso
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown
}
