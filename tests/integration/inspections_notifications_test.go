package integration

import (
	"strings"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/go-pg/pg/types"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCreateInspectionHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspecteur1 := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	inspecteur2 := models.User{
		Email:   "inspecteur2@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	assert.NoError(application.DB.Insert(&inspecteur1, &inspecteur2))

	etablissement := models.Etablissement{
		Nom: "Équipement de pression",
	}
	assert.NoError(application.DB.Insert(&etablissement))

	inspection := models.Inspection{
		Date: util.Date("2019-02-08"),
		Type: models.TypePonctuel,
		Themes: []string{
			"Sanitaire",
		},
		Annonce:         true,
		Origine:         models.OriginePlanControle,
		Etat:            models.EtatEnCours,
		Contexte:        "Inspection en cours",
		EtablissementId: etablissement.Id,
		Inspecteurs: []models.User{
			inspecteur1,
			inspecteur2,
		},
	}

	ctx1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctx2 := &domain.UserContext{
		User: &inspecteur2,
	}

	idInspection, err := application.Repo.CreateInspection(ctx1, inspection)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementCreationInspection, notification.Evenement.Type)
	assert.Equal(idInspection, notification.Evenement.InspectionId)
	assert.Equal(inspecteur1.Id, notification.Evenement.AuteurId)
}

func TestUpdateInspectionHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspecteur1 := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	inspecteur2 := models.User{
		Email:   "inspecteur2@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	assert.NoError(application.DB.Insert(&inspecteur1, &inspecteur2))

	inspection := models.Inspection{
		Date: util.Date("2019-02-08"),
		Type: models.TypePonctuel,
		Themes: []string{
			"Sanitaire",
		},
		Annonce:  true,
		Origine:  models.OriginePlanControle,
		Etat:     models.EtatEnCours,
		Contexte: "Inspection en cours",
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
		},
		Inspecteurs: []models.User{
			inspecteur1,
			inspecteur2,
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctx1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctx2 := &domain.UserContext{
		User: &inspecteur2,
	}

	inspection.Contexte = "test"

	err := application.Repo.UpdateInspection(ctx1, inspection)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementModificationInspection, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur1.Id, notification.Evenement.AuteurId)
}

