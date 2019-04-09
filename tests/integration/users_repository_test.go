package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCheckUsersInProfilInspecteur(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	users := []models.User{
		models.User{
			Email:   "inspecteur1@filharmonic.com",
			Profile: models.ProfilInspecteur,
		},
		models.User{
			Email:   "inspecteur2@filharmonic.com",
			Profile: models.ProfilInspecteur,
		},
	}
	assert.NoError(application.DB.Insert(&users))

	checked, err := application.Repo.CheckUsersInspecteurs([]int64{1, 2})
	assert.True(checked)
	assert.NoError(err)
}

func TestCheckUsersWithProfilExploitant(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	users := []models.User{
		models.User{
			Email:   "inspecteur1@filharmonic.com",
			Profile: models.ProfilInspecteur,
		},
		models.User{
			Email:   "exploitant1@filharmonic.com",
			Profile: models.ProfilExploitant,
		},
	}
	assert.NoError(application.DB.Insert(&users))

	checked, err := application.Repo.CheckUsersInspecteurs([]int64{1, 2})
	assert.False(checked)
	assert.NoError(err)
}

func TestCheckUsersNotExisting(t *testing.T) {
	assert, application := tests.InitEmptyDB(t)

	checked, err := application.Repo.CheckUsersInspecteurs([]int64{4000})
	assert.False(checked)
	assert.NoError(err)
	checked, err = application.Repo.CheckUsersInspecteurs([]int64{4000, 4000})
	assert.False(checked)
	assert.NoError(err)
}
