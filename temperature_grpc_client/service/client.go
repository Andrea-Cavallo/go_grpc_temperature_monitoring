package service

import (
	"context"
	"go_with_grpc/pkg/temperature"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"google.golang.org/grpc"
)

// Client è una struttura che rappresenta il client gRPC per il servizio TemperatureService.
type Client struct {
	service temperature.TemperatureServiceClient
}

// NewClient crea una nuova istanza del client gRPC per il servizio TemperatureService.
// address è l'indirizzo del server gRPC (es. "localhost:50051").
func NewClient(address string) (*Client, error) {
	// Definisci un contesto con un timeout per la connessione gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Crea la connessione gRPC utilizzando il contesto
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Errore durante la connessione al server gRPC: %v", err)
		return nil, err
	}

	client := temperature.NewTemperatureServiceClient(conn)
	log.Println("Connessione al server gRPC riuscita")
	return &Client{service: client}, nil
}

// GetCurrentTemperature richiede la temperatura attuale di una località specificata.
func (c *Client) GetCurrentTemperature(location string) {
	log.Printf("Richiesta della temperatura attuale per la località: %s\n", location)

	resp, err := c.service.GetCurrentTemperature(context.Background(), &temperature.GetCurrentTemperatureRequest{
		Location: location,
	})
	if err != nil {
		log.Fatalf("Errore durante la richiesta della temperatura attuale: %v", err)
		return
	}

	log.Printf("Temperatura attuale per %s: %.2f°C, rilevata il %s\n", location, resp.Reading.Temperature, time.Unix(resp.Reading.Timestamp, 0).Format(time.RFC3339))
}
