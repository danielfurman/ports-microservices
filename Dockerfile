# Services builder image
FROM golang:1.19 as builder

ENV CGO_ENABLED 0

COPY cmd /src/cmd
COPY internal /src/internal
COPY go.mod /src/go.mod
COPY go.sum /src/go.sum

WORKDIR /src
RUN go build -o /build/portssvc ./cmd/portssvc
RUN go build -o /build/ingestsvc ./cmd/ingestsvc

# Ports service image
FROM alpine:3.16 as portssvc

# Set up non-root user
RUN addgroup -g 1000 -S service && \
    adduser -u 1000 -h /app -G service -S service
USER service

COPY --from=builder --chown=service:service /build/portssvc /app/

WORKDIR /app
CMD ["./portssvc"]

# Ingest service image
FROM alpine:3.16 as ingestsvc

# Set up non-root user
RUN addgroup -g 1000 -S service && \
    adduser -u 1000 -h /app -G service -S service
USER service

COPY --from=builder --chown=service:service /build/ingestsvc /app/

WORKDIR /app
CMD ["./ingestsvc"]
