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
	themes.First().Object().ValueEqual("nom", "COV")
	themes.Last().Object().ValueEqual("nom", "Sûreté")
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
	themes.First().Object().ValueEqual("nom", "COV")
	themes.Element(5).Object().ValueEqual("nom", "Sanitaire")
	themes.Last().Object().ValueEqual("nom", "Sûreté")
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
