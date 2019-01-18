package tests

import (
	"testing"

	httpexpect "gopkg.in/gavv/httpexpect.v1"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/MTES-MCT/filharmonic-api/authentication/hash"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func Init(t *testing.T, initDbFunc func(db *database.Database, assert *require.Assertions)) (*httpexpect.Expect, func()) {
	assert := require.New(t)
	config := app.LoadConfig()
	config.Database.InitSchema = true
	config.Http.Host = "localhost"
	config.Http.Logger = false
	config.LogLevel = ""
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
			log.Fatal().Msg(err.Error())
		}
	}
}

func initTestUsersDB(db *database.Database, assert *require.Assertions) {
	encodedpassword1, err := hash.GenerateFromPassword("password1")
	assert.NoError(err)
	encodedpassword2, err := hash.GenerateFromPassword("password2")
	assert.NoError(err)
	users := []interface{}{
		&models.User{
			Id:       1,
			Email:    "exploitant1@filharmonic.com",
			Password: encodedpassword1,
			Profile:  models.ProfilExploitant,
		},
		&models.User{
			Id:       2,
			Email:    "exploitant2@filharmonic.com",
			Password: encodedpassword2,
			Profile:  models.ProfilExploitant,
		},
		&models.User{
			Id:       3,
			Email:    "inspecteur1@filharmonic.com",
			Password: encodedpassword1,
			Profile:  models.ProfilInspecteur,
		},
		&models.User{
			Id:       4,
			Email:    "inspecteur2@filharmonic.com",
			Password: encodedpassword2,
			Profile:  models.ProfilInspecteur,
		},
	}
	err = db.Insert(users...)
	assert.NoError(err)
}

func Auth(request *httpexpect.Request) *httpexpect.Request {
	return request.WithHeader(httpserver.AuthorizationHeader, "token-1")
}
