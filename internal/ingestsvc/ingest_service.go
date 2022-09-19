// Package ingestsvc contains source code for Ingest service that allows to read port resources from input JSON file
// and store them in Ports service.
package ingestsvc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/danielfurman/ports-microservices/internal/logs"
	"github.com/danielfurman/ports-microservices/internal/portsclient"
	"github.com/danielfurman/ports-microservices/internal/portssvc/portsgrpc"
	"github.com/sirupsen/logrus"
)

// Service is an Ingest service.
type Service struct {
	cfg Config

	portsClient portsclient.GRPC
	log         *logrus.Entry
}

// Config is a config for Ingest service.
type Config struct {
	// PortsFilePath is a path to the JSON file containing input ports. Env var: PORTS_FILE_PATH. Required.
	PortsFilePath string `env:"PORTS_FILE_PATH,notEmpty"`
	// PortsServiceAddress is a TCP address of the Ports service. Env var: PORTS_SVC_ADDRESS. Default: ":9090".
	PortsServiceAddress string `env:"PORTS_SVC_ADDRESS" envDefault:":9090"`
}

// NewService creates new Ingest service with given configuration.
func NewService(cfg Config) (Service, error) {
	logs.Configure()
	log := logs.NewLogger("ingest-service")
	log.WithField("config", fmt.Sprintf("%+v", cfg)).Debug("Creating ingest service")

	client, err := portsclient.NewGRPC(cfg.PortsServiceAddress)
	if err != nil {
		return Service{}, fmt.Errorf("new ports gRPC client: %w", err)
	}

	return Service{
		cfg:         cfg,
		portsClient: client,
		log:         log,
	}, nil
}

// Run reads all port resources from specified in input JSON file and transmits them to Ports service via gRPC.
//
// The example JSON file with a format expected by the service is located in ./testdata/ports.json.
// Resources are read from the file one-by-one with a stream to reduce memory consumption and support large files.
// Run can be stopped by context cancel/timeout.
// This function is meant to be called only once, because it closes Ports client connection.
func (s Service) Run(ctx context.Context) (err error) {
	defer func() {
		cErr := s.portsClient.Close()
		if cErr != nil && err == nil {
			err = fmt.Errorf("failed to close client connection: %w", cErr)
		}
	}()

	file, err := os.Open(s.cfg.PortsFilePath)
	if err != nil {
		return fmt.Errorf("open ports file: %w", err)
	}
	defer func() {
		cErr := file.Close()
		if cErr != nil && err == nil {
			err = fmt.Errorf("failed to close ports file: %w", cErr)
		}
	}()

	return s.decodeAndIngestPorts(ctx, file)
}

func (s Service) decodeAndIngestPorts(ctx context.Context, reader io.Reader) error {
	decoder := json.NewDecoder(reader)

	err := s.readOpeningBracket(decoder)
	if err != nil {
		return err
	}

	for decoder.More() {
		err = s.decodeAndIngestPort(ctx, decoder)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Service) readOpeningBracket(decoder *json.Decoder) error {
	openingT, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("decoder.Token(): %w", err)
	}

	if openingT != json.Delim('{') {
		return fmt.Errorf("opening JSON token is %v, expected '{'", openingT)
	}
	return nil
}

func (s Service) decodeAndIngestPort(ctx context.Context, decoder *json.Decoder) error {
	portKeyT, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("decoder.Token(): %w", err)
	}

	portKey, ok := portKeyT.(string)
	if !ok {
		return fmt.Errorf("port key is expected to be string, got %+#v", portKey)
	}

	var port Port
	err = decoder.Decode(&port)
	if err != nil {
		return fmt.Errorf("decode port object: %w", err)
	}

	s.log.WithField("port-id", portKey).Debug("Storing port in ports service")
	err = s.portsClient.StorePort(ctx, portToPayload(port, portKey))
	if err != nil {
		return fmt.Errorf("store port with ID %v in ports service: %w", portKey, err)
	}
	return nil
}

// Port models a JSON representation of the port.
type Port struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

func portToPayload(p Port, portID string) *portsgrpc.Port {
	return &portsgrpc.Port{
		Id:          portID,
		Name:        p.Name,
		City:        p.City,
		Country:     p.Country,
		Alias:       p.Alias,
		Regions:     p.Regions,
		Coordinates: p.Coordinates,
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		Code:        p.Code,
	}
}
