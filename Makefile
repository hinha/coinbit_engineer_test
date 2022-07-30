proto-gen:
	protoc --go_out=. ./proto/*.proto

run-dev-api:
	go run cmd/main.go -service=api

run-dev-broker-balance:
	go run cmd/main.go -service=broker -broker_name=balance-copy

run-dev-broker-flagger:
	go run cmd/main.go -service=broker -broker_name=flagger-copy

run-dev-broker-detector:
	go run cmd/main.go -service=broker -broker_name=threshold-copy