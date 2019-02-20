package templates

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/stretchr/testify/require"
)

func TestRenderEmailTemplate(t *testing.T) {
	assert := require.New(t)

	service, err := New(Config{
		Dir: "templates/",
	})
	assert.NoError(err)

	data := domain.NouveauxMessagesUser{
		Destinataire: models.User{
			Email: "test@localhost",
			Nom:   "Michel Exploitant1",
		},
		Messages: []domain.NouveauMessage{
			domain.NouveauMessage{
				DateInspection:       "2018-02-24",
				RaisonEtablissement:  "Etablissement 1",
				SujetPointDeControle: "Rejets Eau",
				Message:              "Il faut des photos",
				NomAuteur:            "Alain Champion",
			},
			domain.NouveauMessage{
				DateInspection:       "2018-02-26",
				RaisonEtablissement:  "Etablissement 2",
				SujetPointDeControle: "Rejets Air",
				Message:              "Il faut des documents",
				NomAuteur:            "Alain Champion",
			},
		},
	}

	htmlPart, err := service.RenderHTMLEmailNouveauxMessages(data)
	assert.NoError(err)
	assert.Contains(htmlPart, "Il faut des photos")
	assert.Contains(htmlPart, "Il faut des documents")
}

func TestRenderLettreAnnonce(t *testing.T) {
	assert := require.New(t)

	service, err := New(Config{
		Dir: "templates/",
	})
	assert.NoError(err)

	inspection := &models.Inspection{
		Date: util.Date("2019-02-24"),
		Etablissement: &models.Etablissement{
			S3IC:       "1234",
			Nom:        "Nom 1",
			Raison:     "Raison sociale",
			Adresse1:   "1 rue des fleurs",
			Adresse2:   "",
			CodePostal: "75000",
			Commune:    "Paris",
			Exploitants: []models.User{
				models.User{
					Prenom:  "Michel",
					Nom:     "Exploitant1",
					Email:   "exploitant1@filharmonic.com",
					Profile: models.ProfilExploitant,
				},
			},
		},
		Inspecteurs: []models.User{
			models.User{
				Prenom:  "Alain",
				Nom:     "Champion",
				Email:   "inspecteur1@filharmonic.com",
				Profile: models.ProfilInspecteur,
			},
			models.User{
				Prenom:  "Corine",
				Nom:     "Dupont",
				Email:   "inspecteur2@filharmonic.com",
				Profile: models.ProfilInspecteur,
			},
			models.User{
				Prenom:  "Michel",
				Nom:     "Gérard",
				Email:   "inspecteur3@filharmonic.com",
				Profile: models.ProfilInspecteur,
			},
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Sujet: "Mesure des émissions atmosphériques canalisées par un organisme extérieur",
				ReferencesReglementaires: []string{
					"Article 3.2.3. de l'arrêté préfectoral du 28 juin 2017",
					"Article 3.2.8. de l'arrêté préfectoral du 28 juin 2017",
					"Article 8.2.1.2. de l'arrêté préfectoral du 28 juin 2017",
				},
			},
			models.PointDeControle{
				Sujet: "Autosurveillance des émissions canalisées de COV",
				ReferencesReglementaires: []string{
					"Article 8.2.1.1. de l'arrêté préfectoral du 28 juin 2017",
				},
			},
		},
	}

	_, err = service.RenderLettreAnnonce(domain.NewLettreAnnonce(inspection))
	assert.NoError(err)
}
