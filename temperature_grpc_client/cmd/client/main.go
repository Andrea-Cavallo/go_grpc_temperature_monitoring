package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go_with_grpc/pkg/logger"
	"go_with_grpc/pkg/telemetry"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_with_grpc/temperature_grpc_client/service"
)

const (
	grpcServerAddress = "localhost:50051"
	initialDelay      = 0 * time.Second
	pollingInterval   = 60 * time.Second
	location          = "Rome"
	outputFile        = "temp_graph.png"
)

func main() {
	lokiHook, err := logger.NewLokiHook("http://localhost:3100/loki/api/v1/push", logrus.AllLevels)
	if err != nil {
		logrus.Errorf("Failed to create Loki hook: %v", err)
	}
	logrus.AddHook(lokiHook)

	// Inizializza OpenTelemetry
	shutdown := telemetry.InitTelemetry()
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			logrus.Errorf("Failed to shut down OpenTelemetry: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := service.NewClient(grpcServerAddress)
	if err != nil {
		logrus.Errorf("Errore nella creazione del client gRPC: %v", err)
	}

	logrus.Println("Polling temperatures...")

	tracer := otel.Tracer("temperature-monitoring")
	pollCtx, span := tracer.Start(ctx, "Temperature Polling")
	client.GetCurrentTemperature(pollCtx, location)
	span.End()

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Println("Interruzione ricevuta, fermo il polling.")
			logrus.Printf("Disegno il grafico delle temperature... in %s", outputFile)
			client.PlotTemperatureGraph(outputFile)
			return
		case <-ticker.C:
			logrus.Println("Polling temperatures...")

			pollCtx, span := tracer.Start(ctx, "Temperature Polling")
			client.GetCurrentTemperature(pollCtx, location)
			span.End()
		}
	}
}
