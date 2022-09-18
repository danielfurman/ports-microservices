package ingestsvc_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/danielfurman/ports-microservices/internal/ingestsvc"
	"github.com/danielfurman/ports-microservices/internal/portsclient"
	"github.com/danielfurman/ports-microservices/internal/portssvc"
	"github.com/danielfurman/ports-microservices/internal/portssvc/portsgrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestService_Run(t *testing.T) {
	tests := []struct {
		name          string
		filePath      string
		expectedError bool
		expectedPorts map[string]*portsgrpc.Port
	}{
		{
			name:          "foo",
			filePath:      filepath.Join("testdata", "3-ports.json"),
			expectedError: false,
			expectedPorts: map[string]*portsgrpc.Port{
				"AEAJM": {
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
				},
				"AEAUH": {
					Id:          "AEAUH",
					Name:        "Abu Dhabi",
					City:        "Abu Dhabi",
					Country:     "United Arab Emirates",
					Alias:       nil,
					Regions:     nil,
					Coordinates: []float64{54.37, 24.47},
					Province:    "Abu Z¸aby [Abu Dhabi]",
					Timezone:    "Asia/Dubai",
					Unlocs:      []string{"AEAUH"},
					Code:        "52001",
				},
				"AEDXB": {
					Id:          "AEDXB",
					Name:        "Dubai",
					City:        "Dubai",
					Country:     "United Arab Emirates",
					Alias:       nil,
					Regions:     nil,
					Coordinates: []float64{55.27, 25.25},
					Province:    "Dubayy [Dubai]",
					Timezone:    "Asia/Dubai",
					Unlocs:      []string{"AEDXB"},
					Code:        "52005",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			ctx, cancel := context.WithCancel(context.Background())
			server := portssvc.NewServer(portssvc.Config{GRPCServerAddress: ":0"})
			go func() {
				err := server.Serve(ctx)
				assert.NoError(t, err)
			}()
			defer cancel()

			client, err := portsclient.NewGRPC(server.Address().String())
			require.NoError(t, err)

			s, err := ingestsvc.NewService(ingestsvc.Config{
				PortsFilePath:       tt.filePath,
				PortsServiceAddress: server.Address().String(),
			})
			require.NoError(t, err)

			// When
			err = s.Run(ctx)

			// Then
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			ports, err := client.ListPorts(ctx)
			t.Logf("Listed ports: %v", ports)

			assert.NoError(t, err)
			// Ports are returned in arbitrary order
			assert.Equal(t, len(tt.expectedPorts), len(ports), "invalid number of listed ports")
			for _, p := range ports {
				expectedPort := tt.expectedPorts[p.Id]
				assertProtoEqual(t, expectedPort, p)
			}
		})
	}
}

func assertProtoEqual(t testing.TB, expected, actual proto.Message) {
	assert.True(
		t,
		proto.Equal(expected, actual),
		fmt.Sprintf("Protobuf messages are not equal:\nexpected: %v\nactual: %v", expected, actual),
	)
}
