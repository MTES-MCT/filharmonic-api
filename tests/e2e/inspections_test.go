package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
	"github.com/MTES-MCT/filharmonic-api/util"
)

func TestListAllInspections(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	results := tests.AuthUser(e.GET("/inspections"), 4).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	results.Length().Equal(5)
	results.First().Object().ValueEqual("id", 1)
	results.First().Object().Value("etablissement").Object().ValueEqual("id", 1)
	results.Element(2).Object().ValueEqual("id", 3)
	results.Element(2).Object().Value("etablissement").Object().ValueEqual("id", 4)
}

func TestListInspectionsOwnedByInspecteur(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	results := tests.AuthUser(e.GET("/inspections"), 4).WithQuery("assigned", "true").
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	results.Length().Equal(2)
	results.First().Object().ValueEqual("id", 1)
	results.First().Object().Value("etablissement").Object().ValueEqual("id", 1)
	results.Last().Object().ValueEqual("id", 3)
	results.Last().Object().Value("etablissement").Object().ValueEqual("id", 4)
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
	etablissement := inspection.Value("etablissement").Object()
	etablissement.ValueEqual("id", 1)
	exploitants := etablissement.Value("exploitants").Array()
	exploitants.Length().Equal(1)
	exploitant := exploitants.First().Object()
	exploitant.ValueEqual("id", 1)
	exploitant.ValueEqual("email", "exploitant1@filharmonic.com")
	inspection.Value("themes").Array().Contains("Produits chimiques")
	inspecteurs := inspection.Value("inspecteurs").Array()
	inspecteurs.Length().Equal(2)
	inspecteurs.First().Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	inspecteurs.Last().Object().ValueEqual("email", "inspecteur2@filharmonic.com")
	suite := inspection.Value("suite").Object()
	suite.ValueEqual("type", "observation")
	suite.ValueEqual("synthese", "Observations à traiter")
	firstPointDeControle := inspection.Value("points_de_controle").Array().First().Object()
	firstPointDeControle.Value("references_reglementaires").Array().Contains("Article 3.2.3. de l'arrêté préfectoral du 28 juin 2017")
	firstPointDeControle.ValueEqual("sujet", "Mesure des émissions atmosphériques canalisées par un organisme extérieur")
	constat := firstPointDeControle.Value("constat").Object()
	constat.ValueEqual("type", "observation")
	constat.ValueEqual("remarques", "Il manque des choses")
	messages := firstPointDeControle.Value("messages").Array()
	messages.Length().Equal(4)
	firstMessage := messages.First().Object()
	firstMessage.ValueEqual("message", "Auriez-vous l'obligeance de me fournir le document approprié ?")
	firstMessage.Value("auteur").Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	firstMessage.Value("auteur").Object().NotContainsKey("password")
	firstCommentaire := inspection.Value("commentaires").Array().First().Object()
	firstCommentaire.ValueEqual("message", "Attention à l'article 243.")
	firstCommentaire.Value("auteur").Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	firstCommentaire.Value("auteur").Object().NotContainsKey("password")
	messageAvecPieceJointe := messages.Element(1).Object()
	piecesJointes := messageAvecPieceJointe.Value("pieces_jointes").Array().NotEmpty()
	firstPieceJointe := piecesJointes.First().Object()
	firstPieceJointe.ValueEqual("id", 1)
	firstPieceJointe.ValueEqual("nom", "photo-cuve.pdf")
	firstPieceJointe.ValueEqual("type", "application/pdf")
	firstPieceJointe.ValueEqual("taille", 7945)
}

func TestGetInspectionValideeAvecRapport(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", "5").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection.ValueEqual("etat", models.EtatTraitementNonConformites)
	rapport := inspection.Value("rapport").Object()
	rapport.ValueEqual("nom", "rapport.pdf")
}

func TestGetInspectionAsExploitantNotAllowed(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthExploitant(e.GET("/inspections/{id}")).WithPath("id", "2").
		Expect().
		Status(http.StatusForbidden).
		JSON().Object()
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
	messages.Length().Equal(3)
	firstMessage := messages.First().Object()
	firstMessage.ValueEqual("message", "Auriez-vous l'obligeance de me fournir le document approprié ?")
	firstMessage.Value("auteur").Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	firstMessage.Value("auteur").Object().NotContainsKey("password")
	lastMessage := messages.Last().Object()
	lastMessage.ValueEqual("message", "Il manque un document.")
	lastMessage.Value("auteur").Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	lastMessage.Value("auteur").Object().NotContainsKey("password")
	inspection.NotContainsKey("commentaires")
}

