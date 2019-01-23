package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAddPointDeControle(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	pointControle := models.PointDeControle{
		Sujet: "Emissions de NOx",
		ReferencesReglementaires: []string{
			"Article 1 de l'arrêté préfectoral du 2/10/2010",
		},
	}

	tests.AuthInspecteur(e.POST("/inspections/{id}/pointsdecontrole")).WithPath("id", 1).WithJSON(pointControle).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	pointsControle := inspection.Value("points_de_controle").Array()
	pointsControle.Length().Equal(3)
	pointsControle.Last().Object().ValueEqual("sujet", "Emissions de NOx")
	pointsControle.Last().Object().ValueEqual("references_reglementaires", []string{"Article 1 de l'arrêté préfectoral du 2/10/2010"})
}

func TestAddPointDeControleWithInspecteurNotAllowed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	pointControle := models.PointDeControle{
		Sujet: "Emissions de NOx",
		ReferencesReglementaires: []string{
			"Article 1 de l'arrêté préfectoral du 2/10/2010",
		},
	}

	tests.AuthUser(e.POST("/inspections/{id}/pointsdecontrole"), 4).WithPath("id", 1).WithJSON(pointControle).
		Expect().
		Status(http.StatusBadRequest)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	pointsControle := inspection.Value("points_de_controle").Array()
	pointsControle.Length().Equal(2)
}
