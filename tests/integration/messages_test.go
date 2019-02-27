package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCreateMessageEnCours(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	ctxExploitant := &domain.UserContext{
		User: &models.User{
			Id:      2,
			Profile: models.ProfilExploitant,
		},
	}
	message := models.Message{
		Message:       "Test",
		PiecesJointes: []models.PieceJointe{},
	}
	_, err := application.Service.CreateMessage(ctxExploitant, 6, message)
	assert.NoError(err)

	inspection, err := application.Service.GetInspection(ctxExploitant, 4)
	assert.NoError(err)
	assert.False(inspection.PointsDeControle[1].Messages[0].EtapeTraitementNonConformites)
}

func TestCreateMessageTraitementNonConformites(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	ctxExploitant := &domain.UserContext{
		User: &models.User{
			Id:      2,
			Profile: models.ProfilExploitant,
		},
	}
	message := models.Message{
		Message:       "Test",
		PiecesJointes: []models.PieceJointe{},
	}
	_, err := application.Service.CreateMessage(ctxExploitant, 7, message)
	assert.NoError(err)

	inspection, err := application.Service.GetInspection(ctxExploitant, 5)
	assert.NoError(err)
	assert.True(inspection.PointsDeControle[0].Messages[0].EtapeTraitementNonConformites)
}
