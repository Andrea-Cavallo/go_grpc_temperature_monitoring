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
	"fmt"
	pb "go_with_grpc/pkg/temperature"
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

	if err := startGRPCServer(ctx, grpcAddress); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	<-ctx.Done()

	log.Println("Server  stopped.")
}

// startGRPCServer starts a gRPC server for handling temperature service requests.
// It takes a context and an address as input parameters.
// It creates a listener on the provided address and initializes a new gRPC server.
// It registers the TemperatureService server and starts the gRPC server to handle incoming requests.
// If the context is canceled, the server is gracefully stopped.
// It returns an error if any error occurs during the process.
func startGRPCServer(ctx context.Context, address string) error {
	listener, err := net.Listen(tcp, address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()

	// Stoppa se il context viene terminato
	go func() {
		<-ctx.Done()
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	pb.RegisterTemperatureServiceServer(grpcServer, service.NewServer())
	log.Printf("gRPC server listening on %s\n", address)

	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}