func TestCreateInspectionWithoutCanevas(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspectionInput := models.Inspection{
		Date:            util.Date("2019-01-22"),
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

func TestCreateInspectionWithCanevas(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspectionInput := models.Inspection{
		Date:            util.Date("2019-01-22"),
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
			"Santé",
		},
		Contexte:  "Contrôles de début d'année",
		CanevasId: 1,
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

	pointsDeControle := inspection.Value("points_de_controle").Array()
	pointsDeControle.Length().Equal(3)

	pointDeControle := pointsDeControle.First().Object()
	pointDeControle.ValueEqual("sujet", "COV")
	pointDeControle.Value("references_reglementaires").Array().Equal([]string{"Article 5 du 01/02/1995"})
	pointDeControle.NotContainsKey("messages")

	pointDeControle = pointsDeControle.Element(1).Object()
	pointDeControle.ValueEqual("sujet", "Rejets dans l'air")
	pointDeControle.Value("references_reglementaires").Array().Equal([]string{
		"Article 1 du 01/02/2003",
		"Article 2 du 01/02/2002",
	})
	messages := pointDeControle.Value("messages").Array()
	messages.Length().Equal(1)
	firstMessage := messages.First().Object()
	firstMessage.ValueEqual("message", "Donnez-moi le document sur l'air.")
	firstMessage.Value("auteur").Object().ValueEqual("email", "inspecteur1@filharmonic.com")

	pointDeControle = pointsDeControle.Last().Object()
	pointDeControle.ValueEqual("sujet", "Rejets dans l'eau")
	pointDeControle.Value("references_reglementaires").Array().Contains("Article 13 du 01/02/2006")
	messages = pointDeControle.Value("messages").Array()
	messages.Length().Equal(1)
	firstMessage = messages.First().Object()
	firstMessage.ValueEqual("message", "Donnez-moi le document sur l'eau.")
	firstMessage.Value("auteur").Object().ValueEqual("email", "inspecteur1@filharmonic.com")
}

func TestCreateCanevasFromInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	canevasInput := models.Canevas{
		Nom: "zest",
	}

	canevasId := tests.AuthInspecteur(e.POST("/inspections/{id}/canevas")).
		WithPath("id", 1).
		WithJSON(canevasInput).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("id").Raw()

	canevasArray := tests.AuthInspecteur(e.GET("/canevas")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	canevasArray.Length().Equal(2)
	lastCanevas := canevasArray.Last().Object()
	lastCanevas.ValueEqual("id", canevasId)
	lastCanevas.ValueEqual("nom", canevasInput.Nom)

	pointsDeControle := lastCanevas.Value("data").Object().Value("points_de_controle").Array()
	pointsDeControle.Length().Equal(2)

	pointDeControle := pointsDeControle.First().Object()
	pointDeControle.ValueEqual("sujet", "Autosurveillance des émissions canalisées de COV")
	pointDeControle.Value("references_reglementaires").Array().Equal([]string{
		"Article 8.2.1.1. de l'arrêté préfectoral du 28 juin 2017",
	})
	pointDeControle.ValueEqual("message", "Merci de me fournir le document.")

	pointDeControle = pointsDeControle.Last().Object()
	pointDeControle.ValueEqual("sujet", "Mesure des émissions atmosphériques canalisées par un organisme extérieur")
	pointDeControle.Value("references_reglementaires").Array().Equal([]string{
		"Article 3.2.3. de l'arrêté préfectoral du 28 juin 2017",
		"Article 3.2.8. de l'arrêté préfectoral du 28 juin 2017",
		"Article 8.2.1.2. de l'arrêté préfectoral du 28 juin 2017",
	})
	pointDeControle.ValueEqual("message", "Auriez-vous l'obligeance de me fournir le document approprié ?")
}

func TestUpdateCanevas(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	canevasInput := models.Canevas{
		Nom: "test",
	}

	canevasId := tests.AuthInspecteur(e.POST("/inspections/{id}/canevas")).
		WithPath("id", 1).
		WithJSON(canevasInput).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("id").Raw()

	canevasArray := tests.AuthInspecteur(e.GET("/canevas")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	canevasArray.Length().Equal(1)
	lastCanevas := canevasArray.Last().Object()
	lastCanevas.ValueEqual("id", canevasId)
	lastCanevas.ValueEqual("nom", canevasInput.Nom)

	pointsDeControle := lastCanevas.Value("data").Object().Value("points_de_controle").Array()
	pointsDeControle.Length().Equal(2)

	pointDeControle := pointsDeControle.First().Object()
	pointDeControle.ValueEqual("sujet", "Autosurveillance des émissions canalisées de COV")
	pointDeControle.Value("references_reglementaires").Array().Equal([]string{
		"Article 8.2.1.1. de l'arrêté préfectoral du 28 juin 2017",
	})
	pointDeControle.ValueEqual("message", "Merci de me fournir le document.")

	pointDeControle = pointsDeControle.Last().Object()
	pointDeControle.ValueEqual("sujet", "Mesure des émissions atmosphériques canalisées par un organisme extérieur")
	pointDeControle.Value("references_reglementaires").Array().Equal([]string{
		"Article 3.2.3. de l'arrêté préfectoral du 28 juin 2017",
		"Article 3.2.8. de l'arrêté préfectoral du 28 juin 2017",
		"Article 8.2.1.2. de l'arrêté préfectoral du 28 juin 2017",
	})
	pointDeControle.ValueEqual("message", "Auriez-vous l'obligeance de me fournir le document approprié ?")
}

func TestUpdateInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspectionInput := models.Inspection{
		Id:      1,
		Date:    util.Date("2019-01-30"),
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
			"Produits chimiques",
			"Incendie",
			"COV",
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
	themes.Length().Equal(3)
	themes.Equal([]string{"Produits chimiques", "Incendie", "COV"})
	inspecteurs := inspection.Value("inspecteurs").Array()
	inspecteurs.Length().Equal(2)
	inspecteurs.First().Object().ValueEqual("id", 3)
	inspecteurs.First().Object().ValueEqual("email", "inspecteur1@filharmonic.com")
	inspecteurs.Last().Object().ValueEqual("id", 5)
	inspecteurs.Last().Object().ValueEqual("email", "inspecteur3@filharmonic.com")
}

func TestGetInspectionAsExploitantPointDeControleNonPublie(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspection := tests.AuthExploitant(e.GET("/inspections/{id}")).WithPath("id", "1").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection.Value("points_de_controle").Array().Length().Equal(1)
}

func TestValidateInspectionSansNonConformites(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspectionInput := models.Inspection{
		MotifRejetValidation: "pas bon",
	}

	tests.AuthApprobateur(e.POST("/inspections/{id}/rejeter")).WithPath("id", 3).WithJSON(inspectionInput).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	suite := models.Suite{
		Type: models.TypeSuiteAucune,
	}
	tests.AuthInspecteur(e.PUT("/inspections/{id}/suite")).WithPath("id", 3).WithJSON(suite).
		Expect().
		Status(http.StatusOK)

	tests.AuthInspecteur(e.POST("/inspections/{id}/demandervalidation")).WithPath("id", 3).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	tests.AuthApprobateur(e.POST("/inspections/{id}/valider")).WithPath("id", 3).
		WithMultipart().WithFile("file", "../testdata/pdf-sample.pdf").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection := tests.AuthApprobateur(e.GET("/inspections/{id}")).WithPath("id", 3).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection.ValueEqual("etat", models.EtatClos)
}

func TestValidateInspectionAvecNonConformites(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthApprobateur(e.POST("/inspections/{id}/valider")).WithPath("id", 3).
		WithMultipart().WithFile("file", "../testdata/pdf-sample.pdf").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection := tests.AuthApprobateur(e.GET("/inspections/{id}")).WithPath("id", 3).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection.ValueEqual("etat", models.EtatTraitementNonConformites)
}

func TestPublishInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.POST("/inspections/{id}/publier")).WithPath("id", 2).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 2).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection.ValueEqual("etat", models.EtatEnCours)
}
func TestAskValidateInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	constat := models.Constat{
		Type: models.TypeConstatConforme,
	}
	tests.AuthInspecteur(e.POST("/pointsdecontrole/{id}/constat")).WithPath("id", 6).WithJSON(constat).
		Expect().
		Status(http.StatusOK)

	suite := models.Suite{
		Type: models.TypeSuiteObservation,
	}
	tests.AuthInspecteur(e.POST("/inspections/{id}/suite")).WithPath("id", 4).WithJSON(suite).
		Expect().
		Status(http.StatusOK)

	tests.AuthInspecteur(e.POST("/inspections/{id}/demandervalidation")).WithPath("id", 4).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 4).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection.ValueEqual("etat", models.EtatAttenteValidation)
}
func TestRejectInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	inspectionInput := models.Inspection{
		MotifRejetValidation: "Contrôle pas assez poussé",
	}
	tests.AuthApprobateur(e.POST("/inspections/{id}/rejeter")).WithPath("id", 3).WithJSON(inspectionInput).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection := tests.AuthApprobateur(e.GET("/inspections/{id}")).WithPath("id", 3).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection.ValueEqual("etat", models.EtatEnCours)
	inspection.ValueEqual("validation_rejetee", true)
	inspection.ValueEqual("motif_rejet_validation", inspectionInput.MotifRejetValidation)

	inspection = tests.AuthUser(e.GET("/inspections/{id}"), 2).WithPath("id", 3).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection.ValueEqual("etat", models.EtatEnCours)
	inspection.NotContainsKey("validation_rejetee")
	inspection.NotContainsKey("motif_rejet_validation")

	inspections := tests.AuthUser(e.GET("/inspections"), 2).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	inspection = inspections.First().Object()
	inspection.ValueEqual("etat", models.EtatEnCours)
	inspection.NotContainsKey("validation_rejetee")
	inspection.NotContainsKey("motif_rejet_validation")

	etablissement := tests.AuthUser(e.GET("/etablissements/{id}"), 2).WithPath("id", 4).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection = etablissement.Value("inspections").Array().First().Object()
	inspection.ValueEqual("etat", models.EtatEnCours)
	inspection.NotContainsKey("validation_rejetee")
	inspection.NotContainsKey("motif_rejet_validation")
}

func TestCloreInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.POST("/pointsdecontrole/{id}/constat/resoudre")).WithPath("id", 7).
		Expect().
		Status(http.StatusOK)

	tests.AuthInspecteur(e.POST("/inspections/{id}/clore")).WithPath("id", 5).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection := tests.AuthInspecteur(e.GET("/inspections/{id}")).WithPath("id", 5).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	inspection.ValueEqual("etat", models.EtatClos)
}

func TestAddFavoriToInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.POST("/inspections/{id}/favori")).WithPath("id", 2).
		Expect().
		Status(http.StatusOK)

	favoris := tests.AuthInspecteur(e.GET("/inspectionsfavorites")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	favoris.Length().Equal(2)
	favori := favoris.Last().Object()
	favori.NotContainsKey("inspecteurs")
	favori.NotContainsKey("commentaires")
	favori.NotContainsKey("points_de_controle")
	favori.NotContainsKey("suite")
	favori.ValueEqual("id", 2)
	favori.ValueEqual("date", "2018-11-15")
	etablissement := favori.Value("etablissement").Object()
	etablissement.ValueEqual("adresse1", "1 rue des cordeliers")
	etablissement.ValueEqual("nom", "Nom 3")
}
func TestRemoveFavoriToInspection(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.DELETE("/inspections/{id}/favori")).WithPath("id", 1).
		Expect().
		Status(http.StatusOK)

	favoris := tests.AuthInspecteur(e.GET("/inspectionsfavorites")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	favoris.Length().Equal(0)
}

func TestListInspectionsFavorites(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	favoris := tests.AuthInspecteur(e.GET("/inspectionsfavorites")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	favoris.Length().Equal(1)
	favori := favoris.First().Object()
	favori.NotContainsKey("inspecteurs")
	favori.NotContainsKey("commentaires")
	favori.NotContainsKey("points_de_controle")
	favori.NotContainsKey("suite")
	favori.ValueEqual("id", 1)
	favori.ValueEqual("date", "2018-09-01")
	etablissement := favori.Value("etablissement").Object()
	etablissement.ValueEqual("adresse1", "1 rue des fleurs")
	etablissement.ValueEqual("nom", "Nom 1")
}
