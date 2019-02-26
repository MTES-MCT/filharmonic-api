package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCreateSuiteHasCreatedNotification(t *testing.T) {
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

	suite := models.Suite{
		// Id: 3,
		Type:     models.TypeSuiteObservation,
		Synthese: "Observations à traiter",
	}

	idSuite, err := application.Repo.CreateSuite(ctx, idInspection, suite)
	assert.NoError(err)
	assert.Equal(int64(4), idSuite)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementCreationSuite, notification.Evenement.Type)
	assert.Equal(int64(1), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
}

func TestUpdateSuiteHasCreatedNotification(t *testing.T) {
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

	suite := models.Suite{
		Id:       1,
		Type:     models.TypeSuiteAucune,
		Synthese: "RAS",
	}

	err := application.Repo.UpdateSuite(ctx, idInspection, suite)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementModificationSuite, notification.Evenement.Type)
	assert.Equal(idInspection, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
}

func TestDeleteSuiteHasCreatedNotification(t *testing.T) {
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

	err := application.Repo.DeleteSuite(ctx, idInspection)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.EvenementSuppressionSuite, notification.Evenement.Type)
	assert.Equal(idInspection, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
}
