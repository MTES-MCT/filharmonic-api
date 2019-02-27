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

	tests.AuthInspecteur(e.POST("/inspections/{id}/pointsdecontrole")).WithPath("id", 4).WithJSON(pointControle).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 4).
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

	tests.AuthUser(e.POST("/inspections/{id}/pointsdecontrole"), 5).WithPath("id", 1).WithJSON(pointControle).
		Expect().
		Status(http.StatusForbidden)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	pointsControle := inspection.Value("points_de_controle").Array()
	pointsControle.Length().Equal(2)
}

func TestUpdatePointDeControle(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	pointControle := models.PointDeControle{
		Id:    2,
		Sujet: "Emissions de NOx",
		ReferencesReglementaires: []string{
			"Article 1 de l'arrêté préfectoral du 2/10/2010",
		},
	}

	tests.AuthInspecteur(e.PUT("/pointsdecontrole/{id}")).WithPath("id", 5).WithJSON(pointControle).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 4).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	pointsControle := inspection.Value("points_de_controle").Array()
	pointsControle.Length().Equal(2)
	pointsControle.Last().Object().ValueEqual("sujet", "Emissions de NOx")
	pointsControle.Last().Object().ValueEqual("references_reglementaires", []string{"Article 1 de l'arrêté préfectoral du 2/10/2010"})
	pointsControle.Last().Object().ValueEqual("publie", true)
}

func TestUpdatePointDeControleNotAllowed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	pointControle := models.PointDeControle{
		Id:    2,
		Sujet: "Emissions de NOx",
		ReferencesReglementaires: []string{
			"Article 1 de l'arrêté préfectoral du 2/10/2010",
		},
	}

	tests.AuthUser(e.PUT("/pointsdecontrole/{id}"), 5).WithPath("id", 2).WithJSON(pointControle).
		Expect().
		Status(http.StatusBadRequest)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	pointsControle := inspection.Value("points_de_controle").Array()
	pointsControle.Length().Equal(2)
	pointsControle.Last().Object().ValueEqual("sujet", "Autosurveillance des émissions canalisées de COV")
	pointsControle.Last().Object().ValueEqual("references_reglementaires", []string{"Article 8.2.1.1. de l'arrêté préfectoral du 28 juin 2017"})
}

func TestDeletePointDeControle(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.DELETE("/pointsdecontrole/{id}")).WithPath("id", 5).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 4).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	pointsControle := inspection.Value("points_de_controle").Array()
	pointsControle.Length().Equal(1)
}

func TestPublishPointDeControle(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	pointControle := models.PointDeControle{
		Sujet: "Emissions de NOx",
		ReferencesReglementaires: []string{
			"Article 1 de l'arrêté préfectoral du 2/10/2010",
		},
	}

	pointDeControleId := tests.AuthInspecteur(e.POST("/inspections/{id}/pointsdecontrole")).WithPath("id", 4).WithJSON(pointControle).
		Expect().
		Status(http.StatusOK).JSON().Object().Value("id").Raw()

	tests.AuthInspecteur(e.POST("/pointsdecontrole/{id}/publier")).WithPath("id", pointDeControleId).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 4).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	pointsControle := inspection.Value("points_de_controle").Array()
	pointsControle.Length().Equal(3)
	pointsControle.Last().Object().ValueEqual("publie", true)
}
