# Cereal

A command-line interface for the Serialized.io API.

## Installation

```
go get -u github.com/marcusolsson/serialized-go/cmd/cereal
```

## Usage

```
$ cereal feed payment --since 67
TIMESTAMP                               AGGREGATE ID                            EVENT TYPE
Wed, 13 Sep 2017 14:31:19 +0200         2c3cf88c-ee88-427e-818a-ab0267511c84    PaymentProcessed
Wed, 13 Sep 2017 18:56:03 +0200         2c3cf88c-ee88-427e-818a-ab0267511c84    PaymentProcessed
```

For more information run `cereal help`.