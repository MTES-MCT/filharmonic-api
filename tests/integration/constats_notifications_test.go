package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
	"github.com/MTES-MCT/filharmonic-api/util"
)

func TestCreateConstatHasCreatedNotification(t *testing.T) {
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
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Sujet:  "test",
				Publie: true,
			},
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctx1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctx2 := &domain.UserContext{
		User: &inspecteur2,
	}

	idPointDeControle := int64(1)

	constat := models.Constat{
		Type:      models.TypeConstatNonConforme,
		Remarques: "Ne respecte pas la réglementation",
	}

	idConstat, err := application.Repo.CreateConstat(ctx1, idPointDeControle, constat)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementCreationConstat, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur1.Id, notification.Evenement.AuteurId)
	assert.Equal(float64(idConstat), notification.Evenement.Data["constat_id"])
	assert.Equal(float64(idPointDeControle), notification.Evenement.Data["point_de_controle_id"])
}

func TestDeleteConstatHasCreatedNotification(t *testing.T) {

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
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Sujet:  "test",
				Publie: true,
				Constat: &models.Constat{
					Type: models.TypeConstatConforme,
				},
			},
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctx1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctx2 := &domain.UserContext{
		User: &inspecteur2,
	}

	idPointDeControle := int64(1)

	idConstat := int64(1)

	err := application.Repo.DeleteConstat(ctx1, idPointDeControle)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctx2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementSuppressionConstat, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur1.Id, notification.Evenement.AuteurId)
	assert.Equal(float64(idConstat), notification.Evenement.Data["constat_id"])
	assert.Equal(float64(idPointDeControle), notification.Evenement.Data["point_de_controle_id"])
}
