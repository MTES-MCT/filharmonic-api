package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestSearchEtablissements(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	ctxInspecteur := &domain.UserContext{
		User: &models.User{
			Id:      3,
			Profile: models.ProfilInspecteur,
		},
	}
	etablissementsInput := []models.Etablissement{
		models.Etablissement{
			Nom: "Équipement de pression",
		},
		models.Etablissement{
			Nom: "equipement de pression",
		},
		models.Etablissement{
			Nom: "pression",
		},
	}
	assert.NoError(application.DB.Insert(&etablissementsInput))
	filter := domain.ListEtablissementsFilter{
		Nom: "équipement",
	}
	etablissements, err := application.Service.ListEtablissements(ctxInspecteur, filter)
	assert.NoError(err)
	assert.Equal(2, etablissements.Total)
}
