package app

import (
	"log"
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

func LoadConfig() Config {
	var c Config
	err := envconfig.Process("filharmonic", &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	return c
}

func Bootstrap(c Config) (*database.Database, *http.Server) {
	logLevel, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		log.Fatal(err.Error())
	}
	zerolog.SetGlobalLevel(logLevel)

	db, err := database.New(c.Database)
	if err != nil {
		log.Fatal("Database error:", err)
	}

	repo := database.NewRepository(db)
	service := domain.New(repo)
	httpsrv := httpserver.New(c.Http, service)
	server := httpsrv.Start()
	return db, server
}
