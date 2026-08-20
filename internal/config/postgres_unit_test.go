package config

import (
	"strings"
	"testing"
	"time"
)

func TestPostgresValidate(t *testing.T) {
	validPostgres := func() Postgres {
		return Postgres{
			Host:            "localhost",
			Port:            5432,
			User:            "postgres",
			Password:        "password",
			Database:        "meeting_rooms",
			MaxOpenConns:    10,
			MinOpenConns:    2,
			ConnMaxLifetime: 30 * time.Minute,
			MaxConnIdleTime: 5 * time.Minute,
		}
	}

	tests := []struct {
		name    string
		prepare func(*Postgres)
		wantErr string
	}{
		{
			name:    "valid config",
			prepare: func(_ *Postgres) {},
		},
		{
			name: "empty host",
			prepare: func(p *Postgres) {
				p.Host = ""
			},
			wantErr: `environment variable "POSTGRES_HOST" is required`,
		},
		{
			name: "port less than minimum",
			prepare: func(p *Postgres) {
				p.Port = 0
			},
			wantErr: `environment variable "POSTGRES_PORT" must be between 1 and 65535`,
		},
		{
			name: "port greater than maximum",
			prepare: func(p *Postgres) {
				p.Port = 65536
			},
			wantErr: `environment variable "POSTGRES_PORT" must be between 1 and 65535`,
		},
		{
			name: "empty user",
			prepare: func(p *Postgres) {
				p.User = ""
			},
			wantErr: `environment variable "POSTGRES_USER" is required`,
		},
		{
			name: "empty password",
			prepare: func(p *Postgres) {
				p.Password = ""
			},
			wantErr: `environment variable "POSTGRES_PASSWORD" is required`,
		},
		{
			name: "empty database",
			prepare: func(p *Postgres) {
				p.Database = ""
			},
			wantErr: `environment variable "POSTGRES_DB" is required`,
		},
		{
			name: "max open conns is zero",
			prepare: func(p *Postgres) {
				p.MaxOpenConns = 0
			},
			wantErr: `environment variable "MAX_OPEN_CONNS"`,
		},
		{
			name: "min open conns is zero",
			prepare: func(p *Postgres) {
				p.MinOpenConns = 0
			},
			wantErr: `environment variable "MIN_OPEN_CONNS"`,
		},
		{
			name: "conn max lifetime is zero",
			prepare: func(p *Postgres) {
				p.ConnMaxLifetime = 0
			},
			wantErr: `environment variable "CONN_MAX_LIFETIME"`,
		},
		{
			name: "max conn idle time is zero",
			prepare: func(p *Postgres) {
				p.MaxConnIdleTime = 0
			},
			wantErr: `environment variable "MAX_CONN_IDLE_TIME"`,
		},
		{
			name: "user exceeds 63 bytes",
			prepare: func(p *Postgres) {
				p.User = strings.Repeat("a", 64)
			},
			wantErr: `environment variable "POSTGRES_USER" must not exceed 63 bytes`,
		},
		{
			name: "database exceeds 63 bytes",
			prepare: func(p *Postgres) {
				p.Database = strings.Repeat("a", 64)
			},
			wantErr: `environment variable "POSTGRES_DB" must not exceed 63 bytes`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postgres := validPostgres()
			tt.prepare(&postgres)

			err := postgres.validatePostgres()

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validatePostgres() unexpected error: %v", err)
				}

				return
			}

			if err == nil {
				t.Fatalf("validatePostgres() expected error containing %q, got nil", tt.wantErr)
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf(
					"validatePostgres() error = %q, want error containing %q",
					err.Error(),
					tt.wantErr,
				)
			}
		})
	}
}
