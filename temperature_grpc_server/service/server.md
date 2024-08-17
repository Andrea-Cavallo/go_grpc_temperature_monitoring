# gRPC Temperature Server

Il server gRPC per il servizio `TemperatureService` gestisce richieste di temperatura da client gRPC. Questo server è configurato per ascoltare su una porta specificata e può essere interrotto in modo sicuro tramite segnali di sistema.

## Requisiti

- Go 1.22 o superiore
- Chiave API di WeatherAPI configurata come variabile d'ambiente (`WEATHER_API_KEY`)

## Configurazione

Prima di eseguire il server, assicurati di avere configurato correttamente la chiave API di WeatherAPI come descritto nella guida di configurazione.

## Esecuzione del Server

Per eseguire il server, segui questi passaggi:

1. Apri un terminale e naviga alla directory `temperature_grpc_server`:

```bash
cd temperature_grpc_server/cmd/server
```

2. Esegui il server:

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
