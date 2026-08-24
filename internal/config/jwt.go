package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

const (
	pathPrivatKeyEnv        string = "JWT_PATH_PRIVATE_KEY"
	pathPublicKeyEnv        string = "JWT_PATH_PUBLIC_KEY"
	accessExpiresMinutesEnv string = "JWT_ACCESS_EXPIRES"
	refreshExpiresDaysEnv   string = "JWT_REFRESH_EXPIRES"
	issuerEnv               string = "JWT_ISSUER"
	audienceEnv             string = "JWT_AUDIENCE"
)

type JWT interface {
	GetPrivateKey() *rsa.PrivateKey
	GetPublicKey() *rsa.PublicKey
	GetAccessExpires() time.Duration
	GetRefreshExpires() time.Duration
	GetIssuer() string
	GetAudience() string
}

type jwt struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	AccessExpires  time.Duration
	RefreshExpires time.Duration

	Issuer   string
	Audience string
}

func (j jwt) GetPrivateKey() *rsa.PrivateKey {
	return j.privateKey
}

func (j jwt) GetPublicKey() *rsa.PublicKey {
	return j.publicKey
}

func (j jwt) GetAccessExpires() time.Duration {
	return j.AccessExpires
}

func (j jwt) GetRefreshExpires() time.Duration {
	return j.RefreshExpires
}

func (j jwt) GetIssuer() string {
	return j.Issuer
}

func (j jwt) GetAudience() string {
	return j.Audience
}

func (_ jwt) readKeyFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read key file %q: %w", path, err)
	}

	return data, nil
}

func (_ jwt) parsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM")
	}

	if block.Type != "RSA PRIVATE KEY" && block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM type %q", block.Type)
	}

	// PKCS#1
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// PKCS#8
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}

	return rsaKey, nil
}

func (_ jwt) parsePublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM")
	}

	if block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("unexpected PEM type %q", block.Type)
	}

	// PKIX
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		rsaKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not RSA")
		}

		return rsaKey, nil
	}

	// PKCS#1
	rsaKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	return rsaKey, nil
}

func (j *jwt) loadKeys() error {
	privateKeyPath := os.Getenv(pathPrivatKeyEnv)
	publicKeyPath := os.Getenv(pathPublicKeyEnv)

	if privateKeyPath == "" {
		return fmt.Errorf(
			"environment variable %q is required",
			pathPrivatKeyEnv,
		)
	}

	if publicKeyPath == "" {
		return fmt.Errorf(
			"environment variable %q is required",
			pathPublicKeyEnv,
		)
	}

	privatePEM, err := j.readKeyFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("private key: %w", err)
	}

	publicPEM, err := j.readKeyFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("public key: %w", err)
	}

	j.privateKey, err = j.parsePrivateKey(privatePEM)
	if err != nil {
		return fmt.Errorf("private key: %w", err)
	}

	j.publicKey, err = j.parsePublicKey(publicPEM)
	if err != nil {
		return fmt.Errorf("public key: %w", err)
	}

	return nil
}

func (j jwt) validate() error {
	switch {
	case j.Issuer == "":
		return fmt.Errorf("environment variable %q is required", issuerEnv)
	case j.Audience == "":
		return fmt.Errorf("environment variable %q is required", audienceEnv)
	case j.AccessExpires < 1:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			accessExpiresMinutesEnv,
		)
	case j.RefreshExpires < 1:
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			refreshExpiresDaysEnv,
		)
	}

	return nil
}
