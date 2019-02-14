package app

import (
	"context"
	"time"

	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/authentication/cerbere"
	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/MTES-MCT/filharmonic-api/authentication/stubsso"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Application struct {
	Config                Config
	DB                    *database.Database
	Repo                  *database.Repository
	Sso                   authentication.Sso
	Sessions              sessions.Sessions
	AuthenticationService *authentication.AuthenticationService
	Service               *domain.Service
	Storage               *storage.FileStorage
	Server                *httpserver.HttpServer
}

func New(config Config) *Application {
	return &Application{
		Config: config,
	}
}

func (a *Application) BootstrapDB() error {
	logLevel, err := zerolog.ParseLevel(a.Config.LogLevel)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(logLevel)

	log.Info().Msgf("starting in mode %s", a.Config.Mode)
	if a.Config.Mode == ModeTest {
		a.Config.Database.Seeds = true
	}

	db, err := database.New(a.Config.Database)
	if err != nil {
		return err
	}
	a.DB = db

	a.Repo = database.NewRepository(a.Config.Repository, db)
	return nil
}

func (a *Application) BootstrapServer() error {
	storage, err := storage.New(a.Config.Storage)
	if err != nil {
		return err
	}
	a.Storage = storage
	if a.Config.Mode == ModeDev {
		a.Sso = stubsso.New(a.Repo)
		redisSessions, err := sessions.NewRedis(a.Config.Redis)
		if err != nil {
			return err
		}
		a.Sessions = redisSessions
	} else if a.Config.Mode == ModeTest {
		a.Sso = stubsso.New(a.Repo)
		a.Sessions = sessions.NewMemory()
	} else {
		a.Sso = cerbere.New(a.Config.Sso)
		redisSessions, err := sessions.NewRedis(a.Config.Redis)
		if err != nil {
			return err
		}
		a.Sessions = redisSessions
	}
	a.AuthenticationService = authentication.New(a.Repo, a.Sso, a.Sessions)
	a.Service = domain.New(a.Repo, a.Storage)
	a.Server = httpserver.New(a.Config.Http, a.Service, a.AuthenticationService)
	return a.Server.Start()
}

func (a *Application) Shutdown() error {
	if a.DB != nil {
		err := a.DB.Close()
		if err != nil {
			return errors.Wrap(err, "could not close db")
		}
	}
	if a.Server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := a.Server.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "could not close server")
		}
	}
	return nil
}
