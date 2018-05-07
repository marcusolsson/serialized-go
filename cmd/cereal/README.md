# Cereal

A command-line interface for the Serialized.io API.

## Installation

```
go get -u github.com/marcusolsson/serialized-go/cmd/cereal
```

## Usage

```
export SERIALIZED_ACCESS_KEY=<accessKey>
export SERIALIZED_SECRET_ACCESS_KEY=<secretAccessKey>
cereal feeds
```

## Examples

### Show aggregate information

```
$ cereal aggregate 2c3cf88c-ee88-427e-818a-ab0267511c84 --type payment
Type:      payment
ID:        2c3cf88c-ee88-427e-818a-ab0267511c84
Version:   2

Showing the 10 most recent events:

EVENT ID                               TYPE                 DATA
3ba52c7a-8129-444b-94da-a3dc2549845f   PaymentProcessed   {"paymentMethod":"CARD","amount":1000,"currency":"SEK"}
437e0856-713e-4a28-9d94-0c9489962d39   PaymentProcessed   {"paymentMethod":"CARD","amount":99,"currency":"SEK"}
```

### Show single projection

```
$ cereal projection orderTotal --agg-id 2c3cf88c-ee88-427e-818a-ab0267511c84
{
  "total": 189.0
}
```

### Follow the feed

```
$ cereal feed payment --since 67
TIMESTAMP                               AGGREGATE ID                            EVENT TYPE
Wed, 13 Sep 2017 14:31:19 +0200         2c3cf88c-ee88-427e-818a-ab0267511c84    PaymentProcessed
Wed, 13 Sep 2017 18:56:03 +0200         2c3cf88c-ee88-427e-818a-ab0267511c84    PaymentProcessed
```

For more information run `cereal help`.