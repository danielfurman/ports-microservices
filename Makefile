.DEFAULT_GOAL := build

.PHONY: build
build:
	go build ./...

.PHONY: dev-check
dev-check: build docker-build test lint

.PHONY: docker-build
docker-build:
	# TODO(dfurman): implement

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: generate
generate:
	protoc --proto_path=api/grpc --go_out=internal/portssvc/portsgrpc --go_opt=paths=source_relative \
		--go-grpc_out=internal/portssvc/portsgrpc --go-grpc_opt=paths=source_relative  api/grpc/ports.proto

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test ./...

.PHONY: test-race
test-race:
	go test -race ./...
