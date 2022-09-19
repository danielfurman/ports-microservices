## Ports microservices

This project is a simple microservice developed as a recruitment coding assignment.

## Requirements

- [Go](https://golang.org/doc/install) >= Go 1.18
- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](https://docs.docker.com/engine/install)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Overview

TODO

## Usage

Both services are containerized and can be run with docker compose. Usage:
1. Run services: `docker-compose up`.
2. Stop services by press Ctrl+C.
3. Clean up containers: `docker-compose down`.

## Development

Development tools:
- [Protoc v3](https://grpc.io/docs/protoc-installation/), [Protoc-gen-go, Protoc-gen-go-grpc](https://grpc.io/docs/languages/go/quickstart/) for gRPC code generation
- [Golangci-lint](https://golangci-lint.run/usage/install/#local-installation) for static code analysis
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) for vulnerability analysis: `go install golang.org/x/vuln/cmd/govulncheck@latest`

Development commands:
- Build binaries and Docker image, test and lint code: `make dev-check`
- Build binaries to build directory: `make build`
- Build Docker image of the Ports service: `make docker-build`
- Run tests: `make test`
- Run tests with race detector: `make test-race`
- Run static code analysis: `make lint`
- Format source code: `make fmt`
- Regenerate the source code: `make generate`
