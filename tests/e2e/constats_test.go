package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAddConstat(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	constatInput := models.Constat{
		Type:      models.TypeConstatObservation,
		Remarques: "Il manque des choses",
	}

	tests.AuthInspecteur(e.POST("/pointsdecontrole/{id}/constat")).WithPath("id", 6).WithJSON(constatInput).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 4).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	lastPointDeControle := inspection.Value("points_de_controle").Array().Last().Object()
	constatOutput := lastPointDeControle.Value("constat").Object()
	constatOutput.ValueEqual("type", constatInput.Type)
	constatOutput.ValueEqual("remarques", constatInput.Remarques)
}

func TestDeleteConstat(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.DELETE("/pointsdecontrole/{id}/constat")).WithPath("id", 5).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 4).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	firstPointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	firstPointDeControle.NotContainsKey("constat")
}

func TestResoudreConstat(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.POST("/pointsdecontrole/{id}/constat/resoudre")).WithPath("id", 7).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 5).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	lastPointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	constatOutput := lastPointDeControle.Value("constat").Object()
	constatOutput.ValueNotEqual("date_resolution", "")
}
