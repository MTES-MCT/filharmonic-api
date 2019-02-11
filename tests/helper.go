package tests

import (
	"os"
	"testing"

	httpexpect "gopkg.in/gavv/httpexpect.v1"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/authentication/mocks"
	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/storage"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func Init(t *testing.T) (*httpexpect.Expect, func()) {
	e, close, _ := InitWithSso(t)
	return e, close
}

func InitWithSso(t *testing.T) (*httpexpect.Expect, func(), *mocks.Sso) {

	assert, a := InitDB(t)

	var err error
	a.Storage, err = storage.New(a.Config.Storage)
	assert.NoError(err)
	a.Sessions = sessions.New()
	sso := new(mocks.Sso)
	a.Sso = sso
	a.AuthenticationService = authentication.New(a.Repo, a.Sso, a.Sessions)
	a.Service = domain.New(a.Repo, a.Storage)
	a.Server = httpserver.New(a.Config.Http, a.Service, a.AuthenticationService)
	assert.NoError(a.Server.Start())
	assert.NoError(initSessions(a.Sessions))

	httpexpectConfig := httpexpect.Config{
		BaseURL:  "http://" + a.Config.Http.Host + ":" + a.Config.Http.Port + "/",
		Reporter: httpexpect.NewRequireReporter(t),
	}
	if os.Getenv("DEBUG_HTTP") != "" {
		httpexpectConfig.Printers = []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		}
	}
	e := httpexpect.WithConfig(httpexpectConfig)
	return e, func() {
		err := a.Shutdown()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}, sso
}

var UserSessions = make(map[int64]string, 0)

func initSessions(sessionsStorage sessions.Sessions) error {
	for i := 1; i < 8; i++ {
		sessionToken, err := sessionsStorage.Add(int64(i))
		if err != nil {
			return err
		}
		UserSessions[int64(i)] = sessionToken
	}
	return nil
}

func InitDB(t *testing.T) (*require.Assertions, *app.Application) {
	assert := require.New(t)
	config := app.LoadConfig()
	config.DevMode = true
	config.Database.InitSchema = true
	config.Http.Host = "localhost"
	config.Http.Logger = false
	config.LogLevel = ""
	application := app.New(config)
	err := application.BootstrapDB()
	assert.NoError(err)

	return assert, application
}

func AuthInspecteur(request *httpexpect.Request) *httpexpect.Request {
	return AuthUser(request, 3)
}

func AuthExploitant(request *httpexpect.Request) *httpexpect.Request {
	return AuthUser(request, 1)
}

func AuthApprobateur(request *httpexpect.Request) *httpexpect.Request {
	return AuthUser(request, 6)
}

func AuthUser(request *httpexpect.Request, userId int64) *httpexpect.Request {
	return request.WithHeader(httpserver.AuthorizationHeader, UserSessions[userId])
}
