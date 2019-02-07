package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/MTES-MCT/filharmonic-api/database/importcsv"
	"github.com/rs/zerolog/log"
)

func main() {
	importEtablissements := flag.String("import-etablissements", "", "Importe des établissements à partir d'un CSV")
	importInspecteurs := flag.String("import-inspecteurs", "", "Importe des inspecteurs à partir d'un CSV")
	flag.Parse()

	config := app.LoadConfig()
	application := app.New(config)
	err := application.BootstrapDB()
	if err != nil {
		log.Fatal().Err(err).Msg("could not start db")
	}

	if *importEtablissements != "" {
		err = importcsv.LoadEtablissementsCSV(*importEtablissements, application.DB)
		if err != nil {
			log.Fatal().Err(err).Msg("could not import CSV")
		}
		os.Exit(0)
	}

	if *importInspecteurs != "" {
		err = importcsv.LoadInspecteursCSV(*importInspecteurs, application.DB)
		if err != nil {
			log.Fatal().Err(err).Msg("could not import CSV")
		}
		os.Exit(0)
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
