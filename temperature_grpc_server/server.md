# gRPC Temperature Server

Il server gRPC per il servizio `TemperatureService` gestisce richieste di temperatura da client gRPC. Questo server è configurato per ascoltare su una porta specificata e può essere interrotto in modo sicuro tramite segnali di sistema.

[Click here for the English version](#english-version)

## Requisiti

- Go 1.22 o superiore
- Chiave API di WeatherAPI configurata come variabile d'ambiente (`WEATHER_API_KEY`)

## Configurazione

Prima di eseguire il server, assicurati di avere configurato correttamente la chiave API di WeatherAPI come descritto nella guida di configurazione.

## Esecuzione del Server

Prima di lanciare il server, assicurati di avviare `docker-compose` con l'istanza di MongoDB. Segui questi passaggi:

1. Avvia l'istanza di MongoDB con Docker Compose:

```bash
docker-compose up -d
```

2. Apri un terminale e naviga alla directory `temperature_grpc_server`:

```bash
cd temperature_grpc_server/cmd/server
```

3. Esegui il server:

```bash
go run main.go
```

## Funzionamento del Codice

Il codice del server esegue le seguenti operazioni:

1. **Configurazione**:
- Impostazione dell'indirizzo del server gRPC (`grpcAddress`).
- Definizione del protocollo di rete (`tcp`).

```go
const grpcAddress = ":50051"
const tcp = "tcp"
```

2. **Gestione dei Segnali di Interruzione**:
- Creazione di un contesto che gestisce i segnali di interruzione (es. `SIGTERM` e `os.Interrupt`) per consentire una chiusura sicura del processo.

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
```

3. **Avvio del Server gRPC**:
- Invocazione della funzione `startGRPCServer(ctx, grpcAddress)` che avvia il server gRPC per gestire le richieste del servizio `TemperatureService`.

```go
if err := startGRPCServer(ctx, grpcAddress); err != nil {
log.Fatalf("Failed to start gRPC server: %v", err)
}
```

4. **Funzione di Avvio del Server gRPC**:
- Creazione di un listener sulla porta specificata.
- Inizializzazione di un nuovo server gRPC.
- Registrazione del server del servizio `TemperatureService`.
- Avvio del server gRPC per gestire le richieste in arrivo.
- Gestione della chiusura sicura del server gRPC quando il contesto viene annullato.

```go
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
```

## Interruzione del Server

Per interrompere il server:

- Premere `Ctrl+C` nel terminale in cui il server è in esecuzione.

Il server gestirà l'interruzione in modo sicuro, annullando il contesto e terminando il processo.

# English Version

The gRPC server for the `TemperatureService` handles temperature requests from gRPC clients. This server is configured to listen on a specified port and can be safely interrupted via system signals.

## Requirements

- Go 1.22 or higher
- WeatherAPI API key configured as an environment variable (`WEATHER_API_KEY`)

## Configuration

Before running the server, ensure that the WeatherAPI API key is configured correctly as described in the setup guide.

## Running the Server

Before launching the server, make sure to start `docker-compose` with the MongoDB instance. Follow these steps:

1. Start the MongoDB instance with Docker Compose:

```bash
docker-compose up -d
```

2. Open a terminal and navigate to the `temperature_grpc_server` directory:

```bash
cd temperature_grpc_server/cmd/server
```

3. Run the server:

```bash
go run main.go
```

## Code Operation

The server code performs the following operations:

1. **Configuration**:
- Setting the gRPC server address (`grpcAddress`).
- Defining the network protocol (`tcp`).

```go
const grpcAddress = ":50051"
const tcp = "tcp"
```

2. **Handling Interrupt Signals**:
- Creating a context that handles interrupt signals (e.g., `SIGTERM` and `os.Interrupt`) to allow safe process shutdown.

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
```

3. **Starting the gRPC Server**:
- Calling the `startGRPCServer(ctx, grpcAddress)` function that starts the gRPC server to handle `TemperatureService` requests.

```go
if err := startGRPCServer(ctx, grpcAddress); err != nil {
log.Fatalf("Failed to start gRPC server: %v", err)
}
```

4. **gRPC Server Start Function**:
- Creating a listener on the specified port.
- Initializing a new gRPC server.
- Registering the `TemperatureService` server.
- Starting the gRPC server to handle incoming requests.
- Managing safe shutdown of the gRPC server when the context is canceled.

```go
func startGRPCServer(ctx context.Context, address string) error {
listener, err := net.Listen(tcp, address)
if err != nil {
return fmt.Errorf("failed to listen on %s: %w", address, err)
}
defer listener.Close()

grpcServer := grpc.NewServer()

// Stop if context is done
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
```

## Stopping the Server

To stop the server:

- Press `Ctrl+C` in the terminal where the server is running.

The server will handle the interruption safely, canceling the context and terminating the process.
