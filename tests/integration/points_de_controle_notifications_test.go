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
	ctx2 := &domain.UserContext{
		User: &models.User{
			Id: 4,
		},
	}

	idInspection := int64(1)

	pointDeControle := models.PointDeControle{
		// Id:    8,
		Sujet: "Santé 3",
		ReferencesReglementaires: []string{
			"Article 1 de l'Arrêté ministériel du 28 avril 2014",
		},
	}

	idPointDeControle, err := application.Repo.CreatePointDeControle(ctx, idInspection, pointDeControle)
	assert.NoError(err)
	assert.Equal(int64(8), idPointDeControle)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementCreationPointDeControle, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(float64(8), notification.Evenement.Data["point_de_controle_id"])
}
func TestUpdatePointDeControleHasCreatedNotification(t *testing.T) {
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

	idPointDeControle := int64(1)

	pointDeControle := models.PointDeControle{
		Sujet: "Santé 3",
	}

	err := application.Repo.UpdatePointDeControle(ctx, idPointDeControle, pointDeControle)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementModificationPointDeControle, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(float64(1), notification.Evenement.Data["point_de_controle_id"])
}

func TestDeletePointDeControleHasCreatedNotification(t *testing.T) {
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

	idPointDeControle := int64(1)

	err := application.Repo.DeletePointDeControle(ctx, idPointDeControle)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementSuppressionPointDeControle, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(float64(1), notification.Evenement.Data["point_de_controle_id"])
}

func TestPublishPointDeControleHasCreatedNotification(t *testing.T) {
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

	idPointDeControle := int64(1)

	err := application.Repo.PublishPointDeControle(ctx, idPointDeControle)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementPublicationPointDeControle, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(float64(1), notification.Evenement.Data["point_de_controle_id"])
}
