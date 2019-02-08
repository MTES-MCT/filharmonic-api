package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCreatePointDeControleHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}

	idInspection := int64(1)

	pointDeControle := models.PointDeControle{
		// Id:    7,
		Sujet: "Santé 3",
		ReferencesReglementaires: []string{
			"Article 1 de l'Arrêté ministériel du 28 avril 2014",
		},
	}

	idPointDeControle, err := application.Repo.CreatePointDeControle(ctx, idInspection, pointDeControle)
	assert.NoError(err)
	assert.Equal(int64(7), idPointDeControle)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[len(notifications)-1]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.CreationPointDeControle, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"point_de_controle_id": 7}`, notification.Evenement.Data)
}
func TestUpdatePointDeControleHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}

	idPointDeControle := int64(6)

	pointDeControle := models.PointDeControle{
		Sujet: "Santé 3",
	}

	err := application.Repo.UpdatePointDeControle(ctx, idPointDeControle, pointDeControle)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[len(notifications)-1]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.ModificationPointDeControle, notification.Evenement.Type)
	assert.Equal(int64(4), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"point_de_controle_id": 6}`, notification.Evenement.Data)
}

func TestDeletePointDeControleHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}

	idPointDeControle := int64(6)

	err := application.Repo.DeletePointDeControle(ctx, idPointDeControle)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[len(notifications)-1]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.SuppressionPointDeControle, notification.Evenement.Type)
	assert.Equal(int64(4), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"point_de_controle_id": 6}`, notification.Evenement.Data)
}
func TestPublishPointDeControleHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}

	idPointDeControle := int64(6)

	err := application.Repo.PublishPointDeControle(ctx, idPointDeControle)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[len(notifications)-1]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.PublicationPointDeControle, notification.Evenement.Type)
	assert.Equal(int64(4), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"point_de_controle_id": 6}`, notification.Evenement.Data)
}
