package adapter

import (
	"context"
	"fmt"
	"sync"

	"github.com/danielfurman/ports-microservices/internal/logs"
	"github.com/danielfurman/ports-microservices/internal/portssvc/domain/ports"
	"github.com/sirupsen/logrus"
)

type InMemoryPortsRepository struct {
	ports      map[string]*ports.Port
	portsMutex sync.RWMutex
	log        *logrus.Entry
}

func NewInMemoryPortsRepository() *InMemoryPortsRepository {
	return &InMemoryPortsRepository{
		ports: make(map[string]*ports.Port),
		log:   logs.NewLogger("in-memory-ports-repo"),
	}
}

func (r *InMemoryPortsRepository) StorePort(_ context.Context, port *ports.Port) error {
	r.log.WithField("port", port).Debug("Storing port")
	r.portsMutex.Lock()
	defer r.portsMutex.Unlock()

	r.ports[port.ID] = port
	return nil
}

func (r *InMemoryPortsRepository) ListPorts(_ context.Context) ([]ports.Port, error) {
	r.log.Debug("Listing ports")
	r.portsMutex.RLock()
	defer r.portsMutex.RUnlock()

	ps, err := portsToSlice(r.ports)
	return ps, err
}

func portsToSlice(portsM map[string]*ports.Port) ([]ports.Port, error) {
	result := make([]ports.Port, 0, len(portsM))
	for k, p := range portsM {
		if p == nil {
			return nil, fmt.Errorf("nil port in repository on key %v", k)
		}
		result = append(result, *p)
	}
	return result, nil
}
