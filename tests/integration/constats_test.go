package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestResoudreConstat(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	ctxInspecteur := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}

	err := application.Service.ResolveConstat(ctxInspecteur, int64(7))
	assert.NoError(err)
	inspection, err := application.Service.GetInspection(ctxInspecteur, int64(5))
	assert.NoError(err)
	assert.NotNil(inspection.PointsDeControle[0].Constat.DateResolution)
}

func TestResoudreConstatConforme(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	_, err := application.DB.Exec(`update inspections as i set etat = ? where id = ?`, models.EtatTraitementNonConformites, int64(3))
	assert.NoError(err)

	ctxInspecteur := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}

	err = application.Service.ResolveConstat(ctxInspecteur, int64(4))
	assert.Equal(domain.ErrBesoinTypeConstatNonConforme, err)
}

func TestResoudreConstatHorsTraitementNonConformites(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	ctxInspecteur := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}

	err := application.Service.ResolveConstat(ctxInspecteur, int64(3))
	assert.Equal(domain.ErrBesoinEtatTraitementNonConformites, err)
}
