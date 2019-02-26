package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/util"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCreateInspectionHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

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
		EtablissementId: 4,
		Inspecteurs: []models.User{
			models.User{
				Id: 3,
			},
			models.User{
				Id: 4,
			},
		},
	}

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}
	ctx2 := &domain.UserContext{
		User: &models.User{
			Id: 4,
		},
	}

	idInspection, err := application.Repo.CreateInspection(ctx, inspection)
	assert.NoError(err)
	assert.Equal(int64(6), idInspection)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementCreationInspection, notification.Evenement.Type)
	assert.Equal(idInspection, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
}

func TestUpdateInspectionHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}
	ctx2 := &domain.UserContext{
		User: &models.User{
			Id: 4,
		},
	}

	inspection := models.Inspection{
		Id:       1,
		Contexte: "test",
		Inspecteurs: []models.User{
			models.User{
				Id: 3,
			},
			models.User{
				Id: 4,
			},
		},
	}

	err := application.Repo.UpdateInspection(ctx, inspection)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementModificationInspection, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
}

func TestUpdateEtatInspectionHasCreatedNotification(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	ctxInspecteur1 := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}
	ctxInspecteur2 := &domain.UserContext{
		User: &models.User{
			Id:      5,
			Profile: models.ProfilInspecteur,
		},
	}
	ctxApprobateur := &domain.UserContext{
		User: &models.User{
			Id:      6,
			Profile: models.ProfilApprobateur,
		},
	}

	inspectionId := int64(2)

	err := application.Service.PublishInspection(ctxInspecteur1, inspectionId)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctxInspecteur2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementPublicationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)

	_, err = application.Service.CreateConstat(ctxInspecteur1, 3, models.Constat{
		Type: models.TypeConstatConforme,
	})
	assert.NoError(err)
	_, err = application.Service.CreateSuite(ctxInspecteur1, inspectionId, models.Suite{
		Type: models.TypeSuiteAucune,
	})
	assert.NoError(err)

	err = application.Service.AskValidateInspection(ctxInspecteur1, inspectionId)
	assert.NoError(err)

	notifications, err = application.Repo.ListNotifications(ctxInspecteur2, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(8), notification.Id)
	assert.Equal(models.EvenementDemandeValidationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)

	err = application.Service.RejectInspection(ctxApprobateur, inspectionId)
	assert.NoError(err)

	notifications, err = application.Repo.ListNotifications(ctxInspecteur2, nil)
	assert.NoError(err)
	assert.Equal(5, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(12), notification.Id)
	assert.Equal(models.EvenementRejetValidationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(6), notification.Evenement.AuteurId)

	err = application.Service.AskValidateInspection(ctxInspecteur1, inspectionId)
	assert.NoError(err)
	err = application.Service.ValidateInspection(ctxApprobateur, inspectionId)
	assert.NoError(err)

	notifications, err = application.Repo.ListNotifications(ctxInspecteur2, nil)
	assert.NoError(err)
	assert.Equal(7, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(17), notification.Id)
	assert.Equal(models.EvenementValidationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(6), notification.Evenement.AuteurId)
}
