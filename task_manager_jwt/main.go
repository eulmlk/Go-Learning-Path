package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task_manager/database"
	"task_manager/router"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting server...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	client, err := database.Init()
	if err != nil {
		log.Fatal(err)
	}

	router := router.InitializeRouter(client)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server is running on port 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection closed")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
