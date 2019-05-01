package main

type OrderEvent interface{}

type OrderPlacedEvent struct {
	CustomerID CustomerID `json:"customerId"`
	Amount     Amount     `json:"amount"`
}

type PaymentReceivedEvent struct {
	CustomerID CustomerID `json:"customerId"`
	AmountPaid Amount     `json:"amountPaid"`
}

type OrderFullyPaidEvent struct {
	CustomerID CustomerID `json:"customerId"`
}

type OrderCancelledEvent struct {
	CustomerID CustomerID `json:"customerId"`
	Reason     string     `json:"reason"`
}

type OrderShippedEvent struct {
	CustomerID     CustomerID     `json:"customerId"`
	TrackingNumber TrackingNumber `json:"trackingNumber"`
}
