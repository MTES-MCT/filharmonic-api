package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/util"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAddCommentaireHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspecteur1 := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	inspecteur2 := models.User{
		Email:   "inspecteur2@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	assert.NoError(application.DB.Insert(&inspecteur1, &inspecteur2))

	inspection := models.Inspection{
		Date: util.Date("2019-01-10"),
		Etat: models.EtatEnCours,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
		},
		Inspecteurs: []models.User{
			inspecteur1,
			inspecteur2,
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctx1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctx2 := &domain.UserContext{
		User: &inspecteur2,
	}
	commentaire := models.Commentaire{
		Message: "test",
	}
	idCommentaire, err := application.Repo.CreateCommentaire(ctx1, inspection.Id, commentaire)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementCommentaireGeneral, notification.Evenement.Type)
	assert.Equal(float64(idCommentaire), notification.Evenement.Data["commentaire_id"])
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur1.Id, notification.Evenement.AuteurId)
}
