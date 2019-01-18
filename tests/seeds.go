package tests

import (
	"github.com/MTES-MCT/filharmonic-api/authentication/hash"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/models"
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
			Email:    "exploitant1@filharmonic.com",
			Password: encodedpassword1,
			Profile:  models.ProfilExploitant,
		},
		&models.User{
			Id:       2,
			Email:    "exploitant2@filharmonic.com",
			Password: encodedpassword2,
			Profile:  models.ProfilExploitant,
		},
		&models.User{
			Id:       3,
			Email:    "inspecteur1@filharmonic.com",
			Password: encodedpassword1,
			Profile:  models.ProfilInspecteur,
		},
		&models.User{
			Id:       4,
			Email:    "inspecteur2@filharmonic.com",
			Password: encodedpassword2,
			Profile:  models.ProfilInspecteur,
		},
		&models.User{
			Id:       5,
			Email:    "inspecteur3@filharmonic.com",
			Password: encodedpassword3,
			Profile:  models.ProfilInspecteur,
		},
		&models.User{
			Id:       6,
			Email:    "approbateur1@filharmonic.com",
			Password: encodedpassword1,
			Profile:  models.ProfilApprobateur,
		},
		&models.User{
			Id:       7,
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
}
