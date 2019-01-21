package tests

import (
	"os"
	"testing"

	httpexpect "gopkg.in/gavv/httpexpect.v1"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func InitFunc(t *testing.T, initDbFunc func(db *database.Database, assert *require.Assertions)) (*httpexpect.Expect, func()) {
	assert := require.New(t)
	config := app.LoadConfig()
	config.Database.InitSchema = true
	config.Http.Host = "localhost"
	config.Http.Logger = false
	config.LogLevel = ""
	db, server := app.Bootstrap(config)

	initTestDB(db, assert)

	if initDbFunc != nil {
		initDbFunc(db, assert)
	}
	httpexpectConfig := httpexpect.Config{
		BaseURL:  "http://" + config.Http.Host + ":" + config.Http.Port + "/",
		Reporter: httpexpect.NewAssertReporter(t),
	}
	if os.Getenv("DEBUG_HTTP") != "" {
		httpexpectConfig.Printers = []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		}
	}
	e := httpexpect.WithConfig(httpexpectConfig)
	return e, func() {
		err := server.Close()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}
}

func Init(t *testing.T) (*httpexpect.Expect, func()) {
	return InitFunc(t, nil)
}

func AuthInspecteur(request *httpexpect.Request) *httpexpect.Request {
	return auth(request, 3)
}

func AuthExploitant(request *httpexpect.Request) *httpexpect.Request {
	return auth(request, 1)
}

func AuthApprobateur(request *httpexpect.Request) *httpexpect.Request {
	return auth(request, 6)
}

func auth(request *httpexpect.Request, userId int64) *httpexpect.Request {
	return request.WithHeader(httpserver.AuthorizationHeader, authentication.GenerateToken(userId))
}
