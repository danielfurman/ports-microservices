## Ports microservices

[![Go Reference](https://pkg.go.dev/badge/github.com/danielfurman/ports-microservices.svg)](https://pkg.go.dev/github.com/danielfurman/ports-microservices)

This project contains two simple microservices developed as a recruitment coding assignment.

## Requirements

- [Go](https://golang.org/doc/install) >= Go 1.18
- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](https://docs.docker.com/engine/install)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Overview

This repository contains two microservices:
1. Ports service that exposes a gRPC API that allows to store and list Ports in persistence layer.
2. Ingest service that allows to read Port resources from input JSON file and store them in Ports service via gRPC.

Ingest service reads resources from JSON file one-by-one using a stream, so it does not load all data to its memory and supports large files.
In current implementation the Ingest service writes resources to Ports service sequentially in order not to overload it.

Services are configured with environment variables:
- Ports service config: [portssvc/grpc_server.go -> Config struct](./internal/portssvc/grpc_server.go)
- Ingest service config: [ingestsvc/ingest_service.go -> Config struct](./internal/ingestsvc/ingest_service.go)

## Usage

Both services are containerized and can be run with Docker Compose:
1. Run services in Docker containers: `docker-compose up`.
2. Stop containers by press Ctrl+C.
3. Clean up containers: `docker-compose down`.

### Description

These are steps that should happen on a successful run:
1. Docker compose spawns Ports service and Ingest service.
2. Ingest service connects to Ports service via gRPC.
3. Ingest service reads Port resources from [ports.json](./internal/ingestsvc/testdata/ports.json) file and transmits them to Ports service one-by-one. 
4. Each Ports resources received by Port service is stored in its in-memory storage.
5. After processing all resources, Ingest service shuts down.
6. Ports service continues to operate. It can be gracefully shut down by sending SIGTERM or SIGINT signals.

Note that the [ports.json](./internal/ingestsvc/testdata/ports.json) file used by Ingest service can be modified in [docker-compose.yml](./docker-compose.yml) by changing the volume path.

Stripped logs of a successful run are presented below:

```shell
➜  ports-microservices git:(master) docker-compose up
[+] Running 3/3
 ⠿ Network ports-microservices_default        Created   0.0s
 ⠿ Container ports-microservices-portssvc-1   Created   0.1s
 ⠿ Container ports-microservices-ingestsvc-1  Created   0.1s
Attaching to ports-microservices-ingestsvc-1, ports-microservices-portssvc-1
ports-microservices-portssvc-1   | time="2022-09-19T13:19:57Z" level=debug msg="Creating ports server" config="{GRPCServerAddress::9090}" logger=ports-server
ports-microservices-portssvc-1   | time="2022-09-19T13:19:57Z" level=info msg="Starting gRPC Ports server" address="[::]:9090" logger=ports-server
ports-microservices-ingestsvc-1  | time="2022-09-19T13:19:57Z" level=debug msg="Creating ingest service" config="{PortsFilePath:/app/ports.json PortsServiceAddress:portssvc:9090}" logger=ingest-service
ports-microservices-ingestsvc-1  | time="2022-09-19T13:19:57Z" level=debug msg="Dialing gRPC" logger=ports-client server-address="portssvc:9090"
ports-microservices-ingestsvc-1  | time="2022-09-19T13:19:57Z" level=debug msg="Storing port in ports service" logger=ingest-service port-id=AEAJM
ports-microservices-portssvc-1   | time="2022-09-19T13:19:57Z" level=debug msg="Storing port" logger=ports-service port-id=AEAJM
[...]
ports-microservices-portssvc-1   | time="2022-09-19T13:19:58Z" level=debug msg="Storing port" logger=in-memory-ports-repo port="&{ZWUTA Mutare Mutare Zimbabwe [] [] [32.650351 -18.9757714] Manicaland Africa/Harare [ZWUTA] }"
ports-microservices-ingestsvc-1  | time="2022-09-19T13:19:58Z" level=info msg="Ingest service finished successfully"
ports-microservices-ingestsvc-1 exited with code 0
```

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
