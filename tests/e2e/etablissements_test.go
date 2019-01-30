package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestFindEtablissementsByS3IC(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthExploitant(e.GET("/etablissements")).WithQuery("s3ic", "23").
		Expect().
		Status(http.StatusOK).
		JSON().Array().
		Element(0).Object().ValueEqual("s3ic", "1234")
}

func TestFindEtablissementsOwnedByExploitant(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	results := tests.AuthExploitant(e.GET("/etablissements")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	results.Length().Equal(2)
	results.First().Object().ValueEqual("s3ic", "1234")
	results.Last().Object().ValueEqual("s3ic", "4567")
}

func TestGetEtablissementByIdOwnedByExploitant(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	etablissement := tests.AuthExploitant(e.GET("/etablissements/{id}")).WithPath("id", "1").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	etablissement.ValueEqual("id", 1)
	inspections := etablissement.Value("inspections").Array()
	inspections.Length().Equal(1)
	inspection := inspections.First().Object()
	inspection.ValueEqual("etat", models.EtatEnCours)
	inspection.ValueEqual("date", "2018-09-01")
	exploitants := etablissement.Value("exploitants").Array()
	exploitants.Length().Equal(1)
}
func TestGetEtablissementByIdNotOwnedByExploitant(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthExploitant(e.GET("/etablissements/{id}")).WithPath("id", "3").
		Expect().
		Status(http.StatusNotFound)
}
