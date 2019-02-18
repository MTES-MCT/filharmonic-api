package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestListInpections(t *testing.T) {
	assert, application := tests.InitDB(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}
	inspections, err := application.Repo.ListInspections(ctx, domain.ListInspectionsFilter{})
	assert.NoError(err)
	assert.Len(inspections, 4)
	assert.Equal(0, inspections[0].NbMessagesNonLus)
	assert.Equal(1, inspections[1].NbMessagesNonLus)
	assert.Equal(0, inspections[2].NbMessagesNonLus)
	assert.Equal(0, inspections[3].NbMessagesNonLus)

	ctx = &domain.UserContext{
		User: &models.User{
			Id:      1,
			Profile: models.ProfilExploitant,
		},
	}
	inspections, err = application.Repo.ListInspections(ctx, domain.ListInspectionsFilter{})
	assert.NoError(err)
	assert.Len(inspections, 1)
	assert.Equal(1, inspections[0].NbMessagesNonLus)
}
