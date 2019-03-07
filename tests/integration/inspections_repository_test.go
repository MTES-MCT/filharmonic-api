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
	assert.Len(inspections, 5)
	assert.Equal(0, inspections[0].NbMessagesNonLus)
	assert.Equal(1, inspections[1].NbMessagesNonLus)
	assert.Equal(0, inspections[2].NbMessagesNonLus)
	assert.Equal(0, inspections[3].NbMessagesNonLus)
	assert.Equal(0, inspections[4].NbMessagesNonLus)

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

func TestGetInspectionTypesConstatsSuiteByID(t *testing.T) {
	assert, application := tests.InitDB(t)

	inspection, err := application.Repo.GetInspectionTypesConstatsSuiteByID(1)
	assert.NoError(err)
	assert.Equal(models.TypeSuiteObservation, inspection.Suite.Type)
	assert.Equal(2, len(inspection.PointsDeControle))
	pointDeControle := inspection.PointsDeControle[0]
	assert.Equal(models.TypeConstatObservation, pointDeControle.Constat.Type)
}

func TestGetRecapsValidation(t *testing.T) {
	assert, application := tests.InitDB(t)

	recaps, err := application.Repo.GetRecapsValidation(5)
	assert.NoError(err)
	assert.Len(recaps, 1)
}
