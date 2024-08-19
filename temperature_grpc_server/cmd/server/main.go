package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_with_grpc/pkg/logger"
	"go_with_grpc/pkg/telemetry"
	"go_with_grpc/pkg/temperature"
	"go_with_grpc/temperature_grpc_server/mongodb/config"
	"go_with_grpc/temperature_grpc_server/service"
	"google.golang.org/grpc"
	"net"
	"net/http"
	_ "net/http/pprof" // Importa il package per abilitare pprof
	"os"
	"os/signal"
	"syscall"
)

const (
	grpcAddress = ":50051"
	tcp         = "tcp"
)

var log = logrus.WithFields(logrus.Fields{
	"service": "temperature-grpc-server",
	"env":     "development",
})

func main() {
	lokiHook, err := logger.NewLokiHook("http://localhost:3100/loki/api/v1/push", logrus.AllLevels)
	if err != nil {
		log.Errorf("Failed to create Loki hook: %v", err)
	}
	logrus.AddHook(lokiHook)
	defer lokiHook.Close()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("Starting pprof server on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Errorf("Failed to start pprof server: %v", err)
		}
	}()

	shutdownTelemetry := initializeTelemetry()
	defer func() {
		if err := shutdownTelemetry(context.Background()); err != nil {
			log.Errorf("Failed to shut down OpenTelemetry: %v", err)
		}
	}()

	log.Info("Starting OpenTelemetry..")

	if err := startGRPCServer(ctx, grpcAddress); err != nil {
		log.Errorf("Failed to start gRPC server: %v", err)
	}
	log.Info("Starting GRPC server..")

	<-ctx.Done()
	log.Info("Server stopped.")
}

func initializeTelemetry() func(context.Context) error {
	return telemetry.InitTelemetry()
}

// startGRPCServer avvia un server gRPC sull'indirizzo specificato.
func startGRPCServer(ctx context.Context, address string) error {
	listener, err := net.Listen(tcp, address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()

	go func() {
		<-ctx.Done()
		log.Info("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		log.Info("gRPC server was shutdown.")
		config.CloseMongoClient()
		log.Info("Closing MongoDB client... done. Goodbye! ^^")
	}()

	//Senza questa registrazione, il server gRPC non sarebbe in grado di gestire
	//le chiamate RPC indirizzate al servizio TemperatureServiceServer.
	temperature.RegisterTemperatureServiceServer(grpcServer, service.NewServer())

	log.Infof("gRPC server listening on %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}
