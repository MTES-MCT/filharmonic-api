package tests

import (
	"os"
	"strconv"
	"testing"
	"time"

	httpexpect "gopkg.in/gavv/httpexpect.v1"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/authentication/mocks"
	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func Init(t *testing.T) (*httpexpect.Expect, func()) {
	e, close, _ := InitWithSso(t)
	return e, close
}

func InitWithSso(t *testing.T) (*httpexpect.Expect, func(), *mocks.Sso) {

	assert, a := InitDB(t)

	a.Sessions = sessions.New()
	sso := new(mocks.Sso)
	a.Sso = sso
	a.AuthenticationService = authentication.New(a.Repo, a.Sso, a.Sessions)
	a.Service = domain.New(a.Repo)
	a.Server = httpserver.New(a.Config.Http, a.Service, a.AuthenticationService)
	assert.NoError(a.Server.Start())
	initSessions(a.Sessions)

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

func initSessions(sessionsStorage sessions.Sessions) {
	for i := 1; i < 8; i++ {
		sessionsStorage.Set(GenerateToken(int64(i)), int64(i))
	}
}

func InitDB(t *testing.T) (*require.Assertions, *app.Application) {
	assert := require.New(t)
	config := app.LoadConfig()
	config.Database.InitSchema = true
	config.Http.Host = "localhost"
	config.Http.Logger = false
	config.LogLevel = ""
	application := app.New(config)
	err := application.BootstrapDB()
	assert.NoError(err)

	seedsTestDB(application.DB, assert)

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
	return request.WithHeader(httpserver.AuthorizationHeader, GenerateToken(userId))
}

func Date(datestr string) time.Time {
	date, err := time.Parse("2006-01-02", datestr)
	if err != nil {
		log.Fatal().Msgf("unable to parse date: %v", err)
	}
	return date
}

func DateTime(datestr string) time.Time {
	date, err := time.Parse("2006-01-02T15:04:05", datestr)
	if err != nil {
		log.Fatal().Msgf("unable to parse date: %v", err)
	}
	return date
}

func GenerateToken(id int64) string {
	return "token-" + strconv.FormatInt(id, 10)
}
