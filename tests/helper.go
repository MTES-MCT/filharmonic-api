package tests

import (
	"os"
	"testing"
	"time"

	httpexpect "gopkg.in/gavv/httpexpect.v1"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func InitFunc(t *testing.T, initDbFunc func(db *database.Database, assert *require.Assertions)) (*httpexpect.Expect, func()) {
	assert, application := InitFuncDB(t, initDbFunc)

	assert.NoError(application.BootstrapServer())

	httpexpectConfig := httpexpect.Config{
		BaseURL:  "http://" + application.Config.Http.Host + ":" + application.Config.Http.Port + "/",
		Reporter: httpexpect.NewRequireReporter(t),
	}
	if os.Getenv("DEBUG_HTTP") != "" {
		httpexpectConfig.Printers = []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		}
	}
	e := httpexpect.WithConfig(httpexpectConfig)
	return e, func() {
		err := application.Shutdown()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}
}

func Init(t *testing.T) (*httpexpect.Expect, func()) {
	return InitFunc(t, nil)
}

func InitFuncDB(t *testing.T, initDbFunc func(db *database.Database, assert *require.Assertions)) (*require.Assertions, *app.Application) {
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

	if initDbFunc != nil {
		initDbFunc(application.DB, assert)
	}
	return assert, application
}

func InitDB(t *testing.T) (*require.Assertions, *app.Application) {
	return InitFuncDB(t, nil)
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
	return request.WithHeader(httpserver.AuthorizationHeader, authentication.GenerateToken(userId))
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
