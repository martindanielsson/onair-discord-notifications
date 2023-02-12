package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/martindanielsson/onair-discord-notifications/onair/va"
)

func main() {
	vclient, err := va.New(nil, goDotEnvVariable("API_KEY"))

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fetchNotifications(vclient)
	fetchFlights(vclient)
	fetchCashflow(vclient)

}

func fetchNotifications(vclient *va.Client) {
	notifications, _, err := vclient.Notifications(goDotEnvVariable("VA_ID"))

	if err != nil {
		log.Fatalf("error fetching notifications: %v", err)
	}

	for _, notification := range notifications {
		fmt.Println(notification)
	}
}

func fetchFlights(vclient *va.Client) {
	flights, _, err := vclient.Flights(goDotEnvVariable("VA_ID"))

	if err != nil {
		log.Fatalf("error fetching flights: %v", err)
	}

	for _, flight := range flights {
		fmt.Println(flight)
	}
}

func fetchCashflow(vclient *va.Client) {
	cashflow, _, err := vclient.CashFlow(goDotEnvVariable("VA_ID"))

	if err != nil {
		log.Fatalf("error fetching cashflow: %v", err)
	}

	for _, entry := range cashflow.Entries {
		fmt.Println(entry)
	}
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("unable to read .env file: %s", err)
	}

	return os.Getenv(key)
}
