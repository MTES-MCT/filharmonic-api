package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestGenererLettreAnnonce(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	lettreBody := tests.AuthInspecteur(e.GET("/inspections/{id}/generer/lettreannonce")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).Body()

	lettreBody.NotEmpty()
	lettreBody.Contains("Mesure des émissions atmosphériques canalisées par un organisme extérieur")
	lettreBody.NotContains("Autosurveillance des émissions canalisées de COV")
}

func TestGenererLettreSuite(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthApprobateur(e.POST("/inspections/{id}/valider")).WithPath("id", 3).
		WithMultipart().WithFile("file", "../testdata/pdf-sample.pdf").
		Expect().
		Status(http.StatusOK).Body()

	lettreBody := tests.AuthInspecteur(e.GET("/inspections/{id}/generer/lettresuite")).WithPath("id", 3).
		Expect().
		Status(http.StatusOK).Body()

	lettreBody.NotEmpty()
}

func TestGenererRapport(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	lettreBody := tests.AuthInspecteur(e.GET("/inspections/{id}/generer/rapport")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).Body()

	lettreBody.NotEmpty()
	lettreBody.Contains("Mesure des émissions atmosphériques canalisées par un organisme extérieur")
	lettreBody.NotContains("Autosurveillance des émissions canalisées de COV")
}
