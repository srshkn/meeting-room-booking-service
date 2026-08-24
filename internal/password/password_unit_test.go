package password

import (
	"errors"
	"strings"
	"testing"
)

func mustHash(t *testing.T, password string) string {
	t.Helper()

	got, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash(%q): %v", password, err)
	}

	return got
}

func TestHash(t *testing.T) {
	const password = "my-secret-password"

	first := mustHash(t, password)
	second := mustHash(t, password)

	if !strings.HasPrefix(first, "$argon2id$") {
		t.Errorf("Hash() = %q, want argon2id hash", first)
	}
	if first == second {
		t.Error("Hash() produced identical hashes")
	}

	ok, err := Compare(password, first)
	if err != nil {
		t.Fatalf("Compare(): %v", err)
	}
	if !ok {
		t.Error("Compare() = false, want true")
	}
}

func TestCompare(t *testing.T) {
	const password = "my-secret-password"
	encoded := mustHash(t, password)

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"correct password", password, true},
		{"incorrect password", "wrong-password", false},
		{"empty password", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compare(tt.password, encoded)
			if err != nil {
				t.Fatalf("Compare(): %v", err)
			}
			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompare_InvalidHash(t *testing.T) {
	tests := []struct {
		name string
		hash string
		want error
	}{
		{"empty hash", "", ErrInvalidHash},
		{"too few parts", "$argon2id$v=19$m=65536,t=3,p=4", ErrInvalidHash},
		{"unsupported algorithm", "$bcrypt$v=19$m=65536,t=3,p=4$YWJj$YWJj",
			ErrUnsupportedHash},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Compare("password", tt.hash)

			if !errors.Is(err, tt.want) {
				t.Errorf("Compare() error = %v, want %v", err, tt.want)
			}
		})
	}
}
