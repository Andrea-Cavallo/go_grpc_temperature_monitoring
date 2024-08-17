/*
Copyright 2024 Andrea Cavallo

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go_with_grpc/pkg/telemetry"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_with_grpc/temperature_grpc_client/service"
)

const (
	grpcServerAddress = "localhost:50051"
	pollingInterval   = 25 * time.Second
	location          = "Rome"
	outputFile        = "temp_graph.png"
)

func main() {

	// Inizializza OpenTelemetry
	shutdown := telemetry.InitTelemetry()
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shut down OpenTelemetry: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := service.NewClient(grpcServerAddress)
	if err != nil {
		log.Fatalf("Errore nella creazione del client gRPC: %v", err)
	}

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Interruzione ricevuta, fermo il polling.")
			log.Printf("Disegno il grafico delle temperature... in %s...\n", outputFile)
			client.PlotTemperatureGraph(outputFile)
			return
		case <-ticker.C:
			log.Println("Polling temperatures...")

			tracer := otel.Tracer("temperature-monitoring")
			ctx, span := tracer.Start(ctx, "Temperature Polling")
			client.GetCurrentTemperature(ctx, location)
			span.End()
		}
	}
}
