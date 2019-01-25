package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/authentication/sessions"
	"github.com/MTES-MCT/filharmonic-api/httpserver"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestLoginSuccessful(t *testing.T) {
	e, close := tests.Init(t)
	defer close()
	// TODO init cerbere mock

	e.POST("/login").WithJSON(&httpserver.LoginHTTPRequest{Ticket: "ticket-exploitant1"}).
		Expect().Status(http.StatusOK).JSON().Object().ContainsKey("token")

}
func TestLoginFailed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()
	// TODO init cerbere mock

	e.POST("/login").WithJSON(&httpserver.LoginHTTPRequest{Ticket: "invalid-ticket"}).
		Expect().Status(http.StatusUnauthorized)
}

func TestAuthenticateSuccessful(t *testing.T) {
	e, close := tests.Init(t)
	defer close()
	sessions.Set("valid-token", 3)

	e.POST("/authenticate").WithJSON(&httpserver.AuthenticateHTTPRequest{Token: "valid-token"}).
		Expect().Status(http.StatusOK)
}

func TestAuthenticateFailed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	e.POST("/authenticate").WithJSON(&httpserver.AuthenticateHTTPRequest{Token: "invalid-token"}).
		Expect().Status(http.StatusUnauthorized)
}
