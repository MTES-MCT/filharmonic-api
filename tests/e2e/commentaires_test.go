package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAddCommentaire(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	commentaire := models.Commentaire{
		Message: "Commentaire général",
	}

	tests.AuthInspecteur(e.POST("/inspections/{id}/commentaires")).WithPath("id", 1).WithJSON(commentaire).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentaires := inspection.Value("commentaires").Array()
	commentaires.Length().Equal(3)
	commentaires.Last().Object().ValueEqual("message", "Commentaire général")
}

func TestAddCommentaireNotAllowed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	commentaire := models.Commentaire{
		Message: "Commentaire général",
	}

	tests.AuthUser(e.POST("/inspections/{id}/commentaires"), 4).WithPath("id", 1).WithJSON(commentaire).
		Expect().
		Status(http.StatusBadRequest)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentaires := inspection.Value("commentaires").Array()
	commentaires.Length().Equal(2)
}
