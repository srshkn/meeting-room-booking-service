package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestJWTValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     jwt
		wantErr string
	}{
		{
			name: "valid",
			cfg: jwt{
				Issuer:         "mrb-service",
				Audience:       "mrb-client",
				AccessExpires:  15 * time.Minute,
				RefreshExpires: 30 * 24 * time.Hour,
			},
		},
		{
			name: "empty issuer",
			cfg: jwt{
				Audience:       "mrb-client",
				AccessExpires:  time.Minute,
				RefreshExpires: time.Hour,
			},
			wantErr: `environment variable "JWT_ISSUER" is required`,
		},
		{
			name: "empty audience",
			cfg: jwt{
				Issuer:         "mrb-service",
				AccessExpires:  time.Minute,
				RefreshExpires: time.Hour,
			},
			wantErr: `environment variable "JWT_AUDIENCE" is required`,
		},
		{
			name: "invalid access expires",
			cfg: jwt{
				Issuer:         "mrb-service",
				Audience:       "mrb-client",
				RefreshExpires: time.Hour,
			},
			wantErr: `environment variable "JWT_ACCESS_EXPIRES" must be greater than 0`,
		},
		{
			name: "invalid refresh expires",
			cfg: jwt{
				Issuer:        "mrb-service",
				Audience:      "mrb-client",
				AccessExpires: time.Minute,
			},
			wantErr: `environment variable "JWT_REFRESH_EXPIRES" must be greater than 0`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.validate()

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}

			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestJWTParsePrivateKey(t *testing.T) {
	key := generateRSAKey(t)

	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal PKCS8: %v", err)
	}

	tests := []struct {
		name    string
		data    []byte
		wantErr string
	}{
		{
			name: "PKCS1",
			data: pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(key),
			}),
		},
		{
			name: "PKCS8",
			data: pem.EncodeToMemory(&pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: pkcs8Bytes,
			}),
		},
		{
			name:    "invalid PEM",
			data:    []byte("invalid"),
			wantErr: "invalid PEM",
		},
		{
			name: "unexpected type",
			data: pem.EncodeToMemory(&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: []byte("invalid"),
			}),
			wantErr: "unexpected PEM type",
		},
	}

	var j jwt

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := j.parsePrivateKey(tt.data)

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if key == nil {
					t.Fatal("expected key, got nil")
				}
				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want contains %q", err, tt.wantErr)
			}
		})
	}
}

func TestJWTParsePublicKey(t *testing.T) {
	privateKey := generateRSAKey(t)
	publicKey := &privateKey.PublicKey

	pkixBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("marshal PKIX: %v", err)
	}

	tests := []struct {
		name    string
		data    []byte
		wantErr string
	}{
		{
			name: "PKIX",
			data: pem.EncodeToMemory(&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: pkixBytes,
			}),
		},
		{
			name: "PKCS1",
			data: pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(publicKey),
			}),
		},
		{
			name:    "invalid PEM",
			data:    []byte("invalid"),
			wantErr: "invalid PEM",
		},
		{
			name: "unexpected type",
			data: pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: []byte("invalid"),
			}),
			wantErr: "unexpected PEM type",
		},
	}

	var j jwt

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := j.parsePublicKey(tt.data)

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if key == nil {
					t.Fatal("expected key, got nil")
				}
				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want contains %q", err, tt.wantErr)
			}
		})
	}
}

func TestJWTLoadKeys(t *testing.T) {
	privateKey := generateRSAKey(t)

	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}

	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	})

	tests := []struct {
		name    string
		setup   func(t *testing.T)
		wantErr string
	}{
		{
			name: "success",
			setup: func(t *testing.T) {
				t.Setenv(pathPrivatKeyEnv, writeTempFile(t, privatePEM))
				t.Setenv(pathPublicKeyEnv, writeTempFile(t, publicPEM))
			},
		},
		{
			name: "private env missing",
			setup: func(t *testing.T) {
				t.Setenv(pathPrivatKeyEnv, "")
				t.Setenv(pathPublicKeyEnv, "public.pem")
			},
			wantErr: `environment variable "JWT_PATH_PRIVATE_KEY" is required`,
		},
		{
			name: "public env missing",
			setup: func(t *testing.T) {
				t.Setenv(pathPrivatKeyEnv, "private.pem")
				t.Setenv(pathPublicKeyEnv, "")
			},
			wantErr: `environment variable "JWT_PATH_PUBLIC_KEY" is required`,
		},
		{
			name: "key file not found",
			setup: func(t *testing.T) {
				t.Setenv(pathPrivatKeyEnv, "missing.pem")
				t.Setenv(pathPublicKeyEnv, "public.pem")
			},
			wantErr: "private key: read key file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			var j jwt
			err := j.loadKeys()

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if j.privateKey == nil || j.publicKey == nil {
					t.Fatal("keys were not loaded")
				}

				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want contains %q", err, tt.wantErr)
			}
		})
	}
}

func generateRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}

	return key
}

func writeTempFile(t *testing.T, data []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "key.pem")

	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	return path
}
