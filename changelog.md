# Changelog

Tutte le modifiche rilevanti apportate a questo progetto saranno documentate in questo file.

## [1.0.0] - 2024-08-17

### Aggiunto
- Integrazione di gRPC per la comunicazione tra client e server con supporto per i quattro tipi di RPC: Unario, Streaming Server, Streaming Client, e Streaming Bidirezionale.
- Implementazione del servizio `TemperatureService` per ottenere i dati della temperatura attuale di una località specificata.
- Configurazione del server gRPC e del client gRPC per inviare e ricevere richieste gRPC.
- Supporto per la creazione di grafici della temperatura alla chiusura del client.
- Documentazione dettagliata per l'installazione e l'uso del progetto, inclusa la configurazione di chiavi API e variabili d'ambiente.
- Supporto per l'invio di avvisi tramite Twilio con configurazione opzionale.
- Docker Compose configurato con servizi per MongoDB, Jaeger, OpenTelemetry Collector, Prometheus, Grafana, Loki, Node Exporter, cAdvisor, e K6 per il test delle prestazioni.
- Integrazione con Prometheus per la raccolta delle metriche e Grafana per la visualizzazione delle stesse.
- Supporto per tracing distribuito con Jaeger tramite OpenTelemetry.

### Modificato
- Aggiornamento della documentazione per includere dettagli sull'uso di Grafana con credenziali predefinite (`admin/admin`).
- Miglioramento della configurazione Docker Compose per una maggiore stabilità e facilità di utilizzo.

### Risolto
- Problemi di connessione tra client e server gRPC ora risolti, garantendo un flusso di dati stabile e continuo.
- Risolto un problema con la mancata individuazione del file `.proto` da parte di K6 durante l'esecuzione dei test delle prestazioni.

### Rilascio Iniziale
- Rilascio iniziale del progetto con tutte le funzionalità di base per il monitoraggio delle temperature utilizzando gRPC e l'integrazione di strumenti di osservabilità.
