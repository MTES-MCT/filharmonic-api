package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestGenererLettreAnnonce(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	lettreBody := tests.AuthInspecteur(e.GET("/inspections/{id}/lettreannonce")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).Body()

	lettreBody.NotEmpty()
	lettreBody.Contains("Mesure des émissions atmosphériques canalisées par un organisme extérieur")
}

func TestGenererLettreSuite(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	lettreBody := tests.AuthInspecteur(e.GET("/inspections/{id}/lettresuite")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).Body()

	lettreBody.NotEmpty()
}

func TestGenererRapport(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	lettreBody := tests.AuthInspecteur(e.GET("/inspections/{id}/rapport")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).Body()

	lettreBody.NotEmpty()
	lettreBody.Contains("Mesure des émissions atmosphériques canalisées par un organisme extérieur")
}
