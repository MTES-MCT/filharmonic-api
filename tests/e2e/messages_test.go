package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAddMessage(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	message := models.Message{
		Message: "Message publique",
	}

	tests.AuthInspecteur(e.POST("/pointsdecontrole/{id}/messages")).WithPath("id", 1).WithJSON(message).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	firstPointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	messages := firstPointDeControle.Value("messages").Array()
	messages.Length().Equal(5)
	lastMessage := messages.Last().Object()
	lastMessage.ValueEqual("message", "Message publique")
	lastMessage.Value("auteur").Object().ValueEqual("email", "inspecteur1@filharmonic.com")
}

func TestAddMessageWithPieceJointe(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	message := models.Message{
		Message: "Message publique",
		PiecesJointes: []models.PieceJointe{
			models.PieceJointe{
				Id: 2,
			},
		},
	}

	tests.AuthExploitant(e.POST("/pointsdecontrole/{id}/messages")).WithPath("id", 1).WithJSON(message).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	firstPointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	messages := firstPointDeControle.Value("messages").Array()
	messages.Length().Equal(5)
	lastMessage := messages.Last().Object()
	lastMessage.ValueEqual("message", "Message publique")
	lastMessage.Value("auteur").Object().ValueEqual("email", "exploitant1@filharmonic.com")
	piecesJointes := lastMessage.Value("pieces_jointes").Array().NotEmpty()
	firstPieceJointe := piecesJointes.First().Object()
	firstPieceJointe.ValueEqual("id", 2)
	firstPieceJointe.ValueEqual("nom", "photo-cuve-2.pdf")
	firstPieceJointe.ValueEqual("type", "application/pdf")
	firstPieceJointe.ValueEqual("taille", 2262000)
}

func TestAddMessageWithPieceJointeBadOwner(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	message := models.Message{
		Message: "Message publique",
		PiecesJointes: []models.PieceJointe{
			models.PieceJointe{
				Id: 2,
			},
		},
	}

	tests.AuthUser(e.POST("/pointsdecontrole/{id}/messages"), 2).WithPath("id", 1).WithJSON(message).
		Expect().
		Status(http.StatusBadRequest)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	firstPointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	messages := firstPointDeControle.Value("messages").Array()
	messages.Length().Equal(4)
}

func TestAddMessageAsInspecteurNotAllowed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	message := models.Message{
		Message: "Message publique",
	}

	tests.AuthUser(e.POST("/pointsdecontrole/{id}/messages"), 4).WithPath("id", 1).WithJSON(message).
		Expect().
		Status(http.StatusBadRequest)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	firstPointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	messages := firstPointDeControle.Value("messages").Array()
	messages.Length().Equal(4)
}

func TestAddMessageAsExploitantNotAllowed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	message := models.Message{
		Message: "Message publique",
	}

	tests.AuthUser(e.POST("/pointsdecontrole/{id}/messages"), 2).WithPath("id", 1).WithJSON(message).
		Expect().
		Status(http.StatusBadRequest)
}

func TestLireMessageAsInspecteur(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.POST("/messages/{id}/lire")).WithPath("id", 7).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 2).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	pointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	messages := pointDeControle.Value("messages").Array()
	messages.Length().Equal(2)
	lastMessage := messages.Last().Object()
	lastMessage.ValueEqual("lu", true)
}

func TestLireMessageAsExploitant(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthExploitant(e.POST("/messages/{id}/lire")).WithPath("id", 6).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthExploitant(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	pointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	messages := pointDeControle.Value("messages").Array()
	messages.Length().Equal(3)
	lastMessage := messages.Last().Object()
	lastMessage.ValueEqual("lu", true)
}

func TestLireMessageAsInspecteurNotAllowed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.POST("/messages/{id}/lire")).WithPath("id", 4).
		Expect().
		Status(http.StatusBadRequest)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	pointDeControle := inspection.Value("points_de_controle").Array().Element(1).Object()
	messages := pointDeControle.Value("messages").Array()
	messages.Length().Equal(1)
	lastMessage := messages.Last().Object()
	lastMessage.ValueEqual("lu", false)
}
