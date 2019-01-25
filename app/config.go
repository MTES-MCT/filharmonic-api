package app

import (
	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Database database.Config
	Http     httpserver.Config
	Sso      authentication.Config
	LogLevel string `default:"info"`
}

func LoadConfig() Config {
	var c Config
	err := envconfig.Process("filharmonic", &c)
	if err != nil {
		log.Fatal().Err(err).Msg("could not load config")
	}
	return c
}
