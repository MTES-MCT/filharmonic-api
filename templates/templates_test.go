package templates

import (
	"io/ioutil"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/stretchr/testify/require"
)

func initTest(t *testing.T) (*require.Assertions, *TemplateService) {
	assert := require.New(t)
	service, err := New(Config{
		Dir:     "./",
		BaseURL: "http://localhost:8080",
	})
	assert.NoError(err)
	return assert, service
}
func TestRenderEmailNouveauxMessages(t *testing.T) {
	assert, service := initTest(t)

	data := domain.NouveauxMessagesUser{
		Destinataire: models.User{
			Email: "test@localhost",
			Nom:   "Michel Exploitant1",
		},
		Messages: []domain.NouveauMessage{
			domain.NouveauMessage{
				DateInspection:       "24/02/2018",
				RaisonEtablissement:  "Etablissement 1",
				SujetPointDeControle: "Rejets Eau",
				Message:              "Il faut des photos",
				NomAuteur:            "Alain Champion",
				DateMessage:          "02/01/2018 à 15h04",
				InspectionId:         3,
				PointDeControleId:    5,
				MessageId:            7,
			},
			domain.NouveauMessage{
				DateInspection:       "26/02/2018",
				RaisonEtablissement:  "Etablissement 2",
				SujetPointDeControle: "Rejets Air",
				Message:              "Il faut des documents",
				DateMessage:          "16/01/2018 à 18h01",
				NomAuteur:            "Alain Champion",
				InspectionId:         1,
				PointDeControleId:    2,
				MessageId:            2,
			},
		},
	}

	result, err := service.RenderEmailNouveauxMessages(data)
	assert.NoError(err)
	assert.Contains(result.HTML, "Il faut des photos")
	assert.Contains(result.HTML, "Il faut des documents")
	assert.NoError(ioutil.WriteFile("../.tmp/email-new-messages.html", []byte(result.HTML), 0644))
	assert.NoError(ioutil.WriteFile("../.tmp/email-new-messages.txt", []byte(result.Text), 0644))
}

func TestRenderEmailRecapValidation(t *testing.T) {
	assert, service := initTest(t)

	data := domain.RecapValidationInspection{
		Destinataire: models.User{
			Email: "test@localhost",
			Nom:   "Michel Exploitant1",
		},
		InspectionId:         3,
		NonConformites:       true,
		DateInspection:       "01/01/2019",
		RaisonEtablissement:  "Etablissement 1",
		AdresseEtablissement: "1 rue des Fleurs 75000 Paris",
	}

	result, err := service.RenderEmailRecapValidation(data)
	assert.NoError(err)
	assert.Contains(result.HTML, data.RaisonEtablissement)
	assert.Contains(result.HTML, "échéances de résolution")
	assert.NoError(ioutil.WriteFile("../.tmp/email-recap-validation.html", []byte(result.HTML), 0644))
	assert.NoError(ioutil.WriteFile("../.tmp/email-recap-validation.txt", []byte(result.Text), 0644))
}

func TestRenderHTMLEmailExpirationDelais(t *testing.T) {
	assert, service := initTest(t)

	data := domain.InspectionExpirationDelais{
		Destinataire: models.User{
			Email: "test@localhost",
			Nom:   "Michel Exploitant1",
		},
		InspectionId:         3,
		DateInspection:       "01/01/2019",
		RaisonEtablissement:  "Etablissement 1",
		AdresseEtablissement: "1 rue des Fleurs 75000 Paris",
	}

	result, err := service.RenderEmailExpirationDelais(data)
	assert.NoError(err)
	assert.Contains(result.HTML, "lever vos non-conformités")
	assert.Contains(result.HTML, data.AdresseEtablissement)
	assert.NoError(ioutil.WriteFile("../.tmp/email-expiration-delais.html", []byte(result.HTML), 0644))
	assert.NoError(ioutil.WriteFile("../.tmp/email-expiration-delais.txt", []byte(result.Text), 0644))
}

func TestRenderHTMLEmailRappelEcheances(t *testing.T) {
	assert, service := initTest(t)

	data := domain.InspectionEcheancesProches{
		Destinataire: models.User{
			Email: "test@localhost",
			Nom:   "Michel Exploitant1",
		},
		InspectionId:         3,
		DateInspection:       "01/01/2019",
		RaisonEtablissement:  "Etablissement 1",
		AdresseEtablissement: "1 rue des Fleurs 75000 Paris",
	}

	result, err := service.RenderEmailRappelEcheances(data)
	assert.NoError(err)
	assert.Contains(result.HTML, "échéances de résolution sont proches")
	assert.Contains(result.Text, "échéances de résolution sont proches")
	assert.Contains(result.HTML, data.AdresseEtablissement)
	assert.Contains(result.Text, data.AdresseEtablissement)
	assert.NoError(ioutil.WriteFile("../.tmp/email-rappel-echeances.html", []byte(result.HTML), 0644))
	assert.NoError(ioutil.WriteFile("../.tmp/email-rappel-echeances.txt", []byte(result.Text), 0644))
}

func TestRenderLettre(t *testing.T) {
	assert, service := initTest(t)

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

	result, err := service.RenderODTLettreAnnonce(domain.NewLettre(inspection))
	assert.NoError(err)
	assert.NoError(ioutil.WriteFile("../.tmp/lettre-annonce.fodt", []byte(result.Text), 0644))
}

