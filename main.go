package main

import (
	"log"
	"os"

	"go-backend-valos-id/core/server"
)

func main() {
	app := server.NewApp()

	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
		os.Exit(1)
	}

	addr := getServerAddr()
	log.Printf("Starting server on %s", addr)

	if err := app.Run(addr); err != nil {
		log.Fatalf("Failed to run application: %v", err)
		os.Exit(1)
	}
}

func getServerAddr() string {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return ":" + port
	}
	return ":3210"
}
