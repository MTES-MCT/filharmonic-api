package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestGetRapportInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthApprobateur(e.POST("/inspections/{id}/valider")).WithPath("id", 3).
		WithMultipart().WithFile("file", "../testdata/pdf-sample.pdf").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	tests.AuthInspecteur(e.GET("/inspections/{id}/rapport")).WithPath("id", 3).
		Expect().
		Status(http.StatusOK).ContentType("application/octet-stream")

	tests.AuthUser(e.GET("/inspections/{id}/rapport"), 2).WithPath("id", 3).
		Expect().
		Status(http.StatusOK).ContentType("application/octet-stream")
}
