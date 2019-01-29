package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestGetUser(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	user := tests.AuthInspecteur(e.GET("/user")).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	user.ValueEqual("id", 3)
	user.ValueEqual("email", "inspecteur1@filharmonic.com")
	user.ValueEqual("profile", "inspecteur")

}
