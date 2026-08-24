package config

import (
	"testing"
	"time"
)

func TestServerValidate(t *testing.T) {
	tests := []struct {
		name    string
		server  server
		wantErr bool
	}{
		{
			name: "valid",
			server: server{
				Host:            "localhost",
				Port:            8080,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty host",
			server: server{
				Port:            8080,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "port too low",
			server: server{
				Host:            "localhost",
				Port:            0,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "port too high",
			server: server{
				Host:            "localhost",
				Port:            65536,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid shutdown timeout",
			server: server{
				Host:            "localhost",
				Port:            8080,
				ShutdownTimeout: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.validate()

			if (err != nil) != tt.wantErr {
				t.Fatalf("validateServer() error = %v", err)
			}
		})
	}
}
