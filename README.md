# coinbit_engineer_test

This project was created as a stage of completing the submission test. the main goal of creating a service using an implementation of event-driven architecture using Kafka as a message broker

## Getting Started

Here some preparations that you might need before using this project.

### Application Config
Before start the appliaction, first create configuration file. Name configuration is `config.env`
```env
PORT=8000
BROKER="localhost:9092"
```

### Topics

We define brokers includes topic-group it represents implement interface. When the Execute() function is called it will run the consumer process of kafka. See code
`provider.go`

```golang
type Group string

const (
	BalanceGroup   Group = "balance"
	FlaggerGroup   Group = "flagger"
	ThresholdGroup Group = "threshold"
)

type Broker interface {
	Execute(ctx context.Context, groupName Group) error
}
```

### Protobuf

we collect the need for transform protobuf as encoding/decoding payload for kafka 
broker. See code `stream_encoders.go`

```golang
var (
	DepositStream    goka.Stream = "deposit"
	FlagWalletStream goka.Stream = "flag_wallet"
)

type DepositEncoder struct{}

func (c *DepositEncoder) Encode(value interface{}) ([]byte, error) {
	...
}
func (c *DepositEncoder) Decode(data []byte) (interface{}, error) {
	...
}

type CounterEncoder struct{}

func (c *CounterEncoder) Encode(value interface{}) ([]byte, error) {
	...
}

func (c *CounterEncoder) Decode(data []byte) (interface{}, error) {
	...
}
```

## API Endpoints

> Goka provides three components to build systems: emitters, processors, and views. The following figure depicts our initial design using these three components together with Kafka and the endpoints.

### Deposit endpoint
If John wants to deposit money, he would send a request to the send endpoint with the wallet_id and the amount of the message. For example:

```bash
curl -X POST \
    -d '{"wallet_id": "69da5a28bfee68a8b0bfeba076ee43f8ca4", "amount": 3500}' \
    http://localhost:8080/api/deposit
{
    "code": 200,
    "data": null,
    "message": "success"
}
```

### Check Endpoint

When John wants to check his wallet balance, he requests that from the check endpoint. For example:

```bash
curl --location --request GET 'http://localhost:8000/api/check/69da5a28bfee68a8b0bfeba076ee43f8ca4'
{
	"wallet_id": "69da5a28bfee68a8b0bfeba076ee43f8ca4",
	"balance": 3500,
	"above_threshold": true
}
```

## Running the project

In this project, we can put the endpoint handlers and, therefore, emitter and view in the same Go program. we start the collector processor. In another Go program can start the api service

```makefile
proto-gen:
	protoc --go_out=. ./proto/*.proto

run-dev-api:
	go run cmd/main.go -service=api

run-dev-broker-balance:
	go run cmd/main.go -service=broker -broker_name=balance

run-dev-broker-flagger:
	go run cmd/main.go -service=broker -broker_name=flagger

run-dev-broker-detector:
	go run cmd/main.go -service=broker -broker_name=threshold
```