package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestListCanevas(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	canevas := tests.AuthInspecteur(e.GET("/canevas")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	canevas.Length().Equal(1)
	firstCanevas := canevas.First().Object()
	firstCanevas.ValueEqual("id", 1)
	firstCanevas.ValueEqual("nom", "test")
}

func TestDeleteCanevas(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.DELETE("/canevas/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK)

	canevas := tests.AuthInspecteur(e.GET("/canevas")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	canevas.Length().Equal(0)
}
