package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestListInspecteurs(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	users := tests.AuthInspecteur(e.GET("/inspecteurs")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	users.Length().Equal(4)
	firstUser := users.First().Object()
	firstUser.ValueEqual("id", 3)
	firstUser.ValueEqual("email", "inspecteur1@filharmonic.com")
	firstUser.ValueEqual("profile", "inspecteur")
}
