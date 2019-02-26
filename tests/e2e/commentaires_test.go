package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAddCommentaireAsInspecteur(t *testing.T) {
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

func TestAddCommentaireAsApprobateur(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	commentaire := models.Commentaire{
		Message: "Commentaire général",
	}

	tests.AuthApprobateur(e.POST("/inspections/{id}/commentaires")).WithPath("id", 1).WithJSON(commentaire).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthApprobateur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentaires := inspection.Value("commentaires").Array()
	commentaires.Length().Equal(3)
	commentaires.Last().Object().ValueEqual("message", "Commentaire général")
}

func TestAddCommentaireWithPieceJointe(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	commentaire := models.Commentaire{
		Message: "Commentaire général",
		PiecesJointes: []models.PieceJointe{
			models.PieceJointe{
				Id: 4,
			},
		},
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
	lastCommentaire := commentaires.Last().Object()
	lastCommentaire.ValueEqual("message", "Commentaire général")
	piecesJointes := lastCommentaire.Value("pieces_jointes").Array().NotEmpty()
	firstPieceJointe := piecesJointes.First().Object()
	firstPieceJointe.ValueEqual("id", 4)
	firstPieceJointe.ValueEqual("nom", "article-loi-v2.pdf")
	firstPieceJointe.ValueEqual("type", "application/pdf")
	firstPieceJointe.ValueEqual("taille", 10000000)
}

func TestAddCommentaireNotAllowed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	commentaire := models.Commentaire{
		Message: "Commentaire général",
	}

	tests.AuthUser(e.POST("/inspections/{id}/commentaires"), 5).WithPath("id", 1).WithJSON(commentaire).
		Expect().
		Status(http.StatusForbidden)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentaires := inspection.Value("commentaires").Array()
	commentaires.Length().Equal(2)
}
