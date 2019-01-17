package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/authentication/hash"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
	"github.com/stretchr/testify/require"
)

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

func TestLoginSuccessful(t *testing.T) {
	e, close := tests.Init(t, initTestUsersDB)
	defer close()

	e.POST("/login").WithJSON(&httpserver.Credentials{Email: "existing-user@filharmonic.com", Password: "password"}).
		Expect().Status(http.StatusOK).JSON().Object().ContainsKey("token")

}
func TestLoginFailed(t *testing.T) {
	e, close := tests.Init(t, initTestUsersDB)
	defer close()

	e.POST("/login").WithJSON(&httpserver.Credentials{Email: "missing-user@filharmonic.com", Password: "notpassword"}).
		Expect().Status(http.StatusUnauthorized)

}
