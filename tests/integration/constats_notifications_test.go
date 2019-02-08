package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCreateConstatHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}

	idPointDeControle := int64(1)

	constat := models.Constat{
		// Id: 4,
		Type:      models.TypeConstatNonConforme,
		Remarques: "Ne respecte pas la réglementation",
	}

	idConstat, err := application.Repo.CreateConstat(ctx, idPointDeControle, constat)
	assert.NoError(err)
	assert.Equal(int64(4), idConstat)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[len(notifications)-1]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.CreationConstat, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"constat_id": 4, "point_de_controle_id": 1}`, notification.Evenement.Data)
}

func TestDeleteConstatHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}

	idPointDeControle := int64(5)

	err := application.Repo.DeleteConstat(ctx, idPointDeControle)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[len(notifications)-1]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.SuppressionConstat, notification.Evenement.Type)
	assert.Equal(int64(4), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"constat_id": 3, "point_de_controle_id": 5}`, notification.Evenement.Data)
}
