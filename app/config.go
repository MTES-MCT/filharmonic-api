package app

import (
	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/storage"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type ModeEnv string

const (
	ModeDev  ModeEnv = "dev"
	ModeTest ModeEnv = "test"
	ModeProd ModeEnv = "prod"
)

type Config struct {
	Database   database.Config
	Repository database.RepositoryConfig
	Http       httpserver.Config
	Sso        authentication.SsoConfig
	Redis      sessions.RedisConfig
	Storage    storage.Config
	LogLevel   string  `default:"info"`
	Mode       ModeEnv `default:"prod"`
}

func LoadConfig() Config {
	var c Config
	err := envconfig.Process("filharmonic", &c)
	if err != nil {
		log.Fatal().Err(err).Msg("could not load config")
	}
	return c
}
