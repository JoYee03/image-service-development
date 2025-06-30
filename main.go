package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JoYee03/image-service-development/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port for Cloud Run
	}

	bucketName := os.Getenv("FIREBASE_BUCKET")
	if bucketName == "" {
		log.Fatal("FIREBASE_BUCKET environment variable not set")
	}

	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") // Empty in Cloud Run

	handlers.InitUploadHandler(bucketName, credentialsFile)

	router := http.NewServeMux()
	router.HandleFunc("/testImageUpload", handlers.ImageUploadHTTP)
	router.HandleFunc("/testWatermarkImage", handlers.WatermarkHTTP)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
		// Timeouts to prevent resource exhaustion
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	shutdownComplete := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Received shutdown signal, gracefully terminating...")

		// Give pending requests 20 seconds to complete
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			// Error from closing listeners, or context timeout
			log.Printf("HTTP server shutdown error: %v", err)
		}
		close(shutdownComplete)
	}()

	log.Printf("Starting server on port %s", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}

	<-shutdownComplete
	log.Println("Server shutdown complete")

}