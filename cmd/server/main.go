package main

import (
	"fmt"

	"mrb-service/internal/config"
)

func main() {
	cfg, help, err := config.New()
	if err != nil {
		fmt.Printf("help: %s;\nerr: %s\n", help, err)
		return
	}

	fmt.Printf("Server: host: %s, port: %d", cfg.Server.Host, cfg.Server.Port)
}
