package tests

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/authentication/hash"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func initTestDB(db *database.Database, assert *require.Assertions) {
	encodedpassword1, err := hash.GenerateFromPassword("password1")
	assert.NoError(err)
	encodedpassword2, err := hash.GenerateFromPassword("password2")
	assert.NoError(err)
	encodedpassword3, err := hash.GenerateFromPassword("password3")
	assert.NoError(err)
	users := []interface{}{
		&models.User{
			Id:       1,
			Prenom:   "Michel",
			Nom:      "Exploitant1",
			Email:    "exploitant1@filharmonic.com",
			Password: encodedpassword1,
			Profile:  models.ProfilExploitant,
		},
		&models.User{
			Id:       2,
			Prenom:   "Bernard",
			Nom:      "Exploitant2",
			Email:    "exploitant2@filharmonic.com",
			Password: encodedpassword2,
			Profile:  models.ProfilExploitant,
		},
		&models.User{
			Id:       3,
			Prenom:   "Alain",
			Nom:      "Champion",
			Email:    "inspecteur1@filharmonic.com",
			Password: encodedpassword1,
			Profile:  models.ProfilInspecteur,
		},
		&models.User{
			Id:       4,
			Prenom:   "Corine",
			Nom:      "Dupont",
			Email:    "inspecteur2@filharmonic.com",
			Password: encodedpassword2,
			Profile:  models.ProfilInspecteur,
		},
		&models.User{
			Id:       5,
			Prenom:   "Bernard",
			Nom:      "Mars",
			Email:    "inspecteur3@filharmonic.com",
			Password: encodedpassword3,
			Profile:  models.ProfilInspecteur,
		},
		&models.User{
			Id:       6,
			Prenom:   "Albert",
			Nom:      "Approbe",
			Email:    "approbateur1@filharmonic.com",
			Password: encodedpassword1,
			Profile:  models.ProfilApprobateur,
		},
		&models.User{
			Id:       7,
			Prenom:   "Gilbert",
			Nom:      "Approbe",
			Email:    "approbateur2@filharmonic.com",
			Password: encodedpassword2,
			Profile:  models.ProfilApprobateur,
		},
	}
	err = db.Insert(users...)
	assert.NoError(err)

	etablissements := []interface{}{
		&models.Etablissement{
			Id:      1,
			S3IC:    "1234",
			Raison:  "Raison sociale",
			Adresse: "1 rue des fleurs 75000 Paris",
		},
		&models.Etablissement{
			Id:      2,
			S3IC:    "4567",
			Raison:  "Raison sociale 2",
			Adresse: "1 rue des plantes 44000 Nantes",
		},
		&models.Etablissement{
			Id:      3,
			S3IC:    "3335655",
			Raison:  "Raison sociale 3",
			Adresse: "1 rue des cordeliers 69000 Lyon",
		},
		&models.Etablissement{
			Id:      4,
			S3IC:    "4444213",
			Raison:  "Raison sociale 4",
			Adresse: "1 place de l'église 63000 Clermont-Ferrand",
		},
	}
	assert.NoError(db.Insert(etablissements...))

	etablissementToExploitants := []interface{}{
		&models.EtablissementToExploitant{
			EtablissementId: 1,
			UserId:          1,
		},
		&models.EtablissementToExploitant{
			EtablissementId: 2,
			UserId:          1,
		},
		&models.EtablissementToExploitant{
			EtablissementId: 3,
			UserId:          2,
		},
	}
	assert.NoError(db.Insert(etablissementToExploitants...))

	inspections := []interface{}{
		&models.Inspection{
			Id:   1,
			Date: date("2018-09-01"),
			Type: models.TypeApprofondi,
			Themes: []string{
				"Produits chimiques",
			},
			Annonce:         true,
			Origine:         models.OriginePlanControle,
			Etat:            models.EtatEnCours,
			Contexte:        "Emissions de NOx dépassant les seuils le 2/04/2005",
			EtablissementId: 1,
		},
		&models.Inspection{
			Id:   2,
			Date: date("2018-11-15"),
			Type: models.TypeCourant,
			Themes: []string{
				"Rejets dans l'eau",
				"Rejets dans l'air",
			},
			Annonce:         true,
			Origine:         models.OriginePlanControle,
			Etat:            models.EtatPreparation,
			Contexte:        "Incident cuve gaz le 30/12/2017",
			EtablissementId: 3,
		},
	}
	assert.NoError(db.Insert(inspections...))

	inspectionToInspecteurs := []interface{}{
		&models.InspectionToInspecteur{
			InspectionId: 1,
			UserId:       3,
		},
		&models.InspectionToInspecteur{
			InspectionId: 2,
			UserId:       3,
		},
	}
	assert.NoError(db.Insert(inspectionToInspecteurs...))

	pointsDeControle := []interface{}{
		&models.PointDeControle{
			Id:    1,
			Sujet: "Mesure des émissions atmosphériques canalisées par un organisme extérieur",
			ReferencesReglementaires: []string{
				"Article 3.2.3. de l'arrêté préfectoral du 28 juin 2017",
				"Article 3.2.8. de l'arrêté préfectoral du 28 juin 2017",
				"Article 8.2.1.2. de l'arrêté préfectoral du 28 juin 2017",
			},
			Publie:       true,
			InspectionId: 1,
		},
		&models.PointDeControle{
			Id:    2,
			Sujet: "Autosurveillance des émissions canalisées de COV",
			ReferencesReglementaires: []string{
				"Article 8.2.1.1. de l'arrêté préfectoral du 28 juin 2017",
			},
			Publie:       false,
			InspectionId: 1,
		},
		&models.PointDeControle{
			Id:    3,
			Sujet: "Eau - Air",
			ReferencesReglementaires: []string{
				"Article 1 de l'Arrêté ministériel du 28 avril 2014",
			},
			Publie:       false,
			InspectionId: 2,
		},
	}
	assert.NoError(db.Insert(pointsDeControle...))

	commentaires := []interface{}{
		&models.Commentaire{
			Id:           1,
			Message:      "Attention à l'article 243.",
			Date:         dateTime("2018-11-14T08:50:00"),
			AuteurId:     3,
			InspectionId: 1,
		},
		&models.Commentaire{
			Id:           2,
			Message:      "L'article 843 s'applique également.",
			Date:         dateTime("2018-11-16T16:50:00"),
			AuteurId:     4,
			InspectionId: 1,
		},
		&models.Commentaire{
			Id:           3,
			Message:      "Attention au précédent contrôle.",
			Date:         dateTime("2018-11-18T16:50:00"),
			AuteurId:     3,
			InspectionId: 2,
		},
	}
	assert.NoError(db.Insert(commentaires...))

	messages := []interface{}{
		&models.Message{
			Id:                1,
			Message:           "Auriez-vous l'obligeance de me fournir le document approprié ?",
			Date:              dateTime("2018-11-14T08:50:00"),
			Lu:                true,
			Interne:           false,
			AuteurId:          3,
			PointDeControleId: 1,
		},
		&models.Message{
			Id:                2,
			Message:           "Voici le document.",
			Date:              dateTime("2018-11-16T16:50:00"),
			Lu:                true,
			Interne:           false,
			AuteurId:          1,
			PointDeControleId: 1,
		},
		&models.Message{
			Id:                3,
			Message:           "Attention au précédent contrôle.",
			Date:              dateTime("2018-11-20T16:50:00"),
			Lu:                false,
			Interne:           true,
			AuteurId:          3,
			PointDeControleId: 1,
		},
		&models.Message{
			Id:                4,
			Message:           "Merci de me fournir le document.",
			Date:              dateTime("2018-11-21T16:50:00"),
			Lu:                false,
			Interne:           true,
			AuteurId:          4,
			PointDeControleId: 2,
		},
		&models.Message{
			Id:                5,
			Message:           "Auriez-vous l'obligeance de me fournir une photo de la cuve ?",
			Date:              dateTime("2018-11-18T17:50:00"),
			Lu:                true,
			Interne:           false,
			AuteurId:          3,
			PointDeControleId: 3,
		},
	}
	assert.NoError(db.Insert(messages...))
}

func date(datestr string) time.Time {
	date, err := time.Parse("2006-01-02", datestr)
	if err != nil {
		log.Fatal().Msgf("unable to parse date: %v", err)
	}
	return date
}

func dateTime(datestr string) time.Time {
	date, err := time.Parse("2006-01-02T15:04:05", datestr)
	if err != nil {
		log.Fatal().Msgf("unable to parse date: %v", err)
	}
	return date
}
