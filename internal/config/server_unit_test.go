package config

import (
	"testing"
	"time"
)

func TestServerValidate(t *testing.T) {
	tests := []struct {
		name    string
		server  Server
		wantErr bool
	}{
		{
			name: "valid",
			server: Server{
				Host:            "localhost",
				Port:            8080,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty host",
			server: Server{
				Port:            8080,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "port too low",
			server: Server{
				Host:            "localhost",
				Port:            0,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "port too high",
			server: Server{
				Host:            "localhost",
				Port:            65536,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid shutdown timeout",
			server: Server{
				Host:            "localhost",
				Port:            8080,
				ShutdownTimeout: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.validateServer()

			if (err != nil) != tt.wantErr {
				t.Fatalf("validateServer() error = %v", err)
			}
		})
	}
}
