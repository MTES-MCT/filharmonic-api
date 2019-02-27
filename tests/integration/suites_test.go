package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCreateSuite(t *testing.T) {
	assert, application, _ := tests.InitService(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}

	constat := models.Constat{
		// Id: 5,
		Type: models.TypeConstatObservation,
	}

	_, err := application.Service.CreateConstat(ctx, int64(6), constat)
	assert.NoError(err)

	suite := models.Suite{
		// Id: 4,
		Type:     models.TypeSuiteObservation,
		Synthese: "Observations à traiter",
	}

	_, err = application.Service.CreateSuite(ctx, int64(4), suite)
	assert.NoError(err)
}

func TestCreateSuitePointsDeControleNonPublies(t *testing.T) {
	assert, application, _ := tests.InitService(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}

	pointDeControle := models.PointDeControle{
		Sujet:                    "test",
		ReferencesReglementaires: []string{"1", "2"},
	}
	_, err := application.Service.CreatePointDeControle(ctx, int64(4), pointDeControle)
	assert.NoError(err)

	constat := models.Constat{
		// Id: 5,
		Type: models.TypeConstatObservation,
	}
	_, err = application.Service.CreateConstat(ctx, int64(6), constat)
	assert.NoError(err)

	suite := models.Suite{
		// Id: 4,
		Type:     models.TypeSuiteObservation,
		Synthese: "Observations à traiter",
	}

	_, err = application.Service.CreateSuite(ctx, int64(4), suite)
	assert.Equal(models.ErrPointDeControleNonPublie, err)
}

func TestCreateSuitePointsDeControleSansConstat(t *testing.T) {
	assert, application, _ := tests.InitService(t)

	ctx := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}

	suite := models.Suite{
		// Id: 4,
		Type:     models.TypeSuiteObservation,
		Synthese: "Observations à traiter",
	}

	_, err := application.Service.CreateSuite(ctx, int64(4), suite)
	assert.Equal(models.ErrConstatManquant, err)
}
