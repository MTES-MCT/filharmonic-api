package app

import (
	"context"
	"time"

	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/authentication/cerbere"
	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/MTES-MCT/filharmonic-api/authentication/stubsso"
	"github.com/MTES-MCT/filharmonic-api/cron"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/emails"
	"github.com/MTES-MCT/filharmonic-api/events"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/redis"
	"github.com/MTES-MCT/filharmonic-api/storage"
	"github.com/MTES-MCT/filharmonic-api/templates"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Application struct {
	Config                Config
	DB                    *database.Database
	Redis                 *redis.Client
	Repo                  *database.Repository
	EventsManager         events.EventsManager
	EmailService          domain.EmailService
	Cron                  *cron.CronManager
	Sso                   authentication.Sso
	Sessions              sessions.Sessions
	AuthenticationService *authentication.AuthenticationService
	TemplateService       *templates.TemplateService
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

	log.Info().Msgf("starting in %s mode", a.Config.Mode)
	if a.Config.Mode == ModeDev {
		a.Config.Database.ApplyMigrations = false
	} else if a.Config.Mode == ModeTest {
		a.Config.Database.InitSchema = true
		a.Config.Database.Seeds = true
	}

	db, err := database.New(a.Config.Database)
	if err != nil {
		return err
	}
	a.DB = db

	if a.Config.Mode == ModeTest {
		a.EventsManager = events.NewStub()
	} else {
		a.Redis, err = redis.New(a.Config.Redis)
		if err != nil {
			return err
		}
		a.EventsManager = events.New(a.Config.Events, a.Redis)
	}

	a.Repo = database.NewRepository(a.Config.Repository, db, a.EventsManager)
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
		a.Sessions = sessions.NewRedis(a.Config.Sessions, a.Redis)
	} else if a.Config.Mode == ModeTest {
		a.Sso = stubsso.New(a.Repo)
		a.Sessions = sessions.NewMemory()
	} else {
		a.Sso = cerbere.New(a.Config.Sso)
		a.Sessions = sessions.NewRedis(a.Config.Sessions, a.Redis)
	}

	if a.Config.EnableSmtp {
		emailService, err := emails.New(a.Config.Emails)
		if err != nil {
			return err
		}
		a.EmailService = emailService
	} else {
		a.EmailService = emails.NewStub()
	}

	a.AuthenticationService = authentication.New(a.Repo, a.Sso, a.Sessions)
	a.TemplateService, err = templates.New(a.Config.Templates)
	if err != nil {
		return err
	}
	a.Service = domain.New(a.Config.Service, a.Repo, a.Storage, a.TemplateService, a.EmailService)
	a.Cron, err = cron.New(a.Config.Cron, a.Service)
	if err != nil {
		return err
	}
	a.Server = httpserver.New(a.Config.Http, a.Service, a.AuthenticationService, a.EventsManager)
	return a.Server.Start()
}

func (a *Application) Shutdown() error {
	if a.DB != nil {
		err := a.DB.Close()
		if err != nil {
			return errors.Wrap(err, "could not close db")
		}
	}
	if a.Sessions != nil {
		err := a.Sessions.Close()
		if err != nil {
			return errors.Wrap(err, "could not close sessions")
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
