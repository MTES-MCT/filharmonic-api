package main

import (
	"os"
	"os/signal"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/rs/zerolog/log"
)

func main() {
	config := app.LoadConfig()
	application := app.New(config)
	err := application.BootstrapDB()
	if err != nil {
		log.Fatal().Err(err).Msg("could not start db")
	}

	err = application.BootstrapServer()
	if err != nil {
		log.Fatal().Err(err).Msg("could not start http server")
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Warn().Msg("shutting down application")

	if err := application.Shutdown(); err != nil {
		log.Fatal().Err(err).Msg("application shutdown error")
	}
	log.Warn().Msg("application shutdown")
}
