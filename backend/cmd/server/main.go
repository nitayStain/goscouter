package main

import (
	"log"

	"goscouter/backend/internal/server"
)

func main() {
	srv := server.New()
	if err := srv.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

