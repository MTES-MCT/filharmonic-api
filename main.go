package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/MTES-MCT/filharmonic-api/app"
)

func main() {
	// chargement config
	config := app.LoadConfig()
	db, server := app.Bootstrap(config)
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal("Unable to close DB", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
