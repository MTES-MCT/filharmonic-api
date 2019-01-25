package app

import (
	"context"
	"time"

	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Application struct {
	Config  Config
	DB      *database.Database
	Repo    *database.Repository
	Sso     *authentication.Sso
	Service *domain.Service
	Server  *httpserver.HttpServer
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

	db, err := database.New(a.Config.Database)
	if err != nil {
		return err
	}
	a.DB = db

	a.Repo = database.NewRepository(db)
	return nil
}

func (a *Application) BootstrapServer() error {
	a.Sso = authentication.New(a.Config.Sso, a.Repo)
	a.Service = domain.New(a.Repo)
	a.Server = httpserver.New(a.Config.Http, a.Service, a.Sso)
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
