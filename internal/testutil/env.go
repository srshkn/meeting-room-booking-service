package testutil

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

const (
	testPathPrivateKeyEnv string = "JWT_PATH_PRIVATE_KEY"
	testPathPublicKeyEnv  string = "JWT_PATH_PUBLIC_KEY"
	testPathPrivatePem    string = "/secrets/test-private.pem"
	testPathPublicPem     string = "/secrets/test-public.pem"
)

func NewTestEnv() error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		slog.Error(
			"failed to caller",
			slog.Any("error", errors.New("failed to get current file path")),
		)
		os.Exit(1)
	}

	root := filepath.Join(filepath.Dir(filename), "../..")

	envPath := filepath.Join(root, ".test.env")

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return err
	}

	fmt.Printf("envPath: %v\n", envPath)

	if err := godotenv.Overload(envPath); err != nil {
		slog.Error(
			"failed to caller",
			slog.Any("error", errors.New("failed to get current file path")),
		)
		os.Exit(1)
	}

	os.Setenv(
		testPathPrivateKeyEnv,
		filepath.ToSlash(
			filepath.Join(
				envPath,
				"..",
				testPathPrivatePem,
			),
		),
	)

	os.Setenv(
		testPathPublicKeyEnv,
		filepath.ToSlash(
			filepath.Join(
				envPath,
				"..",
				testPathPublicPem,
			),
		),
	)

	return nil
}
