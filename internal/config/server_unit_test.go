package config

import (
	"testing"
	"time"
)

func TestServerValidate(t *testing.T) {
	tests := []struct {
		name    string
		server  ServerHTTP
		wantErr bool
	}{
		{
			name: "valid",
			server: ServerHTTP{
				Host:            "localhost",
				Port:            8080,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty host",
			server: ServerHTTP{
				Port:            8080,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "port too low",
			server: ServerHTTP{
				Host:            "localhost",
				Port:            0,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "port too high",
			server: ServerHTTP{
				Host:            "localhost",
				Port:            65536,
				ShutdownTimeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid shutdown timeout",
			server: ServerHTTP{
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
