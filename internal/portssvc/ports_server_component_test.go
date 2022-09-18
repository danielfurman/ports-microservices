package portssvc_test

import (
	"context"
	"testing"

	"github.com/danielfurman/ports-microservices/internal/portsclient"
	"github.com/danielfurman/ports-microservices/internal/portssvc"
	"github.com/danielfurman/ports-microservices/internal/portssvc/portsgrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPortsServer_StorePorts(t *testing.T) {
	for _, tt := range []struct {
		name      string
		inputPort *portsgrpc.Port
		// TODO(dfurman): verify gRPC error codes
		expectedError bool
		expectedPorts []*portsgrpc.Port
	}{
		{
			name:          "nil port given",
			expectedError: true,
		}, {
			name:          "empty port given",
			expectedError: true,
		}, {
			name:      "valid port given",
			inputPort: newAjmanPort(),
			expectedPorts: []*portsgrpc.Port{
				newAjmanPort(),
			},
		}, {
			name: "port with empty ID given",
			inputPort: func() *portsgrpc.Port {
				p := newAjmanPort()
				p.Id = ""
				return p
			}(),
			expectedError: true,
		}, {
			name: "port with empty name given",
			inputPort: func() *portsgrpc.Port {
				p := newAjmanPort()
				p.Name = ""
				return p
			}(),
			expectedError: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			ctx, cancel := context.WithCancel(context.Background())
			server := portssvc.NewServer(portssvc.Config{
				GRPCServerAddress: ":0",
			})
			go func() {
				err := server.Serve(ctx)
				assert.NoError(t, err)
			}()
			defer cancel()

			client, err := portsclient.NewGRPC(server.Address().String())
			require.NoError(t, err)

			// When
			err = client.StorePort(ctx, tt.inputPort)

			// Then
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			ports, err := client.ListPorts(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPorts, ports)
		})
	}
}

func newAjmanPort() *portsgrpc.Port {
	return &portsgrpc.Port{
		Id:          "AEAJM",
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
