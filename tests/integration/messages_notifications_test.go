package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCreateMessageHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctxInspecteur := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}
	ctxExploitant := &domain.UserContext{
		User: &models.User{
			Id: 1,
		},
	}

	idPointDeControle := int64(1)

	message := models.Message{
		Message:           "test",
		Interne:           false,
		PointDeControleId: idPointDeControle,
	}

	idMessage, err := application.Repo.CreateMessage(ctxInspecteur, idPointDeControle, message)
	assert.NoError(err)
	assert.Equal(int64(8), idMessage)

	notifications, err := application.Repo.ListNotifications(ctxExploitant, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)

	message.Interne = true
	idMessage, err = application.Repo.CreateMessage(ctxInspecteur, idPointDeControle, message)
	assert.NoError(err)
	assert.Equal(int64(9), idMessage)

	notifications, err = application.Repo.ListNotifications(ctxExploitant, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification = notifications[0]
	assert.Equal(int64(4), notification.Id)
}
func TestLireMessageHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctxInspecteur := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}
	ctxExploitant := &domain.UserContext{
		User: &models.User{
			Id: 2,
		},
	}

	idMessage := int64(7)

	err := application.Repo.LireMessage(ctxInspecteur, idMessage)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctxExploitant, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementLectureMessage, notification.Evenement.Type)
	assert.Equal(int64(2), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(float64(7), notification.Evenement.Data["message_id"])
	assert.Equal(float64(3), notification.Evenement.Data["point_de_controle_id"])
}
