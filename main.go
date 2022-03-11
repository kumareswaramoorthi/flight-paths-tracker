package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kumareswaramoorthi/flight-paths-tracker/api/router"
)

func main() {

	ginEngine := router.SetupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: ginEngine,
	}

	// Graceful shut down of server
	graceful := make(chan os.Signal)
	signal.Notify(graceful, syscall.SIGINT)
	signal.Notify(graceful, syscall.SIGTERM)
	go func() {
		<-graceful
		log.Println("Shutting down ctrl...")
		ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelFunc()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Could not do graceful shutdown: %v\n", err)
		}
	}()

	log.Println("Listening server on 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("")
	}
	log.Println("Server gracefully stopped...")
}
