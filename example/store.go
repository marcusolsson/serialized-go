package main

import (
	"flag"
	"log"
	"math/rand"

	serialized "github.com/marcusolsson/serialized-go"
	uuid "github.com/satori/go.uuid"
)

type PaymentProcessed struct {
	PaymentMethod string `json:"paymentMethod"`
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
}

func main() {
	var (
		accessKey       = flag.String("access-key", "", "Serialized.io access key")
		secretAccessKey = flag.String("secret-access-key", "", "Serialized.io secret access key")
	)

	flag.Parse()

	client := serialized.NewClient(
		serialized.WithAccessKey(*accessKey),
		serialized.WithSecretAccessKey(*secretAccessKey),
	)

	pp := PaymentProcessed{
		PaymentMethod: "CARD",
		Amount:        rand.Intn(1000),
		Currency:      "SEK",
	}

	err := client.Store("payment", "2c3cf88c-ee88-427e-818a-ab0267511c84", 0,
		serialized.NewEvent(uuid.NewV4().String(), "PaymentProcessed", pp),
		serialized.NewEvent(uuid.NewV4().String(), "PaymentProcessed", pp))
	if err != nil {
		log.Fatal(err)
	}
}
