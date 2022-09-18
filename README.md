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

TODO

## Development

Optional development tools:
- [Golangci-lint](https://golangci-lint.run/usage/install/#local-installation) for static code analysis
- [Protoc v3](https://grpc.io/docs/protoc-installation/), [Protoc-gen-go, Protoc-gen-go-grpc](https://grpc.io/docs/languages/go/quickstart/) for gRPC code generation

Development commands:
- Build code and Docker image, test and lint code: `make dev-check`
- Build code: `make build`
- Build Docker image of the Ports service: `make docker-build`
- Run tests: `make test`
- Run tests with race detector: `make test-race`
- Run static code analysis: `make lint`
- Format source code: `make fmt`
- Regenerate the source code: `make generate`

TODO:
- Dockerize the service
- Add docstrings and documentation
