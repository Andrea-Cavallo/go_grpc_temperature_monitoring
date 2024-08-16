# gRPC Temperature Client

Il client gRPC per il polling delle temperature utilizza il servizio `TemperatureService` per ottenere le temperature correnti di una località specificata. Questo client è configurato per eseguire il polling a intervalli regolari e può essere interrotto in modo sicuro tramite segnali di sistema.

## Requisiti

- Go 1.22 o superiore
- Chiave API di WeatherAPI configurata come variabile d'ambiente (`WEATHER_API_KEY`)

## Configurazione

Prima di eseguire il client, assicurati di avere configurato correttamente la chiave API di WeatherAPI come descritto nella guida di configurazione.

## Esecuzione del Client

Per eseguire il client, segui questi passaggi:

1. Apri un terminale e naviga alla directory `temperature_grpc_client`:

    ```bash
    cd temperature_grpc_client/cmd/client
    ```

2. Esegui il client:

    ```bash
    go run main.go
    ```

## Funzionamento del Codice

Il codice del client esegue le seguenti operazioni:

1. **Configurazione**:
    - Impostazione dell'indirizzo del server gRPC (`grpcServerAddress`).
    - Definizione dell'intervallo di polling (`pollingInterval`).
    - Specificazione della località per il polling delle temperature (`location`).

    ```go
    const (
        grpcServerAddress = "localhost:50051"
        pollingInterval   = 10 * time.Second
        location          = "Rome"
    )
    ```

2. **Gestione dei Segnali di Interruzione**:
    - Creazione di un contesto che gestisce i segnali di interruzione (es. `SIGTERM` e `os.Interrupt`) per consentire una chiusura sicura del processo.

    ```go
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()
    ```

3. **Avvio del Polling**:
    - Invocazione della funzione `startPolling(ctx)` che avvia il polling delle temperature.

    ```go
    startPolling(ctx)
    ```

4. **Funzione di Polling**:
    - Creazione del client gRPC utilizzando la funzione `NewClient`.
    - Configurazione di un ticker per eseguire il polling a intervalli regolari.
    - Esecuzione del polling finché il contesto non viene cancellato, interrompendo il polling in modo sicuro quando viene ricevuto un segnale di interruzione.

    ```go
    func startPolling(ctx context.Context) {
        client, err := service.NewClient(grpcServerAddress)
        if (err != nil) {
            log.Fatalf("Errore nella creazione del client gRPC: %v", err)
        }

        ticker := time.NewTicker(pollingInterval)
        defer ticker.Stop()

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
    ```

## Interruzione del Client

Per interrompere il client:

- Premere `Ctrl+C` nel terminale in cui il client è in esecuzione.

Il client gestirà l'interruzione in modo sicuro, annullando il polling e terminando il processo.

