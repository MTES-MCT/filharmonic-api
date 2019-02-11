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
	ctx2 := &domain.UserContext{
		User: &models.User{
			Id: 4,
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

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementCreationConstat, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(float64(4), notification.Evenement.Data["constat_id"])
	assert.Equal(float64(1), notification.Evenement.Data["point_de_controle_id"])
}

func TestDeleteConstatHasCreatedNotification(t *testing.T) {
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

	err := application.Repo.DeleteConstat(ctx, idPointDeControle)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementSuppressionConstat, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(float64(1), notification.Evenement.Data["constat_id"])
	assert.Equal(float64(1), notification.Evenement.Data["point_de_controle_id"])
}
