// Package portsclient implements Ports service API client.
package portsclient

import (
	"context"
	"fmt"

	"github.com/danielfurman/ports-microservices/internal/logs"
	"github.com/danielfurman/ports-microservices/internal/portssvc/portsgrpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GRPC is Ports service gRPC client.
type GRPC struct {
	client     portsgrpc.PortServiceClient
	connection *grpc.ClientConn
	log        *logrus.Entry
}

// NewGRPC creates new Ports service gRPC client for given server address.
// GRPC.Close() should be called when client is no longer needed.
func NewGRPC(serverAddress string) (GRPC, error) {
	log := logs.NewLogger("ports-client")

	log.WithField("server-address", serverAddress).Debug("Dialing gRPC")
	connection, err := grpc.Dial(serverAddress, dialOptions()...)
	if err != nil {
		return GRPC{}, fmt.Errorf("gRPC dial on %v: %w", serverAddress, err)
	}

	return GRPC{
		client:     portsgrpc.NewPortServiceClient(connection),
		connection: connection,
		log:        log,
	}, nil
}

func dialOptions() []grpc.DialOption {
	// TODO(dfurman): support TLS, timeout, logging, retries
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

// StorePort stores given port in Ports service.
func (g GRPC) StorePort(ctx context.Context, port *portsgrpc.Port) error {
	_, err := g.client.StorePort(ctx, &portsgrpc.StorePortRequest{
		Port: port,
	})
	return err
}

// ListPorts lists all ports stored in Ports service.
func (g GRPC) ListPorts(ctx context.Context) ([]*portsgrpc.Port, error) {
	response, err := g.client.ListPorts(ctx, &emptypb.Empty{})
	return response.GetPorts(), err
}

// Close closes the client connection.
func (g GRPC) Close() error {
	return g.connection.Close()
}
