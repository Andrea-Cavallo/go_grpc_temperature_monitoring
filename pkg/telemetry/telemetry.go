package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTelemetry() func(context.Context) error {
	ctx := context.Background()

	// Configura l'exporter OTLP per inviare i tracciamenti a Jaeger
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("localhost:16685"), otlptracegrpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to create OTLP trace exporter: %v", err)
	}

	// Configura il TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("temperature-monitoring-service"),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown
}
