// Package serialized provides a client for interacting with https://serialized.io.
//
// Getting started
//
// First, we create a new Client with the Serialized.io credentials:
//
//   var (
//   	accessKey       = os.Getenv("SERIALIZED_ACCESS_KEY")
//   	secretAccessKey = os.Getenv("SERIALIZED_SECRET_ACCESS_KEY")
//   )
//
//   client := serialized.NewClient(
//   	serialized.WithAccessKey(accessKey),
//   	serialized.WithSecretAccessKey(secretAccessKey),
//   )
//
// Storing events
//
// To store a new event, call the Store method on the Client you just created:
//
//   var (
//   	aggregateType = "order"
//   	aggregateID   = uuid.New().String()
//   	eventID       = uuid.New().String()
//   	version       = 0
//   )
//
//   event := &serialized.Event{
//   	ID:   eventID,
//   	Type: "PaymentReceived",
//   	Data: []byte(`{"amount": 100}`),
//   }
//
//   if err := client.Store(ctx, aggregateType, aggregateID, version, event); err != nil {
//   	log.Fatal(err)
//   }
//
// Loading an aggregate
//
// Once you've started storing your events, at some point you'll want to load
// your aggregate. Let's load the aggregate you just created:
//
//  agg, err := client.LoadAggregate(ctx, aggregateType, aggregateID)
//  if err != nil {
//  	log.Fatal(err)
//  }
//
// Next we can iterate over the events to build the aggregate state.
//
//  for _, event := range agg.Events {
//  	switch event.Type {
//  	case "PaymentReceived":
//  		var e PaymentReceived
//  		json.Unmarshal(event.Data, &e)
//
//  		log.Printf("Received payment of %d", e.Amount)
//  	}
//  }
package serialized
