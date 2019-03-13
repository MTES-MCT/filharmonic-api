package app

import (
	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/MTES-MCT/filharmonic-api/cron"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/emails"
	"github.com/MTES-MCT/filharmonic-api/events"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/redis"
	"github.com/MTES-MCT/filharmonic-api/storage"
	"github.com/MTES-MCT/filharmonic-api/templates"
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
	Emails     emails.Config
	Cron       cron.Config
	Http       httpserver.Config
	Sso        authentication.SsoConfig
	Sessions   sessions.Config
	Redis      redis.Config
	Events     events.Config
	Service    domain.Config
	Storage    storage.Config
	Templates  templates.Config
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
