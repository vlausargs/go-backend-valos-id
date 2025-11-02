package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// App represents the application structure
type App struct {
	server *Server
}

// NewApp creates a new application instance
func NewApp() *App {
	return &App{
		server: NewServer(),
	}
}

// Initialize initializes the application
func (a *App) Initialize() error {
	if err := a.server.Initialize(); err != nil {
		return err
	}
	return nil
}

// Run starts the application
func (a *App) Run(addr string) error {
	// Start server in a goroutine
	go func() {
		if err := a.server.Start(addr); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Close database connection and perform cleanup
	if err := a.server.Close(); err != nil {
		log.Printf("Error during shutdown: %v", err)
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}
