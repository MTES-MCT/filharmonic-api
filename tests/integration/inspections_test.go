package integration

import (
	"testing"

	"github.com/go-pg/pg/types"
	"github.com/stretchr/testify/require"

	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
	"github.com/MTES-MCT/filharmonic-api/util"
)

func TestCloreInspectionConstatsResolus(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	inspectionId := initSeedsTestCloreInspectionConstatsResolus(assert, application.DB)

	ctxInspecteur := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}
	err := application.Service.CloreInspection(ctxInspecteur, inspectionId)
	assert.NoError(err)
}

func initSeedsTestCloreInspectionConstatsResolus(assert *require.Assertions, db *database.Database) int64 {
	inspection := models.Inspection{
		Date:            util.Date("2019-01-10"),
		Type:            models.TypeApprofondi,
		Etat:            models.EtatTraitementNonConformites,
		EtablissementId: 4,
	}
	assert.NoError(db.Insert(&inspection))

	constats := []models.Constat{
		models.Constat{
			Type:               models.TypeConstatNonConforme,
			DelaiNombre:        30,
			DelaiUnite:         "jours",
			EcheanceResolution: util.Date("2027-02-27"),
			DateResolution: types.NullTime{
				Time: util.DateTime("2019-02-01T14:04:05"),
			},
		},
		models.Constat{
			Type: models.TypeConstatConforme,
		},
		models.Constat{
			Type:               models.TypeConstatNonConforme,
			DelaiNombre:        20,
			DelaiUnite:         "jours",
			EcheanceResolution: util.Date("2019-02-17"),
			DateResolution: types.NullTime{
				Time: util.DateTime("2019-01-01T14:04:05"),
			},
		},
	}
	assert.NoError(db.Insert(&constats))

	pointsDeControle := []models.PointDeControle{
		models.PointDeControle{
			InspectionId: inspection.Id,
			Publie:       true,
			Sujet:        "test1",
			ConstatId:    constats[0].Id,
		},
		models.PointDeControle{
			InspectionId: inspection.Id,
			Publie:       true,
			Sujet:        "test2",
			ConstatId:    constats[1].Id,
		},
		models.PointDeControle{
			InspectionId: inspection.Id,
			Publie:       true,
			Sujet:        "test3",
			ConstatId:    constats[2].Id,
		},
	}
	assert.NoError(db.Insert(&pointsDeControle))

	inspectionToInspecteur := []models.InspectionToInspecteur{
		models.InspectionToInspecteur{
			InspectionId: inspection.Id,
			UserId:       3,
		},
	}
	assert.NoError(db.Insert(&inspectionToInspecteur))
	return inspection.Id
}

func TestCloreInspectionErrConstatNonResoluDelaisNonDepasses(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	inspectionId := initSeedsTestCloreInspectionErrConstatNonResoluDelaisNonDepasses(assert, application.DB)

	ctxInspecteur := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}
	err := application.Service.CloreInspection(ctxInspecteur, inspectionId)
	assert.Equal(domain.ErrClotureInspectionImpossible, err)
}

func initSeedsTestCloreInspectionErrConstatNonResoluDelaisNonDepasses(assert *require.Assertions, db *database.Database) int64 {
	inspection := models.Inspection{
		Date:            util.Date("2019-01-10"),
		Type:            models.TypeApprofondi,
		Etat:            models.EtatTraitementNonConformites,
		EtablissementId: 4,
	}
	assert.NoError(db.Insert(&inspection))

	constats := []models.Constat{
		models.Constat{
			Type:               models.TypeConstatNonConforme,
			DelaiNombre:        30,
			DelaiUnite:         "jours",
			EcheanceResolution: util.Date("2027-02-27"),
		},
		models.Constat{
			Type: models.TypeConstatConforme,
		},
		models.Constat{
			Type:               models.TypeConstatNonConforme,
			DelaiNombre:        20,
			DelaiUnite:         "jours",
			EcheanceResolution: util.Date("2019-02-17"),
			DateResolution: types.NullTime{
				Time: util.DateTime("2019-01-01T14:04:05"),
			},
		},
	}
	assert.NoError(db.Insert(&constats))

	pointsDeControle := []models.PointDeControle{
		models.PointDeControle{
			InspectionId: inspection.Id,
			Publie:       true,
			Sujet:        "test1",
			ConstatId:    constats[0].Id,
		},
		models.PointDeControle{
			InspectionId: inspection.Id,
			Publie:       true,
			Sujet:        "test2",
			ConstatId:    constats[1].Id,
		},
		models.PointDeControle{
			InspectionId: inspection.Id,
			Publie:       true,
			Sujet:        "test3",
			ConstatId:    constats[2].Id,
		},
	}
	assert.NoError(db.Insert(&pointsDeControle))

	inspectionToInspecteur := []models.InspectionToInspecteur{
		models.InspectionToInspecteur{
			InspectionId: inspection.Id,
			UserId:       3,
		},
	}
	assert.NoError(db.Insert(&inspectionToInspecteur))
	return inspection.Id
}

func TestCloreInspectionErrConstatNonResoluDelaisDepasses(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	inspectionId := initSeedsTestCloreInspectionErrConstatNonResoluDelaisDepasses(assert, application.DB)

	ctxInspecteur := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}
	err := application.Service.CloreInspection(ctxInspecteur, inspectionId)
	assert.NoError(err)
}

func initSeedsTestCloreInspectionErrConstatNonResoluDelaisDepasses(assert *require.Assertions, db *database.Database) int64 {
	inspection := models.Inspection{
		Date:            util.Date("2019-01-10"),
		Type:            models.TypeApprofondi,
		Etat:            models.EtatTraitementNonConformites,
		EtablissementId: 4,
	}
	assert.NoError(db.Insert(&inspection))

	constats := []models.Constat{
		models.Constat{
			Type:               models.TypeConstatNonConforme,
			DelaiNombre:        30,
			DelaiUnite:         "jours",
			EcheanceResolution: util.Date("2019-02-27"),
		},
		models.Constat{
			Type: models.TypeConstatConforme,
		},
		models.Constat{
			Type:               models.TypeConstatNonConforme,
			DelaiNombre:        20,
			DelaiUnite:         "jours",
			EcheanceResolution: util.Date("2019-02-17"),
			DateResolution: types.NullTime{
				Time: util.DateTime("2019-01-01T14:04:05"),
			},
		},
	}
	assert.NoError(db.Insert(&constats))

	pointsDeControle := []models.PointDeControle{
		models.PointDeControle{
			InspectionId: inspection.Id,
			Publie:       true,
			Sujet:        "test1",
			ConstatId:    constats[0].Id,
		},
		models.PointDeControle{
			InspectionId: inspection.Id,
			Publie:       true,
			Sujet:        "test2",
			ConstatId:    constats[1].Id,
		},
		models.PointDeControle{
			InspectionId: inspection.Id,
			Publie:       true,
			Sujet:        "test3",
			ConstatId:    constats[2].Id,
		},
	}
	assert.NoError(db.Insert(&pointsDeControle))

	inspectionToInspecteur := []models.InspectionToInspecteur{
		models.InspectionToInspecteur{
			InspectionId: inspection.Id,
			UserId:       3,
		},
	}
	assert.NoError(db.Insert(&inspectionToInspecteur))
	return inspection.Id
}
