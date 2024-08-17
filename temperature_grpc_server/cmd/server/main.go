package main

import (
	"context"
	"fmt"
	"go_with_grpc/pkg/telemetry"
	"go_with_grpc/pkg/temperature"
	"go_with_grpc/temperature_grpc_server/mongodb/config"
	"go_with_grpc/temperature_grpc_server/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const grpcAddress = ":50051"
const tcp = "tcp"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownTelemetry := initializeTelemetry()
	defer func() {
		if err := shutdownTelemetry(context.Background()); err != nil {
			log.Fatalf("Failed to shut down OpenTelemetry: %v", err)
		}
	}()

	if err := startGRPCServer(ctx, grpcAddress); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	<-ctx.Done()
	log.Println("Server stopped.")
}

func initializeTelemetry() func(context.Context) error {
	return telemetry.InitTelemetry()
}

func startGRPCServer(ctx context.Context, address string) error {
	listener, err := net.Listen(tcp, address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		log.Println("gRPC server stopped.")
		config.CloseMongoClient()
		log.Println("Cloing MongoDB client... done. Goodbye! ^^")
	}()

	temperature.RegisterTemperatureServiceServer(grpcServer, service.NewServer())
	log.Printf("gRPC server listening on %s\n", address)
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}
