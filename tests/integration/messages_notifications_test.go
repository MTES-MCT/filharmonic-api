package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/util"

	"github.com/MTES-MCT/filharmonic-api/models"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCreateMessageHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspecteur := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	exploitant := models.User{
		Email:   "exploitant@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	assert.NoError(application.DB.Insert(&inspecteur, &exploitant))

	inspection := models.Inspection{
		Date: util.Date("2019-02-08"),
		Etat: models.EtatAttenteValidation,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
			Exploitants: []models.User{
				exploitant,
			},
		},
		Inspecteurs: []models.User{
			inspecteur,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
			},
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctxInspecteur := &domain.UserContext{
		User: &inspecteur,
	}
	ctxExploitant := &domain.UserContext{
		User: &exploitant,
	}

	message := models.Message{
		Message: "test",
		Interne: false,
	}

	idMessage, err := application.Repo.CreateMessage(ctxInspecteur, inspection.PointsDeControle[0].Id, message)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctxExploitant, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementCreationMessage, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur.Id, notification.Evenement.AuteurId)
	assert.Equal(float64(idMessage), notification.Evenement.Data["message_id"])
	assert.Equal(float64(inspection.PointsDeControle[0].Id), notification.Evenement.Data["point_de_controle_id"])
}

func TestCreateMessageInterneHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspecteur1 := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	inspecteur2 := models.User{
		Email:   "inspecteur2@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	exploitant := models.User{
		Email:   "exploitant@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	assert.NoError(application.DB.Insert(&inspecteur1, &inspecteur2, &exploitant))

	inspection := models.Inspection{
		Date: util.Date("2019-02-08"),
		Etat: models.EtatAttenteValidation,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
			Exploitants: []models.User{
				exploitant,
			},
		},
		Inspecteurs: []models.User{
			inspecteur1,
			inspecteur2,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
			},
		},
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))

	ctxInspecteur1 := &domain.UserContext{
		User: &inspecteur1,
	}
	ctxInspecteur2 := &domain.UserContext{
		User: &inspecteur2,
	}
	ctxExploitant := &domain.UserContext{
		User: &exploitant,
	}

	message := models.Message{
		Message: "test",
		Interne: true,
	}

	idMessage, err := application.Repo.CreateMessage(ctxInspecteur1, inspection.PointsDeControle[0].Id, message)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctxExploitant, nil)
	assert.NoError(err)
	assert.Equal(0, len(notifications))

	notifications, err = application.Repo.ListNotifications(ctxInspecteur2, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementCreationCommentaire, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur1.Id, notification.Evenement.AuteurId)
	assert.Equal(float64(idMessage), notification.Evenement.Data["message_id"])
	assert.Equal(float64(inspection.PointsDeControle[0].Id), notification.Evenement.Data["point_de_controle_id"])
}

func TestLireMessageHasCreatedNotification(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	inspecteur := models.User{
		Email:   "inspecteur1@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	exploitant := models.User{
		Email:   "exploitant@filharmonic.com",
		Profile: models.ProfilInspecteur,
	}
	assert.NoError(application.DB.Insert(&inspecteur, &exploitant))

	inspection := models.Inspection{
		Date: util.Date("2019-02-08"),
		Etat: models.EtatAttenteValidation,
		Etablissement: &models.Etablissement{
			Nom: "Équipement de pression",
			Exploitants: []models.User{
				exploitant,
			},
		},
		Inspecteurs: []models.User{
			inspecteur,
		},
		PointsDeControle: []models.PointDeControle{
			models.PointDeControle{
				Publie: true,
				Sujet:  "test1",
				Messages: []models.Message{
					models.Message{
						Message: "test",
						Interne: false,
						Auteur:  &inspecteur,
					},
				},
			},
		},
	}

	ctxInspecteur := &domain.UserContext{
		User: &inspecteur,
	}
	ctxExploitant := &domain.UserContext{
		User: &exploitant,
	}
	assert.NoError(tests.CreateInspection(application.DB, &inspection))
	err := application.Repo.LireMessage(ctxInspecteur, inspection.PointsDeControle[0].Messages[0].Id)
	assert.NoError(err)

	notifications, err := application.Repo.ListNotifications(ctxExploitant, nil)
	assert.NoError(err)
	assert.Equal(1, len(notifications))
	notification := notifications[0]
	assert.Equal(models.EvenementLectureMessage, notification.Evenement.Type)
	assert.Equal(inspection.Id, notification.Evenement.InspectionId)
	assert.Equal(inspecteur.Id, notification.Evenement.AuteurId)
	assert.Equal(float64(inspection.PointsDeControle[0].Messages[0].Id), notification.Evenement.Data["message_id"])
	assert.Equal(float64(inspection.PointsDeControle[0].Id), notification.Evenement.Data["point_de_controle_id"])
}
