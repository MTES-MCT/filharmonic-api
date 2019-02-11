package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAddCommentaireHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	commentaire := models.Commentaire{
		Message: "test",
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

	idCommentaire, err := application.Repo.CreateCommentaire(ctx, 1, commentaire)
	assert.NoError(err)
	assert.Equal(int64(4), idCommentaire)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementCommentaireGeneral, notification.Evenement.Type)
	assert.Equal(float64(4), notification.Evenement.Data["commentaire_id"])
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
}
