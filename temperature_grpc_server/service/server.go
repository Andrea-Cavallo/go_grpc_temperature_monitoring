package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go_with_grpc/pkg/temperature"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Server represents the gRPC server for the TemperatureService, containing a map of temperature readings and a mutex for synchronization.
type Server struct {
	temperature.UnimplementedTemperatureServiceServer
	readings map[string][]*temperature.TemperatureReading
	mu       sync.Mutex
}

// NewServer creates a new instance of the TemperatureService server.
func NewServer() *Server {
	return &Server{
		readings: make(map[string][]*temperature.TemperatureReading),
	}
}

// Constants for WeatherAPI
const (
	envWeatherAPIKey = "WEATHER_API_KEY"
	weatherAPIURL    = "http://api.weatherapi.com/v1/current.json?key=%s&q=%s"
)

// GetCurrentTemperature retrieves the current temperature for a specified location.
// It takes a context and a GetCurentTemperatureRequest as input parameters.
// It returns a GetCurrentTemperatureResponse containing the temperature reading
// and an error if any error occurs during the process.
// The function logs the request and the current temperature value.
func (s *Server) GetCurrentTemperature(ctx context.Context, req *temperature.GetCurrentTemperatureRequest) (*temperature.GetCurrentTemperatureResponse, error) {
	location := req.Location
	log.Printf("Received request for current temperature of %s\n", location)

	temp, err := getWeatherData(location)
	if err != nil {
		log.Printf("Error retrieving weather data for %s: %v\n", location, err)
		return nil, err
	}

	reading := &temperature.TemperatureReading{
		Location:    location,
		Temperature: temp,
		Timestamp:   time.Now().Unix(),
	}

	log.Printf("Current temperature for %s: %.2f°C\n", location, temp)

	return &temperature.GetCurrentTemperatureResponse{
		Reading: reading,
	}, nil
}

// getWeatherData retrieves the current temperature for a specified location by making a request to the WeatherAPI.
// It uses the fetchWeatherAPIResponse function to construct and send the HTTP request,
// and then parses the response to extract the temperature value.
// The location argument specifies the location for which to retrieve the weather data.
// The function returns the current temperature as a float64 value and an error if any error occurs during the process.
// If the weather data is not found for the specified location, the function returns a temperature of 0 and an error.
// The function logs the temperature value and any errors encountered during the process.
func getWeatherData(location string) (float64, error) {
	body, err := fetchWeatherAPIResponse(location)
	if err != nil {
		return 0, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing response: %v", err)
		return 0, err
	}

	if current, ok := result["current"].(map[string]interface{}); ok {
		temp := current["temp_c"].(float64)
		log.Printf("Current temperature for %s is %.2f°C", location, temp)
		return temp, nil
	}

	log.Printf("Weather data not found for %s", location)
	return 0, fmt.Errorf("weather data not found")
}

// fetchWeatherAPIResponse sends an HTTP GET request to the WeatherAPI with the specified location
// and returns the response body as a byte slice ([]byte) and an error.
// It uses the provided location to construct the URL for the request and includes a API key for authentication.
// This function times out after 10 seconds using a context with a deadline.
// If an error occurs while creating or sending the request, it is logged and returned.
// The response body is read and returned as a byte slice. The response body is automatically closed after reading.
func fetchWeatherAPIResponse(location string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	weatherAPIKey := os.Getenv(envWeatherAPIKey)
	if weatherAPIKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY non impostata. Per favore registrati su https://www.weatherapi.com/ per ottenere una API key gratuita e imposta la variabile d'ambiente WEATHER_API_KEY")
	}
	url := fmt.Sprintf(weatherAPIURL, weatherAPIKey, location)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Received response from WeatherAPI with status: %s", resp.Status)

	return io.ReadAll(resp.Body)
}
