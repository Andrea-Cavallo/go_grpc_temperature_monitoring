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
	"go_with_grpc/temperature_grpc_client/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	grpcServerAddress = "localhost:50051"
	pollingInterval   = 10 * time.Second
	location          = "Rome"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startPolling(ctx)
}

// startPolling avvia il polling per le temperature correnti utilizzando un client gRPC.
// Il polling si ferma quando il contesto viene annullato.
func startPolling(ctx context.Context) {

	client, err := service.NewClient(grpcServerAddress)
	if err != nil {
		log.Fatalf("Errore nella creazione del client gRPC: %v", err)
	}

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	// Poll finché il contesto non viene cancellato
	for {
		select {
		case <-ctx.Done():
			log.Println("Interruzione ricevuta, fermo il polling.")
			return
		case <-ticker.C:
			log.Println("Polling temperatures...")
			client.GetCurrentTemperature(location)
		}
	}
}
