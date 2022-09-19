.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o build/ ./...

.PHONY: dev-check
dev-check: build docker-build test lint

.PHONY: docker-build
docker-build:
	docker-compose build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: generate
generate:
	protoc --proto_path=api/grpc --go_out=internal/portssvc/portsgrpc --go_opt=paths=source_relative \
		--go-grpc_out=internal/portssvc/portsgrpc --go-grpc_opt=paths=source_relative  api/grpc/ports.proto

.PHONY: lint
lint:
	# govulncheck in golangci-lint ticket: https://github.com/golangci/golangci-lint/issues/3094
	govulncheck ./...
	golangci-lint run ./...

.PHONY: test
test:
	go test ./...

.PHONY: test-race
test-race:
	go test -race ./...
