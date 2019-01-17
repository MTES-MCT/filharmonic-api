package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestLoginSuccessful(t *testing.T) {
	e, close := tests.Init(t, nil)
	defer close()

	e.POST("/login").WithJSON(&httpserver.Credentials{Email: "exploitant1@filharmonic.com", Password: "password1"}).
		Expect().Status(http.StatusOK).JSON().Object().ContainsKey("token")

}
func TestLoginFailed(t *testing.T) {
	e, close := tests.Init(t, nil)
	defer close()

	e.POST("/login").WithJSON(&httpserver.Credentials{Email: "missing-user@filharmonic.com", Password: "notpassword"}).
		Expect().Status(http.StatusUnauthorized)

}
