package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestCheckUsersInProfilInspecteur(t *testing.T) {
	assert, application := tests.InitDB(t)

	checked, err := application.Repo.CheckUsersInspecteurs([]int64{3, 4})
	assert.True(checked)
	assert.NoError(err)
}

func TestCheckUsersWithProfilExploitant(t *testing.T) {
	assert, application := tests.InitDB(t)

	checked, err := application.Repo.CheckUsersInspecteurs([]int64{1, 4})
	assert.False(checked)
	assert.NoError(err)
}

func TestCheckUsersNotExisting(t *testing.T) {
	assert, application := tests.InitDB(t)

	checked, err := application.Repo.CheckUsersInspecteurs([]int64{4000})
	assert.False(checked)
	assert.NoError(err)
	checked, err = application.Repo.CheckUsersInspecteurs([]int64{4000, 4000})
	assert.False(checked)
	assert.NoError(err)
}
