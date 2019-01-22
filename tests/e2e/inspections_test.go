package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
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

func TestGetInspectionAsInspecteur(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", "1").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	inspection.ValueEqual("id", 1)
	inspection.ValueEqual("date", "2018-09-01")
	inspection.Value("etablissement").Object().ValueEqual("id", 1)
	inspection.Value("themes").Array().Contains("Produits chimiques")
	firstPointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	firstPointDeControle.Value("references_reglementaires").Array().Contains("Article 3.2.3. de l'arrêté préfectoral du 28 juin 2017")
	firstPointDeControle.ValueEqual("sujet", "Mesure des émissions atmosphériques canalisées par un organisme extérieur")
	messages := firstPointDeControle.Value("messages").Array()
	messages.Length().Equal(3)
	firstMessage := messages.First().Object()
	firstMessage.ValueEqual("message", "Auriez-vous l'obligeance de me fournir le document approprié ?")
	firstMessage.Value("auteur").Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	firstMessage.Value("auteur").Object().NotContainsKey("password")
	firstCommentaire := inspection.Value("commentaires").Array().First().Object()
	firstCommentaire.ValueEqual("message", "Attention à l'article 243.")
	firstCommentaire.Value("auteur").Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	firstCommentaire.Value("auteur").Object().NotContainsKey("password")
}

func TestGetInspectionAsExploitantNotAllowed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	result := tests.AuthExploitant(e.GET("/inspections/{id}")).WithPath("id", "2").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object()
	result.ValueEqual("message", "not_found")
}

func TestGetInspectionAsExploitantAllowed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspection := tests.AuthExploitant(e.GET("/inspections/{id}")).WithPath("id", "1").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	inspection.ValueEqual("id", 1)
	inspection.ValueEqual("date", "2018-09-01")
	inspection.Value("etablissement").Object().ValueEqual("id", 1)
	inspection.Value("themes").Array().Contains("Produits chimiques")
	firstPointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	firstPointDeControle.Value("references_reglementaires").Array().Contains("Article 3.2.3. de l'arrêté préfectoral du 28 juin 2017")
	firstPointDeControle.ValueEqual("sujet", "Mesure des émissions atmosphériques canalisées par un organisme extérieur")
	messages := firstPointDeControle.Value("messages").Array()
	messages.Length().Equal(2)
	firstMessage := messages.First().Object()
	firstMessage.ValueEqual("message", "Auriez-vous l'obligeance de me fournir le document approprié ?")
	firstMessage.Value("auteur").Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	firstMessage.Value("auteur").Object().NotContainsKey("password")
	lastMessage := messages.Last().Object()
	lastMessage.ValueEqual("message", "Voici le document.")
	lastMessage.Value("auteur").Object().ValueEqual("email", "exploitant1@filharmonic.com")
	lastMessage.Value("auteur").Object().NotContainsKey("password")
	inspection.NotContainsKey("commentaires")
}

func TestCreateInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspectionInput := models.Inspection{
		Date:            tests.Date("2019-01-22"),
		Type:            models.TypeCourant,
		Annonce:         true,
		Origine:         models.OriginePlanControle,
		EtablissementId: 1,
		Inspecteurs: []models.User{
			models.User{
				Id: 3,
			},
			models.User{
				Id: 4,
			},
		},
		Themes: []string{
			"Incendie",
			"Produits chimiques",
		},
		Contexte: "Contrôles de début d'année",
	}

	inspectionId := tests.AuthInspecteur(e.POST("/inspections")).WithJSON(inspectionInput).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("id").Raw()

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", inspectionId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	inspection.ValueEqual("id", inspectionId)
	inspection.ValueEqual("etat", models.EtatPreparation)
	inspecteurs := inspection.Value("inspecteurs").Array()
	inspecteurs.Length().Equal(2)
	inspecteurs.First().Object().ValueEqual("id", 3)
	inspecteurs.First().Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	inspecteurs.Last().Object().ValueEqual("id", 4)
	inspecteurs.Last().Object().ValueEqual("email", "inspecteur2@filharmonic.com")
}

func TestSaveInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspectionInput := models.Inspection{
		Id:      1,
		Date:    tests.Date("2019-01-30"),
		Type:    models.TypeCourant,
		Annonce: true,
		Origine: models.OriginePlanControle,
		Inspecteurs: []models.User{
			models.User{
				Id: 3,
			},
			models.User{
				Id: 5,
			},
		},
		Themes: []string{
			"Incendie",
			"Produits chimiques",
		},
		Contexte: "Contrôles de début d'année",
	}

	tests.AuthInspecteur(e.PUT("/inspections/{id}")).WithPath("id", 1).WithJSON(inspectionInput).
		Expect().
		Status(http.StatusOK)

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	inspection.ValueEqual("id", 1)
	inspection.ValueEqual("etat", models.EtatEnCours)
	inspection.ValueEqual("date", "2019-01-30")
	themes := inspection.Value("themes").Array()
	themes.Length().Equal(2)
	themes.Equal([]string{"Incendie", "Produits chimiques"})
	inspecteurs := inspection.Value("inspecteurs").Array()
	inspecteurs.Length().Equal(2)
	inspecteurs.First().Object().ValueEqual("id", 3)
	inspecteurs.First().Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	inspecteurs.Last().Object().ValueEqual("id", 5)
	inspecteurs.Last().Object().ValueEqual("email", "inspecteur3@filharmonic.com")
}
