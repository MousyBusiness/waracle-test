package main

import (
	"context"
	"fmt"
	"github.com/mousybusiness/waracle-test/internal/db"
	"github.com/mousybusiness/waracle-test/internal/handler/routes"
	"github.com/mousybusiness/waracle-test/internal/secrets"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	// Connect to Datastore
	database, err := db.ConnectToDatastore(context.Background())
	if err != nil {
		log.Panic(err)
	}

	// Get secret from secret manager
	vault := secrets.NewSecretManagerStore()
	secret, err := vault.GetSecret(context.Background(), "api-key")
	if err != nil {
		log.Panic(err)
	}

	// Store in environment variable for easy access in middleware
	if err := os.Setenv("API_KEY_SECRET", secret); err != nil {
		log.Panic(err)
	}

	// Build routes
	router := routes.BuildRoutes(database)

	fmt.Println("Starting listening on :8080")

	// Start the server with sane default timeouts
	server := http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       90 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
