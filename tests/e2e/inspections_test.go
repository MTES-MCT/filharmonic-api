package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestListInspectionsOwnedByInspecteur(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	results := tests.AuthInspecteur(e.GET("/inspections")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	results.Length().Equal(2)
	results.First().Object().ValueEqual("id", 1)
	results.First().Object().Value("etablissement").Object().ValueEqual("id", 1)
	results.Last().Object().ValueEqual("id", 2)
	results.Last().Object().Value("etablissement").Object().ValueEqual("id", 3)
}

func TestListInspectionsOwnedByExploitant(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	results := tests.AuthExploitant(e.GET("/inspections")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	results.Length().Equal(1)
	results.First().Object().ValueEqual("id", 1)
	results.First().Object().Value("etablissement").Object().ValueEqual("id", 1)
}
