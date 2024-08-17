package service

import (
	"context"
	"go_with_grpc/pkg/temperature"
	"go_with_grpc/temperature_grpc_client/plot"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"google.golang.org/grpc"
)

// Client è una struttura che rappresenta il client gRPC per il servizio TemperatureService.
type Client struct {
	service         temperature.TemperatureServiceClient
	temperatureData []plot.TemperatureData // Memorizza i dati di temperatura per il plotting
}

// NewClient crea una nuova istanza del client gRPC per il servizio TemperatureService.
// address è l'indirizzo del server gRPC (es. "localhost:50051").
func NewClient(address string) (*Client, error) {
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
	return &Client{
		service:         client,
		temperatureData: []plot.TemperatureData{},
	}, nil
}

// GetCurrentTemperature richiede la temperatura attuale di una località specificata.
// Prende in ingresso un contesto (ctx) per la propagazione del tracing.
func (c *Client) GetCurrentTemperature(ctx context.Context, location string) {
	log.Printf("Richiesta della temperatura attuale per la località: %s\n", location)

	resp, err := c.service.GetCurrentTemperature(ctx, &temperature.GetCurrentTemperatureRequest{
		Location: location,
	})
	if err != nil {
		log.Fatalf("Errore durante la richiesta della temperatura attuale: %v", err)
		return
	}

	// Memorizza i dati di temperatura per il plotting
	c.temperatureData = append(c.temperatureData, plot.TemperatureData{
		Timestamp: time.Unix(resp.Reading.Timestamp, 0),
		Value:     resp.Reading.Temperature,
	})

	log.Printf("Temperatura attuale per %s: %.2f°C, rilevata il %s\n", location, resp.Reading.Temperature, time.Unix(resp.Reading.Timestamp, 0).Format(time.RFC3339))
}

// PlotTemperatureGraph genera e salva un grafico delle temperature raccolte.
func (c *Client) PlotTemperatureGraph(outputFile string) {
	if err := plot.PlotTemperatureGraph(c.temperatureData, outputFile); err != nil {
		log.Fatalf("Errore durante la generazione del grafico: %v", err)
	} else {
		log.Printf("Grafico generato: %s\n", outputFile)
	}
}
