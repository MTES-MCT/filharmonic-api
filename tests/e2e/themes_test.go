package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestListThemes(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	themes := tests.AuthInspecteur(e.GET("/themes")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	themes.Length().Equal(6)
	firstTheme := themes.First().Object()
	firstTheme.ValueEqual("id", 1)
	firstTheme.ValueEqual("nom", "Rejets dans l'eau")
}

func TestCreateTheme(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	themeInput := models.Theme{
		Nom: "Sanitaire",
	}

	tests.AuthInspecteur(e.POST("/themes")).WithJSON(themeInput).
		Expect().
		Status(http.StatusOK)

	themes := tests.AuthInspecteur(e.GET("/themes")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	themes.Length().Equal(7)
	lastTheme := themes.Last().Object()
	lastTheme.ValueEqual("id", 7)
	lastTheme.ValueEqual("nom", themeInput.Nom)
}

func TestDeleteTheme(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.DELETE("/themes/{id}")).WithPath("id", 6).
		Expect().
		Status(http.StatusOK)

	themes := tests.AuthInspecteur(e.GET("/themes")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	themes.Length().Equal(5)
}
