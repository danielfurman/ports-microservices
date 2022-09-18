package ports

import (
	"context"
	"fmt"

	"github.com/danielfurman/ports-microservices/internal/logs"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	StorePort(context.Context, *Port) error
	ListPorts(context.Context) ([]Port, error)
}

type Service struct {
	portsRepo Repository
	log       *logrus.Entry
}

func NewService(pr Repository) Service {
	return Service{
		portsRepo: pr,
		log:       logs.NewLogger("ports-service"),
	}
}

func (s Service) StorePort(ctx context.Context, port *Port) error {
	if port == nil {
		return fmt.Errorf("nil port given")
	}
	s.log.WithField("port-id", port.ID).Debug("Storing port")

	err := port.Validate()
	if err != nil {
		return fmt.Errorf("validate port: %w", err)
	}

	return s.portsRepo.StorePort(ctx, port)
}

func (s Service) ListPorts(ctx context.Context) ([]Port, error) {
	s.log.Debug("Listing ports")
	ports, err := s.portsRepo.ListPorts(ctx)
	return ports, err
}
