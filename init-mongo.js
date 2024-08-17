// init-mongo.js

db = db.getSiblingDB('temperatureDB');  // Seleziona il database temperatureDB

db.createCollection('temperatures');  // Crea la collezione temperatures

// Inserisce un documento iniziale per creare il database e la collezione
db.readings.insert({
    "Timestamp": new Date(),
    "Value": 0.0
});
