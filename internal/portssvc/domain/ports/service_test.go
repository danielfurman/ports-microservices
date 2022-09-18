package ports_test

import (
	"context"
	"testing"

	"github.com/danielfurman/ports-microservices/internal/portssvc/adapter"
	"github.com/danielfurman/ports-microservices/internal/portssvc/domain/ports"
	"github.com/stretchr/testify/assert"
)

func TestService_StorePort(t *testing.T) {
	for _, tt := range []struct {
		name          string
		inputPort     *ports.Port
		expectedError bool
		expectedPorts []ports.Port
	}{
		{
			name:          "nil port given",
			expectedError: true,
			expectedPorts: []ports.Port{},
		}, {
			name:          "empty port given",
			expectedError: true,
			expectedPorts: []ports.Port{},
		}, {
			name:      "valid port given",
			inputPort: newAjmanPort(),
			expectedPorts: []ports.Port{
				*newAjmanPort(),
			},
		}, {
			name: "port with empty ID given",
			inputPort: func() *ports.Port {
				p := newAjmanPort()
				p.ID = ""
				return p
			}(),
			expectedError: true,
			expectedPorts: []ports.Port{},
		}, {
			name: "port with empty name given",
			inputPort: func() *ports.Port {
				p := newAjmanPort()
				p.Name = ""
				return p
			}(),
			expectedError: true,
			expectedPorts: []ports.Port{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			ctx := context.Background()
			service := ports.NewService(adapter.NewInMemoryPortsRepository())

			// When
			err := service.StorePort(context.Background(), tt.inputPort)

			// Then
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			ps, err := service.ListPorts(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPorts, ps)
		})
	}
}

func newAjmanPort() *ports.Port {
	return &ports.Port{
		ID:          "AEAJM",
		Name:        "Ajman",
		City:        "Ajman",
		Country:     "United Arab Emirates",
		Alias:       []string{"foo-alias", "bar-alias"},
		Regions:     []string{"foo-region", "bar-region"},
		Coordinates: []float64{55.5136433, 25.4052165},
		Province:    "Ajman",
		Timezone:    "Asia/Dubai",
		Unlocs:      []string{"AEAJM"},
		Code:        "52000",
	}
}
