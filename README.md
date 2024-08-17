#  Real-Time Temperature Monitoring with Go and gRPC

## Introduzione

Questo progetto dimostra l'integrazione di servizi gRPC con il linguaggio Go,
fornendo un esempio completo di un server gRPC che serve dati di temperatura e un client gRPC che li consuma,
il tutto con una gestione avanzata delle metriche, connessione a MongoDB per persistere i dati, funzionalità di allerta tramite Twilio.
Il progetto utilizza Docker per facilitare l'impostazione dell'ambiente, inclusi la configurazione di MongoDB e altre dipendenze tipo jaeger.

<img src="./extra/img/GRPC.png" height="600" width="600">

[Click here for the English version](#english-version)

# Struttura del progetto


```text
├── api --> qua dentro il file .proto
│   ├── protos
│   │   ├── temperature.proto
│   ├── api.md
├── docker ---> tutto cio' relativo a docker
│   ├── docker-compose.yml
│   ├── init-mongo.js
│   ├── loki-config.yml
│   ├── otel-collector-config.yaml
│   ├── prometheus.yml
├── extra --> immagini del progetto
├── pkg --> pkg tutto quello che deve essere visibile e condiviso
│   ├── telemetry
│   ├── temperature
├── temperature_grpc_client -> il Client GRPC
│   ├── cmd
│   ├── plot
│   ├── service
│   ├── client.md
├── temperature_grpc_server --> il Server GRPC
│   ├── alert_twilio
│   │   ├── alert.go
│   ├── cmd
│   │   ├── server
│   │   │   ├── main.go
│   ├── mongodb --> connessione a MONGODB
│   │   ├── config
│   │   ├── temperature_repo
│   ├── service
│   │   ├── server.go
│   ├── server.md
├── .editorconfig
├── .gitignore
├── changelog.md
├── go.mod
├── go.sum
```

# Diagramma UML

![Diagram](<?xml version="1.0" encoding="us-ascii" standalone="no"?><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" contentStyleType="text/css" height="577px" preserveAspectRatio="none" style="width:743px;height:577px;background:#FFFFFF;" version="1.1" viewBox="0 0 743 577" width="743px" zoomAndPan="magnify"><defs/><g><rect fill="none" height="194.377" style="stroke:#000000;stroke-width:1.5;" width="716.5" x="20" y="220.7852"/><rect fill="none" height="19.291" style="stroke:#000000;stroke-width:1.5;" width="86" x="30" y="388.8711"/><line style="stroke:#181818;stroke-width:0.5;stroke-dasharray:5.0,5.0;" x1="41" x2="41" y1="83.6211" y2="494.7441"/><line style="stroke:#181818;stroke-width:0.5;stroke-dasharray:5.0,5.0;" x1="248.5" x2="248.5" y1="83.6211" y2="494.7441"/><line style="stroke:#181818;stroke-width:0.5;stroke-dasharray:5.0,5.0;" x1="384.5" x2="384.5" y1="83.6211" y2="494.7441"/><line style="stroke:#181818;stroke-width:0.5;stroke-dasharray:5.0,5.0;" x1="518.5" x2="518.5" y1="83.6211" y2="494.7441"/><line style="stroke:#181818;stroke-width:0.5;stroke-dasharray:5.0,5.0;" x1="654.5" x2="654.5" y1="83.6211" y2="494.7441"/><text fill="#000000" font-family="sans-serif" font-size="14" lengthAdjust="spacing" textLength="36" x="20" y="80.1074">Client</text><ellipse cx="41" cy="13.5" fill="#E2E2F0" rx="8" ry="8" style="stroke:#181818;stroke-width:0.5;"/><path d="M41,21.5 L41,48.5 M28,29.5 L54,29.5 M41,48.5 L28,63.5 M41,48.5 L54,63.5 " fill="none" style="stroke:#181818;stroke-width:0.5;"/><text fill="#000000" font-family="sans-serif" font-size="14" lengthAdjust="spacing" textLength="36" x="20" y="508.8516">Client</text><ellipse cx="41" cy="520.8652" fill="#E2E2F0" rx="8" ry="8" style="stroke:#181818;stroke-width:0.5;"/><path d="M41,528.8652 L41,555.8652 M28,536.8652 L54,536.8652 M41,555.8652 L28,570.8652 M41,555.8652 L54,570.8652 " fill="none" style="stroke:#181818;stroke-width:0.5;"/><rect fill="#E2E2F0" height="32.6211" rx="2.5" ry="2.5" style="stroke:#181818;stroke-width:0.5;" width="97" x="200.5" y="50"/><text fill="#000000" font-family="sans-serif" font-size="14" lengthAdjust="spacing" textLength="83" x="207.5" y="72.1074">gRPC Server</text><rect fill="#E2E2F0" height="32.6211" rx="2.5" ry="2.5" style="stroke:#181818;stroke-width:0.5;" width="97" x="200.5" y="493.7441"/><text fill="#000000" font-family="sans-serif" font-size="14" lengthAdjust="spacing" textLength="83" x="207.5" y="515.8516">gRPC Server</text><rect fill="#E2E2F0" height="32.6211" rx="2.5" ry="2.5" style="stroke:#181818;stroke-width:0.5;" width="141" x="314.5" y="50"/><text fill="#000000" font-family="sans-serif" font-size="14" lengthAdjust="spacing" textLength="127" x="321.5" y="72.1074">MongoDB Database</text><rect fill="#E2E2F0" height="32.6211" rx="2.5" ry="2.5" style="stroke:#181818;stroke-width:0.5;" width="141" x="314.5" y="493.7441"/><text fill="#000000" font-family="sans-serif" font-size="14" lengthAdjust="spacing" textLength="127" x="321.5" y="515.8516">MongoDB Database</text><rect fill="#E2E2F0" height="32.6211" rx="2.5" ry="2.5" style="stroke:#181818;stroke-width:0.5;" width="107" x="465.5" y="50"/><text fill="#000000" font-family="sans-serif" font-size="14" lengthAdjust="spacing" textLength="93" x="472.5" y="72.1074">Twilio SMS API</text><rect fill="#E2E2F0" height="32.6211" rx="2.5" ry="2.5" style="stroke:#181818;stroke-width:0.5;" width="107" x="465.5" y="493.7441"/><text fill="#000000" font-family="sans-serif" font-size="14" lengthAdjust="spacing" textLength="93" x="472.5" y="515.8516">Twilio SMS API</text><rect fill="#E2E2F0" height="32.6211" rx="2.5" ry="2.5" style="stroke:#181818;stroke-width:0.5;" width="144" x="582.5" y="50"/><text fill="#000000" font-family="sans-serif" font-size="14" lengthAdjust="spacing" textLength="130" x="589.5" y="72.1074">Weather Service API</text><rect fill="#E2E2F0" height="32.6211" rx="2.5" ry="2.5" style="stroke:#181818;stroke-width:0.5;" width="144" x="582.5" y="493.7441"/><text fill="#000000" font-family="sans-serif" font-size="14" lengthAdjust="spacing" textLength="130" x="589.5" y="515.8516">Weather Service API</text><polygon fill="#181818" points="237,112.9121,247,116.9121,237,120.9121,241,116.9121" style="stroke:#181818;stroke-width:1.0;"/><line style="stroke:#181818;stroke-width:1.0;" x1="41" x2="243" y1="116.9121" y2="116.9121"/><text fill="#000000" font-family="sans-serif" font-size="13" lengthAdjust="spacing" textLength="135" x="48" y="111.6494">Send Weather Request</text><polygon fill="#181818" points="642.5,144.2031,652.5,148.2031,642.5,152.2031,646.5,148.2031" style="stroke:#181818;stroke-width:1.0;"/><line style="stroke:#181818;stroke-width:1.0;" x1="249" x2="648.5" y1="148.2031" y2="148.2031"/><text fill="#000000" font-family="sans-serif" font-size="13" lengthAdjust="spacing" textLength="163" x="256" y="142.9404">Fetch Current Weather Data</text><polygon fill="#181818" points="260,175.4941,250,179.4941,260,183.4941,256,179.4941" style="stroke:#181818;stroke-width:1.0;"/><line style="stroke:#181818;stroke-width:1.0;stroke-dasharray:2.0,2.0;" x1="254" x2="653.5" y1="179.4941" y2="179.4941"/><text fill="#000000" font-family="sans-serif" font-size="13" lengthAdjust="spacing" textLength="172" x="266" y="174.2314">Provide Current Weather Data</text><path d="M20,220.7852 L81,220.7852 L81,230.0762 L71,240.0762 L20,240.0762 L20,220.7852 " fill="#EEEEEE" style="stroke:#000000;stroke-width:1.5;"/><rect fill="none" height="194.377" style="stroke:#000000;stroke-width:1.5;" width="716.5" x="20" y="220.7852"/><text fill="#000000" font-family="sans-serif" font-size="13" font-weight="bold" lengthAdjust="spacing" textLength="16" x="35" y="235.8135">alt</text><text fill="#000000" font-family="sans-serif" font-size="11" font-weight="bold" lengthAdjust="spacing" textLength="162" x="96" y="234.6553">[Weather Conditions Not Met]</text><polygon fill="#181818" points="642.5,259.3672,652.5,263.3672,642.5,267.3672,646.5,263.3672" style="stroke:#181818;stroke-width:1.0;"/><line style="stroke:#181818;stroke-width:1.0;" x1="249" x2="648.5" y1="263.3672" y2="263.3672"/><text fill="#000000" font-family="sans-serif" font-size="13" lengthAdjust="spacing" textLength="176" x="256" y="258.1045">Poll for Updated Weather Data</text><polygon fill="#181818" points="260,290.6582,250,294.6582,260,298.6582,256,294.6582" style="stroke:#181818;stroke-width:1.0;"/><line style="stroke:#181818;stroke-width:1.0;stroke-dasharray:2.0,2.0;" x1="254" x2="653.5" y1="294.6582" y2="294.6582"/><text fill="#000000" font-family="sans-serif" font-size="13" lengthAdjust="spacing" textLength="174" x="266" y="289.3955">Return Updated Weather Data</text><line style="stroke:#000000;stroke-width:1.0;stroke-dasharray:2.0,2.0;" x1="20" x2="736.5" y1="303.6582" y2="303.6582"/><text fill="#000000" font-family="sans-serif" font-size="11" font-weight="bold" lengthAdjust="spacing" textLength="141" x="25" y="315.5283">[Weather Conditions Met]</text><polygon fill="#181818" points="507,338.5801,517,342.5801,507,346.5801,511,342.5801" style="stroke:#181818;stroke-width:1.0;"/><line style="stroke:#181818;stroke-width:1.0;" x1="249" x2="513" y1="342.5801" y2="342.5801"/><text fill="#000000" font-family="sans-serif" font-size="13" lengthAdjust="spacing" textLength="131" x="256" y="337.3174">Send SMS Notification</text><polygon fill="#181818" points="260,369.8711,250,373.8711,260,377.8711,256,373.8711" style="stroke:#181818;stroke-width:1.0;"/><line style="stroke:#181818;stroke-width:1.0;stroke-dasharray:2.0,2.0;" x1="254" x2="518" y1="373.8711" y2="373.8711"/><text fill="#000000" font-family="sans-serif" font-size="13" lengthAdjust="spacing" textLength="128" x="266" y="368.6084">SMS Notification Sent</text><path d="M30,388.8711 L111,388.8711 L111,398.1621 L101,408.1621 L30,408.1621 L30,388.8711 " fill="#EEEEEE" style="stroke:#000000;stroke-width:1.5;"/><rect fill="none" height="19.291" style="stroke:#000000;stroke-width:1.5;" width="86" x="30" y="388.8711"/><text fill="#000000" font-family="sans-serif" font-size="13" font-weight="bold" lengthAdjust="spacing" textLength="36" x="45" y="403.8994">break</text><polygon fill="#181818" points="373,441.4531,383,445.4531,373,449.4531,377,445.4531" style="stroke:#181818;stroke-width:1.0;"/><line style="stroke:#181818;stroke-width:1.0;" x1="249" x2="379" y1="445.4531" y2="445.4531"/><text fill="#000000" font-family="sans-serif" font-size="13" lengthAdjust="spacing" textLength="112" x="256" y="440.1904">Log Operation Data</text><polygon fill="#181818" points="52,472.7441,42,476.7441,52,480.7441,48,476.7441" style="stroke:#181818;stroke-width:1.0;"/><line style="stroke:#181818;stroke-width:1.0;stroke-dasharray:2.0,2.0;" x1="46" x2="248" y1="476.7441" y2="476.7441"/><text fill="#000000" font-family="sans-serif" font-size="13" lengthAdjust="spacing" textLength="184" x="58" y="471.4814">Return Weather Request Result</text><!--SRC=[bP91ReCm44NtFiKe-ro0HQMMg5AbqY9erIKozeHOTMng3ydrQvic18fKTO60uVtlFynmGTP1hXKQ0kv1f2VMSqTAg9w7wgQhoXngT2TqcyXqhXUDJ-zpiB2cixi7s77YBVUOw6CiohQHxHn-EokjBDH53VVbKq8fbHf_z7Fq0paTqAKCIk2LFu_rb2NK6HlJm_RkDcCMT4R2nQJ2hm4ziHiY2XPU6JmYYGFaWtFHygAEhR452zlbWQMp9oLnNSsqjJsKLcjbEjY7tt_v_DvBJkWPrvCajqOgKjPuUB441VwTF5edtz5F_ZDyXJGuiov6xN4jsagiQPW_0uYJ3YselH8ygjX4_0wnXAC5iTOX_sJ1B5tB4wTDo-3PTl3Qenl5AVkC3lYuBDT-5_iHxdxGz0q0]--></g></svg>)


# Esecuzione del Progetto

## 1. Prerequisiti

Prima di eseguire il progetto, è necessario ottenere una chiave API per utilizzare il servizio WeatherAPI. Questo servizio è completamente gratuito e puoi ottenere la tua chiave API seguendo questi passi:

1. Visita il sito [WeatherAPI](https://www.weatherapi.com/).
2. Registrati e segui le istruzioni per ottenere la tua chiave API.
3. Una volta ottenuta, salva la chiave API in un luogo sicuro.

## 2. Configurare l'API Key come Variabile d'Ambiente

Per mantenere la tua chiave API sicura e fuori dal codice sorgente, configurala come una variabile d'ambiente. Segui questi passi:

### Su Linux/MacOS:

```bash
export WEATHER_API_KEY=la-tua-api-key
```

### Su Windows (nel terminale PowerShell):

```powershell
$env:WEATHER_API_KEY="la-tua-api-key"
```

Assicurati che la variabile d'ambiente `WEATHER_API_KEY` sia impostata correttamente prima di procedere.

## 3. Build del Docker Compose

Assicurati di avere Docker installato e di avere accesso al file `docker-compose.yml`:

- Posizionati nella directory che contiene `docker-compose.yml`

## 4. Opzionale - Configurazione di Twilio

Se desideri inviare avvisi tramite Twilio, segui i seguenti passi per creare un account su Twilio e configurare le variabili d'ambiente:

Consulta: [Twilio WhatsApp Quickstart (Go)](https://www.twilio.com/docs/whatsapp/quickstart/go)

Configurare le seguenti variabili d'ambiente:

```bash
export TWILIO_ACCOUNT_SID=il-tuo-account-sid
export TWILIO_AUTH_TOKEN=il-tuo-auth-token
export TWILIO_PHONE_NUMBER=il-tuo-numero-di-telefono
export ALERT_PHONE_NUMBER=il-numero-di-telefono-per-gli-avvisi
```

## 5. Eseguire il Server

Una volta configurata la chiave API, puoi avviare il server:

1. Naviga alla directory `temperature_grpc_server`:

    ```bash
    cd temperature_grpc_server/cmd/server
    ```

2. Esegui il server:

    ```bash
    go run main.go
    ```

Il server inizierà ad ascoltare le richieste gRPC.

## 6. Eseguire il Client

Con il server in esecuzione, ora puoi avviare il client per inviare richieste al server:

1. Apri un nuovo terminale e naviga alla directory `temperature_grpc_client`:

    ```bash
    cd temperature_grpc_client/cmd/client
    ```

2. Esegui il client:

    ```bash
    go run main.go
    ```

Il client invierà una richiesta gRPC al server e riceverà una risposta con i dati di temperatura.

## Grafico della Temperatura

Ogni volta che arresti l'applicazione, in particolare il client, vedrai disegnato un grafico che mostrerà la variazione della temperatura.

![Grafico della Temperatura](./temp_graph.png)

## Per il dettaglio del client e del server

- temperature_grpc_client/client.md
- temperature_grpc_server/server.md

## Approfondimenti sulla tecnologia scelta

gRPC può essere considerato un successore di RPC, ed è leggero in termini di peso. Google lo ha sviluppato per la comunicazione tra microservizi e altri sistemi che necessitano di interagire. Ci sono diversi vantaggi nell'usare gRPC.

## Vantaggi di gRPC

- **Utilizza Protocol Buffers (Protobuf)** invece del JSON.
- **Costruito su HTTP/2** invece che su HTTP 1.1.
- **Generazione di codice incorporata**.
- **Alte prestazioni**.
- **Sicurezza SSL**.

Oltre ai vantaggi chiave menzionati sopra, gRPC promuove un design migliore per la tua applicazione. gRPC è orientato all'API, a differenza di REST che è orientato alle risorse. È anche asincrono per impostazione predefinita, il che significa che non blocca il thread su richiesta e può servire milioni di richieste in parallelo, garantendo un'alta scalabilità.

## Vantaggi di gRPC rispetto a REST

gRPC è circa sette volte più veloce di REST nel ricevere dati e circa dieci volte più veloce di REST nel trasmettere dati per un payload specifico. Questo è principalmente dovuto al pacchettamento compatto dei Protocol Buffers e all'uso di HTTP/2 da parte di gRPC.


## Autore

Progetto sviluppato da [Andrea Cavallo].

## Licenza

Questo progetto è distribuito sotto la licenza Apache 2.0. Vedi il file LICENSE per ulteriori dettagli.

# English Version

# Super cool Go with GRPC

<img src="./extra/img/GRPC.png" height="600" width="600">

# Project Execution

## 1. Prerequisites

Before running the project, you need to obtain an API key to use the WeatherAPI service. This service is completely free, and you can get your API key by following these steps:

1. Visit the [WeatherAPI](https://www.weatherapi.com/) website.
2. Sign up and follow the instructions to get your API key.
3. Once obtained, save the API key in a safe place.

## 2. Configure the API Key as an Environment Variable

To keep your API key secure and out of the source code, configure it as an environment variable. Follow these steps:

### On Linux/MacOS:

```bash
export WEATHER_API_KEY=your-api-key
```

### On Windows (in PowerShell):

```powershell
$env:WEATHER_API_KEY="your-api-key"
```

Ensure the `WEATHER_API_KEY` environment variable is set correctly before proceeding.

## 3. Docker Compose Build

Ensure you have Docker installed and access to the `docker-compose.yml` file:

- Navigate to the directory that contains `docker-compose.yml`.

## 4. Optional - Twilio Configuration

If you wish to send alerts via Twilio, follow these steps to create a Twilio account and set up the environment variables:

Consult: [Twilio WhatsApp Quickstart (Go)](https://www.twilio.com/docs/whatsapp/quickstart/go)

Set the following environment variables:

```bash
export TWILIO_ACCOUNT_SID=your-account-sid
export TWILIO_AUTH_TOKEN=your-auth-token
export TWILIO_PHONE_NUMBER=your-phone-number
export ALERT_PHONE_NUMBER=alert-phone-number
```

## 5. Running the Server

Once the API key is configured, you can start the server:

1. Navigate to the `temperature_grpc_server` directory:

    ```bash
    cd temperature_grpc_server/cmd/server
    ```

2. Start the server:

    ```bash
    go run main.go
    ```

The server will start listening for gRPC requests.

## 6. Running the Client

With the server running, you can now start the client to send requests to the server:

1. Open a new terminal and navigate to the `temperature_grpc_client` directory:

    ```bash
    cd temperature_grpc_client/cmd/client
    ```

2. Start the client:

    ```bash
    go run main.go
    ```

The client will send a gRPC request to the server and receive a response with the temperature data.

## Temperature Graph

Whenever you shut down the application, especially the client, you will see a graph showing the temperature variation.

![Temperature Graph](./temp_graph.png)

## Detailed Client and Server Documentation

- temperature_grpc_client/client.md
- temperature_grpc_server/server.md

## Insights into the Chosen Technology

gRPC can be considered a successor to RPC and is lightweight. Google developed it for communication between microservices and other systems that need to interact. There are several advantages to using gRPC.

## Advantages of gRPC

- **Uses Protocol Buffers (Protobuf)** instead of JSON.
- **Built on HTTP/2** instead of HTTP 1.1.
- **Built-in Code Generation**.
- **High Performance**.
- **SSL Security**.

In addition to the key advantages mentioned above, gRPC promotes better design for your application. gRPC is API-oriented, unlike REST, which is resource-oriented. It is also asynchronous by default, meaning it does not block the thread on request and can serve millions of requests in parallel, ensuring high scalability.

## Advantages of gRPC over REST

gRPC is about seven times faster than REST in receiving data and about ten times faster than REST in transmitting data for a specific payload. This is primarily due to the compact packing of Protocol Buffers and the use of HTTP/2 by gRPC.


## Author

Project developed by [Andrea Cavallo].

## License

This project is licensed under the Apache 2.0 License. See the LICENSE file for more details.
