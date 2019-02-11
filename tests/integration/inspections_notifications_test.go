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
	assert.Equal(int64(5), idInspection)

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
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}
	ctx2 := &domain.UserContext{
		User: &models.User{
			Id: 5,
		},
	}

	inspectionId := int64(2)

	err := application.Repo.UpdateEtatInspection(ctx, inspectionId, models.EtatEnCours)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementPublicationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)

	err = application.Repo.UpdateEtatInspection(ctx, inspectionId, models.EtatAttenteValidation)
	assert.NoError(err)

	notifications, err = application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(2, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(6), notification.Id)
	assert.Equal(models.EvenementDemandeValidationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)

	err = application.Repo.UpdateEtatInspection(ctx, inspectionId, models.EtatNonValide)
	assert.NoError(err)

	notifications, err = application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(3, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(9), notification.Id)
	assert.Equal(models.EvenementRejetValidationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)

	err = application.Repo.UpdateEtatInspection(ctx, inspectionId, models.EtatValide)
	assert.NoError(err)

	notifications, err = application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(10), notification.Id)
	assert.Equal(models.EvenementValidationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
}
