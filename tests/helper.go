package tests

import (
	"log"
	"testing"

	httpexpect "gopkg.in/gavv/httpexpect.v1"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/MTES-MCT/filharmonic-api/authentication/hash"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/stretchr/testify/require"
)

func Init(t *testing.T, initDbFunc func(db *database.Database, assert *require.Assertions)) (*httpexpect.Expect, func()) {
	assert := require.New(t)
	config := app.LoadConfig()
	config.Database.InitSchema = true
	config.Http.Host = "localhost"
	config.Http.Logger = false
	db, server := app.Bootstrap(config)

	initTestUsersDB(db, assert)

	if initDbFunc != nil {
		initDbFunc(db, assert)
	}
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://" + config.Http.Host + ":" + config.Http.Port + "/",
		Reporter: httpexpect.NewAssertReporter(t),
	})
	return e, func() {
		err := server.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}

func initTestUsersDB(db *database.Database, assert *require.Assertions) {
	encodedpassword, err := hash.GenerateFromPassword("password")
	assert.NoError(err)
	user := &models.User{
		Email:    "existing-user@filharmonic.com",
		Password: encodedpassword,
	}
	err = db.Insert(user)
	assert.NoError(err)
}

func Auth(request *httpexpect.Request) *httpexpect.Request {
	return request.WithHeader(httpserver.AuthorizationHeader, "token-1")
}
