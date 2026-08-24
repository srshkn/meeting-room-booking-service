package config

import (
	"net/http"
	"testing"
)

func TestCookieValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     cookie
		wantErr string
	}{
		{
			name: "valid",
			cfg: cookie{
				RefreshTokenName: "refresh_token",
				Path:             "/",
				Domain:           "example.com",
				Secure:           true,
				HTTPOnly:         true,
				sameSite:         http.SameSiteLaxMode,
			},
		},
		{
			name: "empty refresh token name",
			cfg: cookie{
				Path:     "/",
				sameSite: http.SameSiteLaxMode,
			},
			wantErr: "refresh token cookie name is empty",
		},
		{
			name: "empty path",
			cfg: cookie{
				RefreshTokenName: "refresh_token",
				sameSite:         http.SameSiteLaxMode,
			},
			wantErr: `invalid cookie path: ""`,
		},
		{
			name: "path without slash",
			cfg: cookie{
				RefreshTokenName: "refresh_token",
				Path:             "auth",
				sameSite:         http.SameSiteLaxMode,
			},
			wantErr: `invalid cookie path: "auth"`,
		},
		{
			name: "domain with spaces",
			cfg: cookie{
				RefreshTokenName: "refresh_token",
				Path:             "/",
				Domain:           " example.com ",
				sameSite:         http.SameSiteLaxMode,
			},
			wantErr: "cookie domain contains leading or trailing spaces",
		},
		{
			name: "domain with invalid characters",
			cfg: cookie{
				RefreshTokenName: "refresh_token",
				Path:             "/",
				Domain:           "example.com/path",
				sameSite:         http.SameSiteLaxMode,
			},
			wantErr: "cookie domain contains invalid characters",
		},
		{
			name: "invalid same site",
			cfg: cookie{
				RefreshTokenName: "refresh_token",
				Path:             "/",
				sameSite:         http.SameSite(100),
			},
			wantErr: "invalid cookie same site",
		},
		{
			name: "same site none without secure",
			cfg: cookie{
				RefreshTokenName: "refresh_token",
				Path:             "/",
				Secure:           false,
				sameSite:         http.SameSiteNoneMode,
			},
			wantErr: "cookie secure must be enabled when same site is none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.validate()

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validate() error = %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("validate() error = nil, want %q", tt.wantErr)
			}

			if err.Error() != tt.wantErr {
				t.Fatalf("validate() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCookieParseSameSite(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    http.SameSite
		wantErr bool
	}{
		{
			name:  "default",
			value: "default",
			want:  http.SameSiteDefaultMode,
		},
		{
			name:  "lax",
			value: "lax",
			want:  http.SameSiteLaxMode,
		},
		{
			name:  "strict",
			value: "strict",
			want:  http.SameSiteStrictMode,
		},
		{
			name:  "none",
			value: "none",
			want:  http.SameSiteNoneMode,
		},
		{
			name:  "case insensitive",
			value: "LaX",
			want:  http.SameSiteLaxMode,
		},
		{
			name:    "invalid",
			value:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := cookie{
				StringSameSite: tt.value,
			}

			err := cfg.parseSameSite()

			if tt.wantErr {
				if err == nil {
					t.Fatal("parseSameSite() error = nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("parseSameSite() error = %v", err)
			}

			if cfg.sameSite != tt.want {
				t.Fatalf(
					"sameSite = %v, want %v",
					cfg.sameSite,
					tt.want,
				)
			}
		})
	}
}
