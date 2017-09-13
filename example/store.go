package main

import (
	"encoding/json"
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

	ev := &serialized.Event{
		ID:   uuid.NewV4().String(),
		Type: "PaymentProcessed",
		Data: mustMarshal(PaymentProcessed{
			PaymentMethod: "CARD",
			Amount:        rand.Intn(1000),
			Currency:      "SEK",
		}),
	}

	err := client.Store("payment", "2c3cf88c-ee88-427e-818a-ab0267511c84", 0, ev)
	if err != nil {
		log.Fatal(err)
	}
}

func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