func TestPublishInspectionHasCreatedNotification(t *testing.T) {
	assert, application, close := tests.InitServiceEmptyDB(t)
	defer close()

	inspecteur1 := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	inspecteur2 := models.User{
		Email:   "inspecteur2@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	assert.NoError(application.DB.Insert(&inspecteur1, &inspecteur2))

	inspection := models.Inspection{
		Date: util.Date("2019-02-08"),
		Etat: models.EtatPreparation,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
		},
		Inspecteurs: []models.User{
			inspecteur1,
			inspecteur2,
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctxInspecteur1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctxInspecteur2 := &domain.UserContext{
		User: &inspecteur2,
	}
	err := application.Service.PublishInspection(ctxInspecteur1, inspection.Id)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctxInspecteur2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementPublicationInspection, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur1.Id, notification.Evenement.AuteurId)
}

func TestAskValidateInspectionHasCreatedNotification(t *testing.T) {
	assert, application, close := tests.InitServiceEmptyDB(t)
	defer close()

	inspecteur1 := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	inspecteur2 := models.User{
		Email:   "inspecteur2@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	approbateur := models.User{
		Email:   "approbateur1@filharmonic.com",
		Profile: models.ProfilApprobateur,
	}
	assert.NoError(application.DB.Insert(&inspecteur1, &inspecteur2, &approbateur))

	inspection := models.Inspection{
		Date: util.Date("2019-02-08"),
		Etat: models.EtatEnCours,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
		},
		Inspecteurs: []models.User{
			inspecteur1,
			inspecteur2,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Constat: &models.Constat{
					Type: models.TypeConstatConforme,
				},
			},
		},
		Suite: &models.Suite{
			Type: models.TypeSuiteAucune,
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctxInspecteur1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctxInspecteur2 := &domain.UserContext{
		User: &inspecteur2,
	}
	ctxApprobateur := &domain.UserContext{
		User: &approbateur,
	}
	err := application.Service.AskValidateInspection(ctxInspecteur1, inspection.Id)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctxInspecteur2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementDemandeValidationInspection, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur1.Id, notification.Evenement.AuteurId)

	notifications, err = application.Repo.ListNotifications(ctxApprobateur, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification = notifications[0]
	assert.Equal(models.EvenementDemandeValidationInspection, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur1.Id, notification.Evenement.AuteurId)
}

func TestRejectInspectionHasCreatedNotification(t *testing.T) {
	assert, application, close := tests.InitServiceEmptyDB(t)
	defer close()

	inspecteur1 := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	inspecteur2 := models.User{
		Email:   "inspecteur2@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	approbateur := models.User{
		Email:   "approbateur1@filharmonic.com",
		Profile: models.ProfilApprobateur,
	}
	assert.NoError(application.DB.Insert(&inspecteur1, &inspecteur2, &approbateur))

	inspection := models.Inspection{
		Date: util.Date("2019-02-08"),
		Etat: models.EtatAttenteValidation,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
		},
		Inspecteurs: []models.User{
			inspecteur1,
			inspecteur2,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Constat: &models.Constat{
					Type: models.TypeConstatConforme,
				},
			},
		},
		Suite: &models.Suite{
			Type: models.TypeSuiteAucune,
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctxInspecteur1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctxInspecteur2 := &domain.UserContext{
		User: &inspecteur2,
	}
	ctxApprobateur := &domain.UserContext{
		User: &approbateur,
	}
	err := application.Service.RejectInspection(ctxApprobateur, inspection.Id, "motif de rejet")
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctxInspecteur1, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementRejetValidationInspection, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(approbateur.Id, notification.Evenement.AuteurId)

	notifications, err = application.Repo.ListNotifications(ctxInspecteur2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification = notifications[0]
	assert.Equal(models.EvenementRejetValidationInspection, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(approbateur.Id, notification.Evenement.AuteurId)
}

func TestValidateInspectionHasCreatedNotification(t *testing.T) {
	assert, application, close := tests.InitServiceEmptyDB(t)
	defer close()

	inspecteur1 := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	inspecteur2 := models.User{
		Email:   "inspecteur2@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	approbateur := models.User{
		Email:   "approbateur1@filharmonic.com",
		Profile: models.ProfilApprobateur,
	}
	assert.NoError(application.DB.Insert(&inspecteur1, &inspecteur2, &approbateur))

	inspection := models.Inspection{
		Date: util.Date("2019-02-08"),
		Etat: models.EtatAttenteValidation,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
		},
		Inspecteurs: []models.User{
			inspecteur1,
			inspecteur2,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Constat: &models.Constat{
					Type: models.TypeConstatConforme,
				},
			},
		},
		Suite: &models.Suite{
			Type: models.TypeSuiteAucune,
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctxInspecteur1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctxInspecteur2 := &domain.UserContext{
		User: &inspecteur2,
	}
	ctxApprobateur := &domain.UserContext{
		User: &approbateur,
	}
	rapportFile := models.File{
		Content: strings.NewReader("MonContenu"),
		Type:    "application/pdf",
		Taille:  int64(len("MonContenu")),
		Nom:     "test.pdf",
	}
	err := application.Service.ValidateInspection(ctxApprobateur, inspection.Id, rapportFile)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctxInspecteur1, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementValidationInspection, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(approbateur.Id, notification.Evenement.AuteurId)

	notifications, err = application.Repo.ListNotifications(ctxInspecteur2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification = notifications[0]
	assert.Equal(models.EvenementValidationInspection, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(approbateur.Id, notification.Evenement.AuteurId)
}

func TestCloreInspectionHasCreatedNotification(t *testing.T) {
	assert, application, close := tests.InitServiceEmptyDB(t)
	defer close()

	inspecteur1 := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	inspecteur2 := models.User{
		Email:   "inspecteur2@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	assert.NoError(application.DB.Insert(&inspecteur1, &inspecteur2))

	inspection := models.Inspection{
		Date: util.Date("2019-02-08"),
		Etat: models.EtatTraitementNonConformites,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
		},
		Inspecteurs: []models.User{
			inspecteur1,
			inspecteur2,
		},
		DateValidation: types.NullTime{Time: util.Date("2019-03-08").Time},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Constat: &models.Constat{
					Type:               models.TypeConstatNonConforme,
					EcheanceResolution: util.Date("2019-04-08"),
					DateResolution: types.NullTime{
						Time: util.Date("2019-03-15").Time,
					},
				},
			},
		},
		Suite: &models.Suite{
			Type: models.TypeSuitePropositionRenforcement,
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctxInspecteur1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctxInspecteur2 := &domain.UserContext{
		User: &inspecteur2,
	}
	err := application.Service.CloreInspection(ctxInspecteur1, inspection.Id)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctxInspecteur2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementClotureInspection, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur1.Id, notification.Evenement.AuteurId)
}
