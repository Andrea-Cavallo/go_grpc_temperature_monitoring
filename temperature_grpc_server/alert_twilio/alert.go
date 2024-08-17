package alert_twilio

import (
	"fmt"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"log"
	"os"
	"strconv"
)

// SendTemperatureAlert sends a temperature alert to a specified phone number using the Twilio API.
// It takes the temperature value as a float64 input parameter.
// The function retrieves the necessary environment variables for the Twilio account,
// creates a Twilio REST client with the specified credentials,
// and sends a message containing the temperature alert to the specified recipient.
// It returns an error if the alert fails to send, or if any necessary environment variables are missing.
// The function logs the successful sending of the alert if no errors occur.
func sendTemperatureAlert(temperature float64) error {
	//see: https://www.twilio.com/docs/whatsapp/quickstart/go
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromPhone := os.Getenv("TWILIO_PHONE_NUMBER")
	toPhone := os.Getenv("ALERT_PHONE_NUMBER")

	if accountSid == "" || authToken == "" || fromPhone == "" || toPhone == "" {
		return fmt.Errorf("missing Twilio environment variables")
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	msg := fmt.Sprintf("Alert! The temperature has exceeded the threshold: %.2f°C", temperature)

	params := &openapi.CreateMessageParams{}
	params.SetTo(toPhone)
	params.SetFrom(fromPhone)
	params.SetBody(msg)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send alert: %v", err)
	}

	log.Printf("Alert sent successfully to %s", toPhone)
	return nil
}

func CheckAndSendAlert(temp float64) {
	const defaultTemperatureThreshold = 35.0

	temperatureThreshold, err := strconv.ParseFloat(os.Getenv("TEMPERATURE_THRESHOLD"), 64)
	if err != nil {
		temperatureThreshold = defaultTemperatureThreshold
	}

	if temp > temperatureThreshold {
		log.Printf("Temperature alert triggered: %.2f°C", temp)
		if err := sendTemperatureAlert(temp); err != nil {
			log.Printf("Failed to send temperature alert: %v", err)
		}
	}
}
