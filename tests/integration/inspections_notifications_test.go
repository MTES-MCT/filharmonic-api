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
	}

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}

	idInspection, err := application.Repo.CreateInspection(ctx, inspection)
	assert.NoError(err)
	assert.Equal(int64(5), idInspection)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.CreationInspection, notification.Evenement.Type)
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

	inspection := models.Inspection{
		Id:       1,
		Contexte: "test",
	}

	err := application.Repo.UpdateInspection(ctx, inspection)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.ModificationInspection, notification.Evenement.Type)
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

	inspectionId := int64(2)

	err := application.Repo.UpdateEtatInspection(ctx, inspectionId, models.EtatEnCours)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.PublicationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)

	err = application.Repo.UpdateEtatInspection(ctx, inspectionId, models.EtatAttenteValidation)
	assert.NoError(err)

	notifications, err = application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(5, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(5), notification.Id)
	assert.Equal(models.DemandeValidationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)

	err = application.Repo.UpdateEtatInspection(ctx, inspectionId, models.EtatNonValide)
	assert.NoError(err)

	notifications, err = application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(6, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(6), notification.Id)
	assert.Equal(models.RejetValidationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)

	err = application.Repo.UpdateEtatInspection(ctx, inspectionId, models.EtatValide)
	assert.NoError(err)

	notifications, err = application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(7, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(7), notification.Id)
	assert.Equal(models.ValidationInspection, notification.Evenement.Type)
	assert.Equal(inspectionId, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
}
