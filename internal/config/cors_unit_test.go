package config

import (
	"strings"
	"testing"
)

func TestCORScfgValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     CORScfg
		wantErr string
	}{
		{
			name: "valid config",
			cfg: CORScfg{
				AllowedOrigins:   []string{"http://localhost:3000", "https://example.com"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Content-Type", "Authorization"},
				AllowCredentials: "true",
			},
		},
		{
			name: "valid credentials false",
			cfg: CORScfg{
				AllowedOrigins:   []string{"http://localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "false",
			},
		},
		{
			name: "method with spaces is valid",
			cfg: CORScfg{
				AllowedOrigins:   []string{"http://localhost"},
				AllowedMethods:   []string{" GET ", " POST "},
				AllowedHeaders:   []string{" Content-Type "},
				AllowCredentials: "true",
			},
		},
		{
			name: "empty origin",
			cfg: CORScfg{
				AllowedOrigins:   []string{""},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "true",
			},
			wantErr: "CORS origin cannot be empty",
		},
		{
			name: "invalid origin URL",
			cfg: CORScfg{
				AllowedOrigins:   []string{"://localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "true",
			},
			wantErr: "invalid CORS origin",
		},
		{
			name: "invalid origin scheme",
			cfg: CORScfg{
				AllowedOrigins:   []string{"ftp://localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "true",
			},
			wantErr: "invalid CORS scheme",
		},
		{
			name: "origin without host",
			cfg: CORScfg{
				AllowedOrigins:   []string{"http:localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "true",
			},
			wantErr: "host must be localhost",
		},
		{
			name: "invalid method",
			cfg: CORScfg{
				AllowedOrigins:   []string{"http://localhost"},
				AllowedMethods:   []string{"TRACE"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: "true",
			},
			wantErr: `invalid CORS method "TRACE"`,
		},
		{
			name: "empty header",
			cfg: CORScfg{
				AllowedOrigins:   []string{"http://localhost"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"   "},
				AllowCredentials: "true",
			},
			wantErr: "CORS header cannot be empty",
		},
		{
			name: "invalid credentials",
			cfg: CORScfg{
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
