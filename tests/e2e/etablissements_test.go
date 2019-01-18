package e2e

import (
	"net/http"
	"testing"

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

func TestGetEtablissementsByIdOwnedByExploitant(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthExploitant(e.GET("/etablissements/{id}")).WithPath("id", "1").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ValueEqual("id", 1)
}
func TestGetEtablissementsByIdNotOwnedByExploitant(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthExploitant(e.GET("/etablissements/{id}")).WithPath("id", "3").
		Expect().
		Status(http.StatusNotFound)
}
