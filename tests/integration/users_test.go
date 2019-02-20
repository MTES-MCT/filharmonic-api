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
	assert.NoError(err)
	assert.NotEmpty(inspecteurs)
	const nbUsersInSeeds = 5
	assert.Len(inspecteurs, 30+nbUsersInSeeds)
	inspecteur := inspecteurs[len(inspecteurs)-1]
	assert.Equal("Rachel", inspecteur.Prenom)
	assert.Equal("BOUVARD", inspecteur.Nom)
	assert.Equal("rachel.bouvard@developpement-durable.gouv.fr", inspecteur.Email)
	assert.Equal(models.ProfilInspecteur, inspecteur.Profile)

	approbateur := inspecteurs[2+nbUsersInSeeds]
	assert.Equal(models.ProfilApprobateur, approbateur.Profile)
}

func TestLoadExploitantsFromCSV(t *testing.T) {
	assert, application := tests.InitDB(t)

	err := importcsv.LoadExploitantsCSV("exploitants.mini.csv", application.DB)
	assert.NoError(err)
	exploitants, err := application.Repo.FindUsers(domain.ListUsersFilters{})
	assert.NoError(err)
	assert.NotEmpty(exploitants)
	const nbUsersInSeeds = 7
	assert.Len(exploitants, 2+nbUsersInSeeds)
	exploitant := exploitants[len(exploitants)-1]
	assert.Equal("Anne", exploitant.Prenom)
	assert.Equal("Exploitant4", exploitant.Nom)
	assert.Equal("exploitant4@filharmonic.com", exploitant.Email)
	assert.Equal(models.ProfilExploitant, exploitant.Profile)

	ctx := &domain.UserContext{
		User: &models.User{
			Id:      exploitant.Id,
			Profile: models.ProfilExploitant,
		},
	}
	etablissementsResults, err := application.Repo.FindEtablissements(ctx, domain.ListEtablissementsFilter{})
	assert.NoError(err)
	assert.Len(etablissementsResults.Etablissements, 1)
	assert.Equal("451267", etablissementsResults.Etablissements[0].S3IC)

}
