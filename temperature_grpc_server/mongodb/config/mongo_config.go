package config

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance    *mongo.Client
	clientInstanceErr error
	mongoOnce         sync.Once
	cfg               Config
)

// Config contiene la configurazione per la connessione a MongoDB.
type Config struct {
	URI            string
	DatabaseName   string
	TimeoutSeconds int
}

// InitializeConfig inizializza la configurazione del database.
func InitializeConfig(uri string, dbName string, timeoutSeconds int) {
	cfg = Config{
		URI:            uri,
		DatabaseName:   dbName,
		TimeoutSeconds: timeoutSeconds,
	}
}

// GetMongoClient ritorna un'istanza singleton della connessione a MongoDB.
func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(cfg.URI)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.TimeoutSeconds)*time.Second)
		defer cancel()

		clientInstance, clientInstanceErr = mongo.Connect(ctx, clientOptions)
		if clientInstanceErr != nil {
			return
		}

		// Verifica la connessione
		clientInstanceErr = clientInstance.Ping(ctx, nil)
		if clientInstanceErr == nil {
			logrus.Println("Connected to MongoDB")
		}
	})

	return clientInstance, clientInstanceErr
}

// CloseMongoClient chiude la connessione a MongoDB.
func CloseMongoClient() error {
	if clientInstance == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return clientInstance.Disconnect(ctx)
}
