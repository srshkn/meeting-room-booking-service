package config

import (
	"strings"
	"testing"
)

func TestCORScfgValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     cors
		wantErr string
	}{
		{
			name: "valid config",
			cfg: cors{
				AllowedOrigins:   []string{"http://localhost:3000", "https://example.com"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Content-Type", "Authorization"},
				AllowCredentials: "true",
			},
		},
		{
			name: "valid credentials false",
			cfg: cors{
				AllowedOrigins:   []string{"http://localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "false",
			},
		},
		{
			name: "method with spaces is valid",
			cfg: cors{
				AllowedOrigins:   []string{"http://localhost"},
				AllowedMethods:   []string{" GET ", " POST "},
				AllowedHeaders:   []string{" Content-Type "},
				AllowCredentials: "true",
			},
		},
		{
			name: "empty origin",
			cfg: cors{
				AllowedOrigins:   []string{""},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "true",
			},
			wantErr: "CORS origin cannot be empty",
		},
		{
			name: "invalid origin URL",
			cfg: cors{
				AllowedOrigins:   []string{"://localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "true",
			},
			wantErr: "invalid CORS origin",
		},
		{
			name: "invalid origin scheme",
			cfg: cors{
				AllowedOrigins:   []string{"ftp://localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "true",
			},
			wantErr: "invalid CORS scheme",
		},
		{
			name: "origin without host",
			cfg: cors{
				AllowedOrigins:   []string{"http:localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "true",
			},
			wantErr: "host must be localhost",
		},
		{
			name: "invalid method",
			cfg: cors{
				AllowedOrigins:   []string{"http://localhost"},
				AllowedMethods:   []string{"TRACE"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "true",
			},
			wantErr: `invalid CORS method "TRACE"`,
		},
		{
			name: "empty header",
			cfg: cors{
				AllowedOrigins:   []string{"http://localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"   "},
				AllowCredentials: "true",
			},
			wantErr: "CORS header cannot be empty",
		},
		{
			name: "invalid credentials",
			cfg: cors{
				AllowedOrigins:   []string{"http://localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "TRUE",
			},
			wantErr: "invalid CORS credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.validate()

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validate() unexpected error: %v", err)
				}

				return
			}

			if err == nil {
				t.Fatalf("validate() expected error containing %q, got nil", tt.wantErr)
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf(
					"validate() error = %q, want error containing %q",
					err.Error(),
					tt.wantErr,
				)
			}
		})
	}
}
