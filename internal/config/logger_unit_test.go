package config

import (
	"testing"
)

func TestLogManagerValidate(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{
			name:    "dev format",
			format:  "dev",
			wantErr: false,
		},
		{
			name:    "local format",
			format:  "local",
			wantErr: false,
		},
		{
			name:    "prod format",
			format:  "prod",
			wantErr: false,
		},
		{
			name:    "empty format",
			format:  "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			format:  "test",
			wantErr: true,
		},
		{
			name:    "uppercase format",
			format:  "DEV",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logger{
				Format: tt.format,
			}

			err := logger.validate()

			if tt.wantErr && err == nil {
				t.Fatalf("validate() expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("validate() unexpected error: %v", err)
			}
		})
	}
}
