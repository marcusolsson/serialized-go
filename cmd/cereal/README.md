# serialized-cli

A command-line interface for the Serialized.io API.

## Installation

```
go get -u github.com/marcusolsson/serialized-go/cmd/serialized-cli
```

## Usage

```
$ serialized-cli feed payment --since 67
SEQUENCE AGGREGATE ID                         NUM EVENTS TIMESTAMP
68       2c3cf88c-ee88-427e-818a-ab0267511c84 1          2017-09-03T11:01:33+02:00
69       2c3cf88c-ee88-427e-818a-ab0267511c84 1          2017-09-03T11:01:34+02:00
70       2c3cf88c-ee88-427e-818a-ab0267511c84 1          2017-09-03T11:01:35+02:00
```

For more information run `serialized-cli help`.