func TestRenderRapport(t *testing.T) {
	assert, service := initTest(t)

	inspection := &models.Inspection{
		Date:                 util.Date("2019-02-24"),
		Themes:               []string{"Rejets dans l'eau", "Rejets dans l'air"},
		Type:                 models.TypeApprofondi,
		Origine:              models.OrigineCirconstancielle,
		Annonce:              true,
		Circonstance:         models.CirconstanceAutre,
		PersonnesRencontrees: "Gérard Bichon, Françoise Denivelle et Brigitte Lombard",
		DetailCirconstance:   "Explosion du 02/10/1978. Vents forts",
		Contexte: `Sed (saepe enim redeo ad Scipionem, cuius omnis sermo erat de amicitia) querebatur, quod omnibus in rebus homines diligentiores essent; capras et oves quot quisque haberet, dicere posse, amicos quot haberet, non posse dicere et in illis quidem parandis adhibere curam, in amicis eligendis neglegentis esse nec habere quasi signa quaedam et notas, quibus eos qui ad amicitias essent idonei, iudicarent. Sunt igitur firmi et stabiles et constantes eligendi; cuius generis est magna penuria. Et iudicare difficile est sane nisi expertum; experiendum autem est in ipsa amicitia. Ita praecurrit amicitia iudicium tollitque experiendi potestatem.
		Quis enim aut eum diligat quem metuat, aut eum a quo se metui putet? Coluntur tamen simulatione dumtaxat ad tempus. Quod si forte, ut fit plerumque, ceciderunt, tum intellegitur quam fuerint inopes amicorum. Quod Tarquinium dixisse ferunt, tum exsulantem se intellexisse quos fidos amicos habuisset, quos infidos, cum iam neutris gratiam referre posset.
		Mensarum enim voragines et varias voluptatum inlecebras, ne longius progrediar, praetermitto illuc transiturus quod quidam per ampla spatia urbis subversasque silices sine periculi metu properantes equos velut publicos signatis quod dicitur calceis agitant, familiarium agmina tamquam praedatorios globos post terga trahentes ne Sannione quidem, ut ait comicus, domi relicto. quos imitatae matronae complures opertis capitibus et basternis per latera civitatis cuncta discurrunt.`,
		Etablissement: &models.Etablissement{
			S3IC:       "1234",
			Nom:        "Nom 1",
			Raison:     "Raison sociale",
			Adresse1:   "1 rue des fleurs",
			Adresse2:   "",
			CodePostal: "75000",
			Commune:    "Paris",
			Activite:   "Abattoirs",
			Regime:     models.RegimeAutorisation,
			Seveso:     "haut",
			Iedmtd:     false,
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
				Messages: []models.Message{
					models.Message{
						Auteur: &models.User{
							Prenom:  "Alain",
							Nom:     "Champion",
							Email:   "inspecteur1@filharmonic.com",
							Profile: models.ProfilInspecteur,
						},
						Message: "Merci de me donner le document.",
						Date:    util.DateTime("2018-02-15T13:57:00"),
					},
					models.Message{
						Auteur: &models.User{
							Prenom:  "Michel",
							Nom:     "Exploitant1",
							Email:   "exploitant1@filharmonic.com",
							Profile: models.ProfilExploitant,
						},
						Message: "Voici le document demandé.",
						Date:    util.DateTime("2018-02-25T13:57:00"),
						PiecesJointes: []models.PieceJointe{
							models.PieceJointe{
								Nom: "rapport-2018.pdf",
							},
							models.PieceJointe{
								Nom: "rapport-2017.pdf",
							},
						},
					},
				},
				Constat: &models.Constat{
					Type:               models.TypeConstatObservation,
					DelaiNombre:        3,
					DelaiUnite:         "mois",
					EcheanceResolution: util.Date("2019-04-30"),
					Remarques:          "Vous devez réparer la cuve.",
				},
			},
			models.PointDeControle{
				Sujet: "Autosurveillance des émissions canalisées de COV",
				ReferencesReglementaires: []string{
					"Article 8.2.1.1. de l'arrêté préfectoral du 28 juin 2017",
				},
				Constat: &models.Constat{
					Type: models.TypeConstatConforme,
				},
			},
		},
		Suite: &models.Suite{
			Type:        models.TypeSuiteObservation,
			Synthese:    "Cette visite a permis de relever des non conformités vis-à-vis des prescriptions examinées, ainsi que des points faisant l’objet d’observations. L’exploitant devra fournir selon les délais mentionnés dans le présent rapport, les éléments permettant de justifier de la mise en œuvre des actions correctives nécessaires pour les lever.",
			PenalEngage: true,
		},
	}

	result, err := service.RenderODTRapport(domain.NewRapport(inspection))
	assert.NoError(err)
	assert.NoError(ioutil.WriteFile("../.tmp/rapport.fodt", []byte(result.Text), 0644))
}

func TestRenderLettreSuite(t *testing.T) {
	assert, service := initTest(t)

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
		Suite: &models.Suite{
			PenalEngage: true,
			Type:        models.TypeSuiteAucune,
			Synthese:    "blah blah",
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

	result, err := service.RenderODTLettreSuite(domain.NewLettre(inspection))
	assert.NoError(err)
	assert.NoError(ioutil.WriteFile("../.tmp/lettre-suite.fodt", []byte(result.Text), 0644))
}
