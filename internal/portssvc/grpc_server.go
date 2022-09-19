// Package portssvc implement Ports domain service that allows to store and list Ports in persistence layer
// via gRPC API.
package portssvc

import (
	"context"
	"fmt"
	"net"

	"github.com/danielfurman/ports-microservices/internal/logs"
	"github.com/danielfurman/ports-microservices/internal/portssvc/adapter"
	"github.com/danielfurman/ports-microservices/internal/portssvc/domain/ports"
	"github.com/danielfurman/ports-microservices/internal/portssvc/portsgrpc"
	emptypb "github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCServer is a gRPC server for Ports domain service.
type GRPCServer struct {
	cfg Config

	portsgrpc.UnimplementedPortServiceServer
	service ports.Service
	log     *logrus.Entry

	listenerAddress net.Addr
	listenerReady   chan struct{}
}

// Config is a configuration for GRPCServer.
type Config struct {
	// GRPCServerAddress is a TCP address of the server. Default: ":9090".
	GRPCServerAddress string `env:"GRPC_SERVER_ADDRESS" envDefault:":9090"`
}

// NewServer creates new GRPCServer with given configuration.
func NewServer(cfg Config) *GRPCServer {
	logs.Configure()
	log := logs.NewLogger("ports-server")
	log.WithField("config", fmt.Sprintf("%+v", cfg)).Debug("Creating ports server")

	return &GRPCServer{
		cfg: cfg,
		// TODO(dfurman): implement and use adapter.NewPostgresRepository()
		service:       ports.NewService(adapter.NewInMemoryPortsRepository()),
		log:           log,
		listenerReady: make(chan struct{}),
	}
}

// Serve starts the gRPC Ports server. The server is gracefully stopped on context cancel/timeout.
// TODO(dfurman): support request/response logging.
func (s *GRPCServer) Serve(ctx context.Context) error {
	grpcServer := grpc.NewServer()
	portsgrpc.RegisterPortServiceServer(grpcServer, s)

	go s.gracefulStopOnCancel(ctx, grpcServer)

	listener, err := net.Listen("tcp", s.cfg.GRPCServerAddress)
	if err != nil {
		return fmt.Errorf("listen TCP: %w", err)
	}

	s.listenerAddress = listener.Addr()
	s.markListenerReady()

	s.log.WithField("address", s.listenerAddress).Info("Starting gRPC Ports server")
	err = grpcServer.Serve(listener)
	s.log.Info("gRPC Ports server stopped")

	return err
}

// gracefulStopOnCancel stops the server on context cancel/timeout.
func (s *GRPCServer) gracefulStopOnCancel(ctx context.Context, grpcServer *grpc.Server) {
	<-ctx.Done()
	s.log.Debug("Stopping the gRPC server")
	grpcServer.GracefulStop()
}

// Address return the TCP address of the server.
func (s *GRPCServer) Address() net.Addr {
	<-s.listenerReady
	return s.listenerAddress
}

func (s *GRPCServer) markListenerReady() {
	close(s.listenerReady)
}

// StorePort handles the store port request.
func (s *GRPCServer) StorePort(ctx context.Context, req *portsgrpc.StorePortRequest) (*emptypb.Empty, error) {
	err := s.service.StorePort(
		ctx,
		portPayloadToDomain(req.GetPort()),
	)
	return &emptypb.Empty{}, err
}

// ListPorts handles the list ports request.
func (s *GRPCServer) ListPorts(ctx context.Context, _ *emptypb.Empty) (*portsgrpc.ListPortsResponse, error) {
	p, err := s.service.ListPorts(ctx)
	return &portsgrpc.ListPortsResponse{
		Ports: domainPortsToPayload(p),
	}, err
}

func portPayloadToDomain(p *portsgrpc.Port) *ports.Port {
	if p == nil {
		return nil
	}

	return &ports.Port{
		ID:          p.Id,
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

func domainPortsToPayload(ports []ports.Port) []*portsgrpc.Port {
	result := make([]*portsgrpc.Port, 0, len(ports))
	for i := range ports {
		result = append(result, domainPortToPayload(ports[i]))
	}
	return result
}

func domainPortToPayload(p ports.Port) *portsgrpc.Port {
	return &portsgrpc.Port{
		Id:          p.ID,
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
