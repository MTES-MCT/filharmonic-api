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
	firstEvenement.ValueEqual("type", models.CreationMessage)
	firstEvenement.ValueEqual("data", `{"messageId": 1, "pointDeControleId": 1}`)
	auteur := firstEvenement.Value("auteur").Object()
	auteur.ValueEqual("id", 3)
}

func TestGetEvenementById(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	evenement := tests.AuthInspecteur(e.GET("/evenements/{id}")).WithPath("id", "1").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	evenement.ValueEqual("id", 1)
	evenement.ValueEqual("type", models.CreationMessage)
	evenement.ValueEqual("data", `{"messageId": 1, "pointDeControleId": 1}`)
	auteur := evenement.Value("auteur").Object()
	auteur.ValueEqual("id", 3)
}
