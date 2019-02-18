package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/database/importcsv"
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestLoadInspecteursFromCSV(t *testing.T) {
	assert, application := tests.InitDB(t)

	err := importcsv.LoadInspecteursCSV("Liste_d_agents.mini.csv", application.DB)
	assert.NoError(err)
	inspecteurs, err := application.Repo.FindUsers(domain.ListUsersFilters{
		Inspecteurs:  true,
		Approbateurs: true,
	})
	const nbUsersInSeeds = 5
	assert.NoError(err)
	assert.NotEmpty(inspecteurs)
	assert.Len(inspecteurs, 30+nbUsersInSeeds)
	inspecteur := inspecteurs[len(inspecteurs)-1]
	assert.Equal("Rachel", inspecteur.Prenom)
	assert.Equal("BOUVARD", inspecteur.Nom)
	assert.Equal("rachel.bouvard@developpement-durable.gouv.fr", inspecteur.Email)
	assert.Equal(models.ProfilInspecteur, inspecteur.Profile)

	approbateur := inspecteurs[2+nbUsersInSeeds]
	assert.Equal(models.ProfilApprobateur, approbateur.Profile)
}
