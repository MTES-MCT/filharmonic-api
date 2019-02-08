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

	idInspection := int64(4)

	suite := models.Suite{
		// Id: 3,
		Type:     models.TypeSuiteObservation,
		Delai:    30,
		Synthese: "Observations à traiter",
	}

	idSuite, err := application.Repo.CreateSuite(ctx, idInspection, suite)
	assert.NoError(err)
	assert.Equal(int64(3), idSuite)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[len(notifications)-1]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.CreationSuite, notification.Evenement.Type)
	assert.Equal(int64(4), notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"suite_id": 3}`, notification.Evenement.Data)
}

func TestUpdateSuiteHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
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

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[len(notifications)-1]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.ModificationSuite, notification.Evenement.Type)
	assert.Equal(idInspection, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"suite_id": 1}`, notification.Evenement.Data)
}

func TestDeleteSuiteHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id: 3,
		},
	}

	idInspection := int64(1)

	err := application.Repo.DeleteSuite(ctx, idInspection)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx, nil)
	assert.NoError(err)
	assert.Equal(4, len(notifications))
	notification := notifications[len(notifications)-1]
	assert.Equal(int64(4), notification.Id)
	assert.Equal(models.SuppressionConstat, notification.Evenement.Type)
	assert.Equal(idInspection, notification.Evenement.InspectionId)
	assert.Equal(int64(3), notification.Evenement.AuteurId)
	assert.Equal(`{"suite_id": 1}`, notification.Evenement.Data)
}
