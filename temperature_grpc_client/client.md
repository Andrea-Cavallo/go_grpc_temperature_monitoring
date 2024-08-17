# gRPC Temperature Client

Il client gRPC per il polling delle temperature utilizza il servizio `TemperatureService` per ottenere le temperature correnti di una località specificata. Questo client è configurato per eseguire il polling a intervalli regolari e può essere interrotto in modo sicuro tramite segnali di sistema.

[Click here for the English version](#english-version)

## Requisiti

- Go 1.22 o superiore
- Chiave API di WeatherAPI configurata come variabile d'ambiente (`WEATHER_API_KEY`)

## Configurazione

Prima di eseguire il client, assicurati di avere configurato correttamente la chiave API di WeatherAPI come descritto nella guida di configurazione.

## Esecuzione del Client

Per eseguire il client, segui questi passaggi:

1. Prima di lanciare il client, assicurati di aver avviato il server:

    ```bash
    cd temperature_grpc_server/cmd/server
    go run main.go
    ```

2. Apri un terminale e naviga alla directory `temperature_grpc_client`:

    ```bash
    cd temperature_grpc_client/cmd/client
    ```

3. Esegui il client:

    ```bash
    go run main.go
    ```

La chiusura del client (premendo `Ctrl+C`) farà sì che il programma disegni un grafico con la variazione delle temperature registrate durante l'esecuzione, tramite il file `plot.go`.

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
        if err != nil {
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

Il client gestirà l'interruzione in modo sicuro, annullando il polling, terminando il processo e disegnando un grafico con la variazione delle temperature tramite `plot.go`.

# English Version

The gRPC client for polling temperatures uses the `TemperatureService` to get current temperatures for a specified location. This client is configured to poll at regular intervals and can be safely interrupted via system signals.

## Requirements

- Go 1.22 or higher
- WeatherAPI API key configured as an environment variable (`WEATHER_API_KEY`)

## Configuration

Before running the client, ensure that the WeatherAPI API key is configured correctly as described in the setup guide.

## Running the Client

To run the client, follow these steps:

1. Before launching the client, make sure to start the server:

    ```bash
    cd temperature_grpc_server/cmd/server
    go run main.go
    ```

2. Open a terminal and navigate to the `temperature_grpc_client` directory:

    ```bash
    cd temperature_grpc_client/cmd/client
    ```

3. Run the client:

    ```bash
    go run main.go
    ```

Stopping the client (by pressing `Ctrl+C`) will cause the program to draw a graph with the temperature variations recorded during the application's runtime, using the `plot.go`.

## Code Operation

The client code performs the following operations:

1. **Configuration**:
    - Setting the gRPC server address (`grpcServerAddress`).
    - Defining the polling interval (`pollingInterval`).
    - Specifying the location for temperature polling (`location`).

    ```go
    const (
        grpcServerAddress = "localhost:50051"
        pollingInterval   = 10 * time.Second
        location          = "Rome"
    )
    ```

2. **Handling Interrupt Signals**:
    - Creating a context that handles interrupt signals (e.g., `SIGTERM` and `os.Interrupt`) to allow safe process shutdown.

    ```go
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()
    ```

3. **Start Polling**:
    - Invoking the `startPolling(ctx)` function that starts temperature polling.

    ```go
    startPolling(ctx)
    ```

4. **Polling Function**:
    - Creating the gRPC client using the `NewClient` function.
    - Setting up a ticker to perform polling at regular intervals.
    - Performing polling until the context is canceled, safely stopping polling when an interrupt signal is received.

    ```go
    func startPolling(ctx context.Context) {
        client, err := service.NewClient(grpcServerAddress)
        if err != nil {
            log.Fatalf("Error creating the gRPC client: %v", err)
        }

        ticker := time.NewTicker(pollingInterval)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                log.Println("Interrupt signal received, stopping polling.")
                return
            case <-ticker.C:
                log.Println("Polling temperatures...")
                client.GetCurrentTemperature(location)
            }
        }
    }
    ```

## Stopping the Client

To stop the client:

- Press `Ctrl+C` in the terminal where the client is running.

The client will handle the interruption safely, cancel the polling, terminate the process, and draw a graph with temperature variations using `plot.go`.
