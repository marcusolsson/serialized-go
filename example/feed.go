package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	serialized "github.com/marcusolsson/serialized-go"
)

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

	feed, err := client.Feed(context.Background(), "payment", 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range feed.Entries {
		fmt.Printf("%+v\n", e)
	}

}
