package temperature_repo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
)

// TemperatureReading rappresenta una lettura di temperatura.
type TemperatureReading struct {
	Timestamp time.Time `bson:"timestamp"`
	Value     float64   `bson:"value"`
}

// InsertTemperatureReading inserisce un nuovo dato di temperatura e traccia l'operazione.
func InsertTemperature(client *mongo.Client, reading TemperatureReading) error {
	tracer := otel.Tracer("mongodb")
	ctx, span := tracer.Start(context.Background(), "Insert-Temperature")
	defer span.End()

	collection := client.Database("temperatureDB").Collection("temperatures")
	_, err := collection.InsertOne(ctx, reading)
	if err != nil {
		span.RecordError(err)
	}
	return err
}
