package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAuthentication(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.GET("/ping")).
		Expect().
		Status(http.StatusOK)
}

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
	favoris := user.Value("favoris").Array()
	favoris.Length().Equal(1)
	favori := favoris.First().Object()
	favori.NotContainsKey("inspecteurs")
	favori.NotContainsKey("commentaires")
	favori.NotContainsKey("points_de_controle")
	favori.NotContainsKey("suite")
	favori.ValueEqual("id", 1)
	favori.ValueEqual("date", "2018-09-01")
	etablissement := favori.Value("etablissement").Object()
	etablissement.ValueEqual("adresse", "1 rue des fleurs 75000 Paris")
	etablissement.ValueEqual("nom", "Nom 1")
}
