package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAddSuite(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	constat := models.Constat{
		// Id: 5,
		Type: models.TypeConstatObservation,
	}

	tests.AuthInspecteur(e.POST("/pointsdecontrole/{id}/constat")).WithPath("id", 6).WithJSON(constat).
		Expect().
		Status(http.StatusOK)

	suiteInput := models.Suite{
		Type:     models.TypeSuiteObservation,
		Synthese: "Il manque des choses",
	}

	tests.AuthInspecteur(e.POST("/inspections/{id}/suite")).WithPath("id", 4).WithJSON(suiteInput).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 4).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	suiteOutput := inspection.Value("suite").Object()
	suiteOutput.ValueEqual("type", suiteInput.Type)
	suiteOutput.ValueEqual("synthese", suiteInput.Synthese)
}

func TestUpdateSuite(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	suiteInput := models.Suite{
		Id:       1,
		Type:     models.TypeSuiteObservation,
		Synthese: "Il manque des choses",
	}

	tests.AuthInspecteur(e.PUT("/inspections/{id}/suite")).WithPath("id", 1).WithJSON(suiteInput).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	suiteOutput := inspection.Value("suite").Object()
	suiteOutput.ValueEqual("type", suiteInput.Type)
	suiteOutput.ValueEqual("synthese", suiteInput.Synthese)
}

func TestDeleteSuite(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.DELETE("/inspections/{id}/suite")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection.NotContainsKey("suite")
}
