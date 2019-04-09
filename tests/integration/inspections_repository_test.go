package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/go-pg/pg/types"
)

func TestListInpections(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)
	exploitant := models.User{
		Prenom:  "Michel",
		Nom:     "Exploitant1",
		Email:   "exploitant1@filharmonic.com",
		Profile: models.ProfilExploitant,
	}
	assert.NoError(application.DB.Insert(&exploitant))

	inspecteur := models.User{
		Prenom:  "Alain",
		Nom:     "Champion",
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	assert.NoError(application.DB.Insert(&inspecteur))

	inspectionAvecMessage := models.Inspection{
		Date: util.Date("2019-01-10"),
		Type: models.TypeApprofondi,
		Etat: models.EtatAttenteValidation,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
			Exploitants: []models.User{
				exploitant,
			},
		},
		Suite: &models.Suite{
			Type: models.TypeSuitePropositionMiseEnDemeure,
		},
		Inspecteurs: []models.User{
			inspecteur,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Messages: []models.Message{
					models.Message{
						Auteur: &exploitant,
					},
				},
			},
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspectionAvecMessage))

	inspectionSansMessage := models.Inspection{
		Date: util.Date("2019-01-10"),
		Type: models.TypeApprofondi,
		Etat: models.EtatAttenteValidation,
		Etablissement: &models.Etablissement{
			Nom: "Équipement 2",
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspectionSansMessage))

	ctxInspecteur := &domain.UserContext{
		User: &inspecteur,
	}
	inspections, err := application.Repo.ListInspections(ctxInspecteur, domain.ListInspectionsFilter{})
	assert.NoError(err)
	assert.Len(inspections, 2)
	assert.Equal(1, inspections[0].NbMessagesNonLus)
	assert.Equal(0, inspections[1].NbMessagesNonLus)

	ctxExploitant := &domain.UserContext{
		User: &exploitant,
	}
	inspections, err = application.Repo.ListInspections(ctxExploitant, domain.ListInspectionsFilter{})
	assert.NoError(err)
	assert.Len(inspections, 1)
	assert.Equal(0, inspections[0].NbMessagesNonLus)
}

func TestGetInspectionTypesConstatsSuiteByID(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspectionAvecConstat := models.Inspection{
		Date: util.Date("2019-01-10"),
		Type: models.TypeApprofondi,
		Etat: models.EtatAttenteValidation,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
		},
		Suite: &models.Suite{
			Type: models.TypeSuiteObservation,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Constat: &models.Constat{
					Type: models.TypeConstatObservation,
				},
			},
			models.PointDeControle{
				Publie: true,
				Sujet:  "test2",
				Constat: &models.Constat{
					Type: models.TypeConstatConforme,
				},
			},
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspectionAvecConstat))
	inspection, err := application.Repo.GetInspectionTypesConstatsSuiteByID(inspectionAvecConstat.Id)
	assert.NoError(err)
	assert.Equal(inspectionAvecConstat.Suite.Type, inspection.Suite.Type)
	assert.Equal(2, len(inspection.PointsDeControle))
	pointDeControle := inspection.PointsDeControle[0]
	assert.Equal(inspectionAvecConstat.PointsDeControle[0].Constat.Type, pointDeControle.Constat.Type)
}

func TestGetRecapsValidation(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspectionAvecRecaps := models.Inspection{
		Date: util.Date("2019-01-10"),
		Type: models.TypeApprofondi,
		Etat: models.EtatTraitementNonConformites,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
			Exploitants: []models.User{
				models.User{
					Prenom:  "Michel",
					Nom:     "Exploitant1",
					Email:   "exploitant1@filharmonic.com",
					Profile: models.ProfilExploitant,
				},
				models.User{
					Prenom:  "Bernard",
					Nom:     "Exploitant2",
					Email:   "exploitant2@filharmonic.com",
					Profile: models.ProfilExploitant,
				},
			},
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Constat: &models.Constat{
					Type: models.TypeConstatObservation,
				},
			},
			models.PointDeControle{
				Publie: true,
				Sujet:  "test2",
				Constat: &models.Constat{
					Type: models.TypeConstatNonConforme,
				},
			},
		},
		Suite: &models.Suite{
			Type: models.TypeSuitePropositionRenforcement,
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspectionAvecRecaps))

	recaps, err := application.Repo.GetRecapsValidation(inspectionAvecRecaps.Id)
	assert.NoError(err)
	assert.Len(recaps, 2)
}

