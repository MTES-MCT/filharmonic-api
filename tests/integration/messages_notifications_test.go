package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCreateMessageHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}

	idPointDeControle := int64(1)

	message := models.Message{
		Message:           "test",
		Interne:           false,
		PointDeControleId: idPointDeControle,
	}

	idMessage, err := application.Repo.CreateMessage(ctx, idPointDeControle, message)
	assert.NoError(err)
	assert.Equal(int64(8), idMessage)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.CreationMessage, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"message_id": 8, "point_de_controle_id": 1}`, notification.Evenement.Data)

	message.Interne = true

	idMessage, err = application.Repo.CreateMessage(ctx, idPointDeControle, message)
	assert.NoError(err)
	assert.Equal(int64(9), idMessage)

	notifications, err = application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(5, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(5), notification.Id)
	assert.Equal(models.CreationCommentaire, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"message_id": 9, "point_de_controle_id": 1}`, notification.Evenement.Data)
}
func TestLireMessageHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}

	idMessage := int64(7)

	err := application.Repo.LireMessage(ctx, idMessage)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.LectureMessage, notification.Evenement.Type)
	assert.Equal(int64(2), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"message_id": 7, "point_de_controle_id": 3}`, notification.Evenement.Data)
}
