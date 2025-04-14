package main

import (
	"cchecker/handlers"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/client"
)

func main() {
	log.Println("Starting healthchecker service...")
	log.Println("Connecting to Docker daemon...")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	// Ping the Docker daemon to verify connection
	_, err = cli.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to ping Docker daemon: %v", err)
	}
	log.Println("Successfully connected to Docker daemon")

	healthHandler := handlers.HealthHandler{
		CLI: cli,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Define routes
	http.HandleFunc("/health", healthHandler.MainHandler)

	// Redirect root to dashboard
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/health", http.StatusSeeOther)
	})

	log.Printf("Starting HTTP server on port %s...", port)
	log.Printf("Dashboard available at http://localhost:%s/health", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