func TestListInspectionsExpirationDelais(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspection1 := models.Inspection{
		Date: util.Date("2019-01-10"),
		Type: models.TypeApprofondi,
		Etat: models.EtatTraitementNonConformites,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
			Exploitants: []models.User{
				models.User{
					Prenom:  "Michel",
					Nom:     "Exploitant1",
					Email:   "exploitant1@filharmonic.com",
					Profile: models.ProfilExploitant,
				},
				models.User{
					Prenom:  "Bernard",
					Nom:     "Exploitant2",
					Email:   "exploitant2@filharmonic.com",
					Profile: models.ProfilExploitant,
				},
			},
		},
		Suite: &models.Suite{
			Type: models.TypeSuitePropositionMiseEnDemeure,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Constat: &models.Constat{
					Type:               models.TypeConstatNonConforme,
					EcheanceResolution: util.Date("2019-01-10"),
				},
			},
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection1))
	inspection2 := models.Inspection{
		Date: util.Date("2019-01-10"),
		Type: models.TypeApprofondi,
		Etat: models.EtatTraitementNonConformites,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
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
			Type: models.TypeSuitePropositionMiseEnDemeure,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Constat: &models.Constat{
					Type:               models.TypeConstatNonConforme,
					EcheanceResolution: util.Date("2019-01-10"),
				},
			},
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection2))

	util.SetTime(util.Date("2019-04-01").Time)
	destinatairesInspections, err := application.Repo.ListInspectionsExpirationDelais()
	assert.NoError(err)
	assert.Len(destinatairesInspections, 3)
}
func TestListInspectionsEcheancesProches(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspection1 := models.Inspection{
		Date: util.Date("2019-01-10"),
		Type: models.TypeApprofondi,
		Etat: models.EtatTraitementNonConformites,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
			Exploitants: []models.User{
				models.User{
					Prenom:  "Michel",
					Nom:     "Exploitant1",
					Email:   "exploitant1@filharmonic.com",
					Profile: models.ProfilExploitant,
				},
				models.User{
					Prenom:  "Bernard",
					Nom:     "Exploitant2",
					Email:   "exploitant2@filharmonic.com",
					Profile: models.ProfilExploitant,
				},
			},
		},
		DateValidation: types.NullTime{Time: util.Date("2019-04-01").Time},
		Suite: &models.Suite{
			Type: models.TypeSuitePropositionMiseEnDemeure,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Constat: &models.Constat{
					Type:               models.TypeConstatNonConforme,
					DelaiNombre:        30,
					DelaiUnite:         "jours",
					EcheanceResolution: util.Date("2019-04-30"),
				},
			},
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection1))

	util.SetTime(util.Date("2019-04-15").Time)
	inspections, err := application.Repo.ListInspectionsEcheancesProches(0.2)
	assert.NoError(err)
	assert.Len(inspections, 0)

	util.SetTime(util.Date("2019-04-26").Time)
	inspections, err = application.Repo.ListInspectionsEcheancesProches(0.2)
	assert.NoError(err)
	assert.Len(inspections, 2)
}

func TestListUsersAssignedToInspection(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspection := models.Inspection{
		Date: util.Date("2019-01-10"),
		Type: models.TypeApprofondi,
		Etat: models.EtatTraitementNonConformites,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
			Exploitants: []models.User{
				models.User{
					Email:   "exploitant1@filharmonic.com",
					Profile: models.ProfilExploitant,
				},
				models.User{
					Email:   "exploitant2@filharmonic.com",
					Profile: models.ProfilExploitant,
				},
			},
		},
		Inspecteurs: []models.User{
			models.User{
				Email:   "inspecteur1@filharmonic.com",
				Profile: models.ProfilInspecteur,
			},
			models.User{
				Email:   "inspecteur2@filharmonic.com",
				Profile: models.ProfilInspecteur,
			},
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	userIds, err := application.Repo.ListUsersAssignedToInspection(inspection.Id)
	assert.NoError(err)
	assert.Len(userIds, 4)
	assert.Equal([]int64{1, 2, 3, 4}, userIds)
}
