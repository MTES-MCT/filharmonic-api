package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestListEvenements(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	evenements := tests.AuthInspecteur(e.GET("/evenements")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	evenements.Length().Equal(4)
	firstEvenement := evenements.First().Object()
	firstEvenement.ValueEqual("id", 1)
	firstEvenement.ValueEqual("type", models.EvenementCreationMessage)
	data := firstEvenement.Value("data").Object()
	data.ValueEqual("message_id", 1)
	data.ValueEqual("point_de_controle_id", 1)
	auteur := firstEvenement.Value("auteur").Object()
	auteur.ValueEqual("id", 3)
}
