package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"task_manager/api/router"
	"task_manager/database"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting server...")

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database connection
	client, err := database.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize router
	router := router.InitializeRouter(client)
	database.CreateRootUser(client)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server is running on port 8080")

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Shutting down server...")

	// Create a context with a timeout for the graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close database connection
	if err := client.Disconnect(context.Background()); err != nil {
		log.Println("Error closing database connection:", err)
	} else {
		log.Println("Database connection closed")
	}

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")
}
